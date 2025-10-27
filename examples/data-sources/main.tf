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

# Look up an existing engine
data "gemctl_engine" "existing_engine" {
  engine_id = "search-engine"
}

# Look up an existing data store
data "gemctl_data_store" "existing_store" {
  data_store_id = "document-store"
}

output "engine_details" {
  value = {
    name          = data.gemctl_engine.existing_engine.name
    display_name  = data.gemctl_engine.existing_engine.display_name
    solution_type = data.gemctl_engine.existing_engine.solution_type
    data_stores   = data.gemctl_engine.existing_engine.data_store_ids
  }
}

output "data_store_details" {
  value = {
    name           = data.gemctl_data_store.existing_store.name
    display_name   = data.gemctl_data_store.existing_store.display_name
    content_config = data.gemctl_data_store.existing_store.content_config
  }
}
