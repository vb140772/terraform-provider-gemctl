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

# Create a data store from GCS
resource "gemctl_data_store" "documents" {
  data_store_id = "document-store"
  display_name  = "Document Store"
  gcs_uri       = "gs://your-bucket/documents/*"
}

# Create an engine with the data store
resource "gemctl_engine" "search" {
  engine_id    = "search-engine"
  display_name = "Search Engine"
  data_stores  = [gemctl_data_store.documents.id]
}

output "engine_info" {
  value = {
    name         = gemctl_engine.search.name
    display_name = gemctl_engine.search.display_name
  }
}

output "data_store_info" {
  value = {
    name = gemctl_data_store.documents.name
    id   = gemctl_data_store.documents.id
  }
}
