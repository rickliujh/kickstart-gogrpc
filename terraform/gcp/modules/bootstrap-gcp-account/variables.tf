variable "gcp_project_id" {
  description = "The GCP project ID where the resources will be created."
  type        = string
}

variable "gcp_region" {
  description = "The GCP region for the GCS bucket."
  type        = string
  default     = "europe-west1" # Example region, change as needed
}

variable "state_file_bucket_name" {
  description = "(Required) Name of the GCS bucket to store the state file"
  type        = string
}

variable "state_file_bucket_prefix" {
  description = "(Required) Key of the GCS bucket to store the state file"
  type        = string
  default     = "terraform/state"
}

variable "state_file_region" {
  description = "(Required) region of the GCS bucket to store the state file"
  type        = string
}

variable "tf_additional_providers" {
  description = "(Optional) List of additional Terraform providers"
  type = list(object({
    name             = string
    provider_source  = string
    provider_version = string
  }))
  default = []
}

variable "override_tags" {
  description = "(Optional) Override tags to apply to resources"
  type        = map(string)
  default     = null
}

variable "override_tf_version" {
  description = "(Optional) Override version of Terraform to use, defaults to 1.9.7 if not set"
  type        = string
  default     = null
}

variable "override_gcp_provider_version" {
  description = "(Optional) Override version of GCP provider to use, defaults to 5.70.0 if not set"
  type        = string
  default     = null
}

variable "override_local_provider_version" {
  description = "(Optional) Override version of local provider to use, defaults to 2.5.2 if not set"
  type        = string
  default     = null
}
