terraform {
  required_providers {
    gemctl = {
      source  = "vb140772/gemctl"
      version = "~> 0.1"
    }
  }
}

provider "gemctl" {
  project_id = var.project_id
  location   = var.location
}

# Variables
variable "project_id" {
  description = "Google Cloud project ID"
  type        = string
}

variable "location" {
  description = "Location for resources"
  type        = string
  default     = "us"
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "data_store_id" {
  description = "Data store ID"
  type        = string
}

variable "display_name" {
  description = "Display name for resources"
  type        = string
}

variable "gcs_uri" {
  description = "GCS URI for data import"
  type        = string
}

# Resources
resource "gemctl_data_store" "main" {
  data_store_id = var.data_store_id
  display_name  = var.display_name
  gcs_uri       = var.gcs_uri
}

resource "gemctl_engine" "main" {
  engine_id    = "${var.environment}-search-engine"
  display_name = "${var.display_name} Search Engine"
  data_stores  = [gemctl_data_store.main.id]
}

# Outputs
output "data_store" {
  value = {
    id   = gemctl_data_store.main.id
    name = gemctl_data_store.main.name
  }
}

output "engine" {
  value = {
    id   = gemctl_engine.main.id
    name = gemctl_engine.main.name
  }
}
