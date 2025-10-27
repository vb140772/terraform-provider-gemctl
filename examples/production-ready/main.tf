terraform {
  required_providers {
    gemctl = {
      source  = "vb140772/gemctl"
      version = "~> 0.1"
    }
  }
}

provider "gemctl" {
  project_id = "your-project-id"
  location   = "us"
}

# Example: Workflow for creating a production-ready search infrastructure

# 1. Create data stores from different GCS buckets
resource "gemctl_data_store" "production_docs" {
  data_store_id = "prod-document-store"
  display_name  = "Production Document Store"
  gcs_uri       = "gs://prod-documents-bucket/*"
}

resource "gemctl_data_store" "staging_docs" {
  data_store_id = "staging-document-store"
  display_name  = "Staging Document Store"
  gcs_uri       = "gs://staging-documents-bucket/*"
}

# 2. Create engines for different environments
resource "gemctl_engine" "production" {
  engine_id    = "prod-search-engine"
  display_name = "Production Search Engine"
  data_stores  = [gemctl_data_store.production_docs.id]
}

resource "gemctl_engine" "staging" {
  engine_id    = "staging-search-engine"
  display_name = "Staging Search Engine"
  data_stores  = [gemctl_data_store.staging_docs.id]
}

# 3. Use outputs to get resource information
output "production_engine" {
  value = {
    engine_id    = gemctl_engine.production.engine_id
    full_name    = gemctl_engine.production.name
    display_name = gemctl_engine.production.display_name
    data_stores  = gemctl_engine.production.data_stores
  }
}

output "staging_engine" {
  value = {
    engine_id    = gemctl_engine.staging.engine_id
    full_name    = gemctl_engine.staging.name
    display_name = gemctl_engine.staging.display_name
    data_stores  = gemctl_engine.staging.data_stores
  }
}

# 4. Use variables for flexibility
variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
  default     = "dev"
}

variable "gcs_bucket" {
  description = "GCS bucket for data"
  type        = string
}

# This example uses variables to create environment-specific resources
resource "gemctl_data_store" "env_docs" {
  count        = var.environment == "prod" ? 1 : 0
  data_store_id = "env-document-store"
  display_name  = "${var.environment} Document Store"
  gcs_uri       = "gs://${var.gcs_bucket}/*"
}
