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

data "google_project" "proj" {
  project_id = local.name
}

provider "google" {
  region  = local.region
  project = data.google_project.proj.id
}

resource "google_project" "proj" {
  name            = local.name
  project_id      = local.name
  billing_account = var.billing_account
}

resource "google_artifact_registry_repository" "gar" {
  location      = local.region
  repository_id = local.name
  description   = "${local.name} docker repository"
  format        = "DOCKER"

  docker_config {
    immutable_tags = true
  }
}

resource "google_cloud_run_v2_service" "kickstart_svc" {
  name                = "${local.name}-service"
  location            = local.region
  deletion_protection = false

  template {
    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello"
      resources {
        limits = {
          cpu    = "1"
          memory = "50Mi"
        }
      }
    }
  }
}

variable "billing_account" {
  type        = string
  description = "0XXX0-0XXX0-0XXX0, Replace with your billing account ID"
}

