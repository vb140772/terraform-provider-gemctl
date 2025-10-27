terraform {
  required_providers {
    gemctl = {
      source  = "vb140772/gemctl"
      version = "~> 0.1"
  }
}

provider "gemctl" {
  project_id = "your-project-id"
  location   = "us"
}

# Example 1: Using a data source to reference an existing engine
data "gemctl_engine" "existing" {
  engine_id = "my-existing-engine"
}

# Example 2: Using a data source to reference an existing data store
data "gemctl_data_store" "existing" {
  data_store_id = "my-existing-store"
}

# Example 3: Output the engine details
output "existing_engine_info" {
  value = {
    name         = data.gemctl_engine.existing.name
    display_name = data.gemctl_engine.existing.display_name
    solution_type = data.gemctl_engine.existing.solution_type
  }
}

# Example 4: Output the data store details
output "existing_data_store_info" {
  value = {
    name          = data.gemctl_data_store.existing.name
    display_name  = data.gemctl_data_store.existing.display_name
    content_config = data.gemctl_data_store.existing.content_config
  }
}

# Example 5: Reference a data source in another resource
resource "gemctl_engine" "new_engine" {
  engine_id    = "new-engine-with-existing-store"
  display_name = "New Engine with Existing Store"
  data_stores  = [data.gemctl_data_store.existing.data_store_id]
}
