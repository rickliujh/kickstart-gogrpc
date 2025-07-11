terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "6.42.0"
    }
  }
}

locals {
  region = "us-central1"
  name   = "go-kickstart"
}

provider "google" {
  region = local.region
}

resource "google_project" "proj" {
  name            = local.name
  project_id      = local.name
  billing_account = var.billing_account
}

resource "google_project_service" "gar" {
  project = google_project.proj.project_id
  service = "artifactregistry.googleapis.com"

  timeouts {
    create = "30m"
    update = "40m"
  }

  disable_on_destroy = true
}

resource "google_project_service" "crun" {
  project = google_project.proj.project_id
  service = "run.googleapis.com"

  timeouts {
    create = "30m"
    update = "40m"
  }

  disable_on_destroy = true
}

module "bootstrap" {
  source                 = "github.com/rickliujh/kickstart-gogrpc//terraform/gcp/modules/bootstrap-gcp-account"
  state_file_region      = local.region
  state_file_bucket_name = "kickstart-tf-state-alpha"
  gcp_project_id         = google_project.proj.project_id
  gcp_region             = local.region
}

resource "google_artifact_registry_repository" "gar" {
  project       = google_project.proj.project_id
  location      = local.region
  repository_id = local.name
  description   = "${local.name} docker repository"
  format        = "DOCKER"

  docker_config {
    immutable_tags = true
  }

  depends_on = [
    google_project_service.crun
  ]
}

resource "google_cloud_run_v2_service" "kickstart_svc" {
  project             = google_project.proj.project_id
  name                = "${local.name}-service"
  location            = local.region
  deletion_protection = false

  template {
    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello"
      resources {
        limits = {
          cpu    = "1.0"
          memory = "128Mi"
        }
        startup_cpu_boost = true
        cpu_idle          = true
      }
    }
  }

  depends_on = [
    google_project_service.crun
  ]
}

variable "billing_account" {
  type        = string
  description = "0XXX0-0XXX0-0XXX0, Replace with your billing account ID"
}

