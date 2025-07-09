output "terraform_state_bucket_name" {
  value       = google_storage_bucket.terraform_state_bucket.name
  description = "The name of the GCS bucket used for Terraform state."
}

