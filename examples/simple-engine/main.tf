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

# Create a basic engine without data stores
resource "gemctl_engine" "basic_engine" {
  engine_id    = "basic-search-engine"
  display_name = "Basic Search Engine"
  data_stores  = []
}

output "engine_name" {
  value = gemctl_engine.basic_engine.name
}
