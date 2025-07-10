#--------------------------------------------#
# Using locals instead of hard-coding strings
#--------------------------------------------#
locals {
  tf_version = coalesce(var.override_tf_version, "1.12.1")

  tags = coalesce(var.override_tags, {
    name    = "tf-bootstrap",
    module  = "bootstrap-gcp-account",
    purpose = "terraform-state"
  })

  provider_config = concat(var.tf_additional_providers, [
    {
      name             = "google"
      provider_source  = "hashicorp/google"
      provider_version = coalesce(var.override_gcp_provider_version, "6.43.0")
    },
    {
      name             = "local"
      provider_source  = "hashicorp/local"
      provider_version = coalesce(var.override_local_provider_version, "2.5.3")
    }
  ])
}

provider "google" {
  project = var.gcp_project_id
  region  = var.gcp_region
}

# Resource to create the GCS bucket for Terraform state
resource "google_storage_bucket" "terraform_state_bucket" {
  name          = var.state_file_bucket_name
  location      = var.state_file_region # Or a multi-region location like "US"
  force_destroy = false                 # Set to true ONLY when you want to delete the bucket and its contents

  # Enable versioning to keep a history of your state files
  versioning {
    enabled = true
  }

  # Enable uniform bucket-level access for simplified IAM permissions
  uniform_bucket_level_access = true

  # Optional: Add lifecycle rules to manage old versions
  lifecycle_rule {
    action {
      type = "Delete"
    }
    condition {
      num_newer_versions = 5 # Keep the 5 most recent versions
    }
  }

  # Optional: Prevent public access
  public_access_prevention = "enforced"

  # Labels for organization
  labels = local.tags
}

# Create the terraform backend configuration - the catch 22 is that you need infrastructure to 
# store the state file before you can automate your infrastructure. The approach needs 2 steps:
# 1. Create the S3 bucket and DynamoDB table to store the state, and generate the backend 
#    config for terraform to use in the terraform.tf file.
# 2. For the 2nd run, it will now use this config and migrate the local state file to S3.
resource "local_file" "terraform_tf" {
  filename = "${path.root}/terraform.tf"
  content = templatefile("${path.module}/templates/terraform.tf.tmpl", {
    state_file_bucket_name   = var.state_file_bucket_name
    state_file_region        = var.state_file_region
    state_file_bucket_prefix = var.state_file_bucket_prefix
    tf_version               = local.tf_version
    providers                = local.provider_config
  })
  directory_permission = "0666"
  file_permission      = "0666"
}

