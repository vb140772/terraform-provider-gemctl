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

output "data_store_name" {
  value = gemctl_data_store.documents.name
}
