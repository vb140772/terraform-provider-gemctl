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

# Multiple data stores for different content types
resource "gemctl_data_store" "documents" {
  data_store_id = "document-store"
  display_name  = "Document Store"
  gcs_uri       = "gs://your-bucket/documents/*"
}

resource "gemctl_data_store" "presentations" {
  data_store_id = "presentation-store"
  display_name  = "Presentation Store"
  gcs_uri       = "gs://your-bucket/presentations/*"
}

resource "gemctl_data_store" "videos" {
  data_store_id = "video-store"
  display_name  = "Video Store"
  gcs_uri       = "gs://your-bucket/videos/*"
}

# Create an engine that connects to all data stores
resource "gemctl_engine" "unified_search" {
  engine_id    = "unified-search-engine"
  display_name = "Unified Search Engine"
  data_stores  = [
    gemctl_data_store.documents.id,
    gemctl_data_store.presentations.id,
    gemctl_data_store.videos.id
  ]
}

# Create a separate engine for document-only search
resource "gemctl_engine" "document_search" {
  engine_id    = "document-search-engine"
  display_name = "Document Search Engine"
  data_stores  = [gemctl_data_store.documents.id]
}

# Output all the resources
output "unified_engine" {
  value = {
    name         = gemctl_engine.unified_search.name
    data_stores = gemctl_engine.unified_search.data_stores
  }
}

output "document_engine" {
  value = {
    name         = gemctl_engine.document_search.name
    data_stores = gemctl_engine.document_search.data_stores
  }
}
