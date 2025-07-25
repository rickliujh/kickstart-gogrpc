terraform {
  # required_providers {
  #   google = {
  #     source  = "hashicorp/google"
  #     version = "6.42.0"
  #   }
  # }
}

locals {
  region      = "us-central1"
  name        = "go-kickstart"
  github_org  = "rickliujh"
  github_repo = "kickstart-gogrpc"
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
  source                 = "github.com/rickliujh/tf-tmpl//gcp/modules/bootstrap-gcp-account"
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
      image = "us-central1-docker.pkg.dev/go-kickstart/go-kickstart/kickstart-server:v0.0.1"

      ports {
        name           = "http1"
        container_port = 8080
      }

      resources {
        limits = {
          cpu    = "1.0"
          memory = "128Mi"
        }
        startup_cpu_boost = true
        cpu_idle          = true
      }

      args = [
        "server",
        "http",
        "-a",
        ":8080",
        "-e",
        "dev",
        "-n",
        "gcp-http",
        "-l",
        "INFO",
        "-c",
        " ",
      ]

      startup_probe {
        initial_delay_seconds = 5
        timeout_seconds       = 1
        period_seconds        = 30
        failure_threshold     = 3
        tcp_socket {
          port = 8080
        }
      }

      liveness_probe {
        http_get {
          path = "/ping"
        }
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

module "github_actions" {
  source                 = "github.com/rickliujh/tf-tmpl//gcp/modules/github-actions"
  project_id             = google_project.proj.project_id
  github_org             = "rickliujh"
  github_repo            = "kickstart-gogrpc"
  github_org_id          = "36358701"
  artifect_repository_id = google_artifact_registry_repository.gar.repository_id
  cloud_run_service_name = google_cloud_run_v2_service.kickstart_svc.name
  override_wif_pool_id   = "github-actions-pool2"
}

data "google_secret_manager_secret_version" "github_token" {
  project = google_project.proj.project_id
  secret  = "kickstart-gogrpc-ci-github-token"
}

