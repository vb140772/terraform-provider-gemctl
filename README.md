# Terraform Provider for Gemini Enterprise (gemctl)

A Terraform provider for managing Google Gemini Enterprise (formerly Agentspace) resources, including engines and data stores.

## Documentation

Full provider documentation is available in the [docs/](docs/) directory:

- [Provider Documentation](docs/index.md)
- [Resources](docs/resources/)
  - [gemctl_engine](docs/resources/engine.md)
  - [gemctl_data_store](docs/resources/data_store.md)
- [Data Sources](docs/data-sources/)
  - [gemctl_engine](docs/data-sources/engine.md)
  - [gemctl_data_store](docs/data-sources/data_store.md)

## Features

- **Engine Management**: Create, read, update, and delete search engines
- **Data Store Management**: Import data from GCS buckets and manage data stores
- **Data Sources**: Look up existing engines and data stores
- **Full CRUD Operations**: Complete lifecycle management of resources

## Installation

### Prerequisites

- Terraform >= 1.0
- Google Cloud credentials configured
- Discovery Engine API enabled

### Configuration

Set up your provider configuration:

```hcl
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
  location   = "us" # or "global"
}
```

### Provider Arguments

- `project_id` (Required): Your Google Cloud project ID
- `location` (Optional): Location for resources. Defaults to "us"
- `collection` (Optional): Collection ID. Defaults to "default_collection"
- `use_service_account` (Optional): Use service account credentials. Defaults to false (uses user credentials)

## Resources

### gemctl_engine

Manages a search engine in Gemini Enterprise.

**Arguments:**

- `engine_id` (Required): Unique identifier for the engine
- `display_name` (Required): Display name for the engine
- `data_stores` (Optional): List of data store IDs to connect to this engine

**Attributes:**

- `id`: The engine ID
- `name`: Full resource name of the engine

**Example:**

```hcl
resource "gemctl_engine" "my_engine" {
  engine_id    = "my-search-engine"
  display_name = "My Search Engine"
  data_stores  = ["my-data-store"]
}
```

### gemctl_data_store

Manages a data store in Gemini Enterprise.

**Arguments:**

- `data_store_id` (Required): Unique identifier for the data store
- `display_name` (Required): Display name for the data store
- `gcs_uri` (Required): GCS URI to import data from (e.g., `gs://bucket/path/*`)

**Attributes:**

- `id`: The data store ID
- `name`: Full resource name of the data store

**Example:**

```hcl
resource "gemctl_data_store" "my_store" {
  data_store_id = "my-data-store"
  display_name  = "My Data Store"
  gcs_uri       = "gs://my-bucket/data/*"
}
```

## Data Sources

### gemctl_engine

Retrieves information about an existing engine.

**Arguments:**

- `engine_id` (Required): Engine ID to look up

**Attributes:**

- `name`: Full resource name
- `display_name`: Display name
- `solution_type`: Solution type (e.g., SOLUTION_TYPE_SEARCH)
- `industry_vertical`: Industry vertical (e.g., GENERIC)
- `data_store_ids`: List of connected data store IDs

**Example:**

```hcl
data "gemctl_engine" "existing" {
  engine_id = "my-engine"
}

output "engine_name" {
  value = data.gemctl_engine.existing.name
}
```

### gemctl_data_store

Retrieves information about an existing data store.

**Arguments:**

- `data_store_id` (Required): Data store ID to look up

**Attributes:**

- `name`: Full resource name
- `display_name`: Display name
- `industry_vertical`: Industry vertical
- `content_config`: Content configuration
- `create_time`: Creation timestamp

**Example:**

```hcl
data "gemctl_data_store" "existing" {
  data_store_id = "my-store"
}

output "store_name" {
  value = data.gemctl_data_store.existing.name
}
```

## Examples

We provide comprehensive examples in the [examples/](examples/) directory:

### Basic Examples

1. **[Simple Engine](examples/simple-engine/main.tf)** - Create a basic search engine
2. **[Simple Data Store](examples/simple-datastore/main.tf)** - Create a data store from GCS
3. **[Engine with Data Store](examples/engine-with-datastore/main.tf)** - Connect an engine to a data store

### Advanced Examples

4. **[Data Sources](examples/data-sources/main.tf)** - Look up existing resources
5. **[Use Existing Resources](examples/use-existing-resources/main.tf)** - Reference existing resources in new ones
6. **[Multiple Data Stores](examples/multiple-stores/main.tf)** - Manage multiple data stores and engines
7. **[Production Ready](examples/production-ready/main.tf)** - Production setup with environments
8. **[Modular Configuration](examples/modular/)** - Use variables and tfvars files

### Complete Example

```hcl
terraform {
  required_providers {
    gemctl = {
      source  = "vb140772/gemctl"
      version = "~> 0.1"
    }
  }
}

provider "gemctl" {
  project_id = "bcbsma-ailab-agentspace-01"
  location   = "us"
}

# Create a data store
resource "gemctl_data_store" "documents" {
  data_store_id = "document-store"
  display_name   = "Document Store"
  gcs_uri        = "gs://my-bucket/documents/*"
}

# Create an engine with the data store
resource "gemctl_engine" "search" {
  engine_id    = "search-engine"
  display_name = "Search Engine"
  data_stores  = [gemctl_data_store.documents.id]
}

# Look up existing resources
data "gemctl_engine" "existing" {
  engine_id = "search-engine"
}

output "engine_details" {
  value = {
    name         = data.gemctl_engine.existing.name
    display_name = data.gemctl_engine.existing.display_name
  }
}
```

## Authentication

The provider supports two authentication methods:

1. **User Credentials (Default)**: Uses `gcloud auth print-access-token`
   ```hcl
   provider "gemctl" {
     use_service_account = false
   }
   ```

2. **Service Account**: Uses Application Default Credentials (ADC)
   ```hcl
   provider "gemctl" {
     use_service_account = true
   }
   ```

## Local Development

For local development, configure dev overrides in `~/.terraformrc`:

```
provider_installation {
  dev_overrides {
    "vb140772/gemctl" = "/path/to/terraform-provider-gemctl"
  }
  direct_install_mirrors = ["registry.terraform.io"]
}
```

## Requirements

| Name | Version |
|------|---------|
| terraform | >= 1.0 |
| go | >= 1.21 |

## Building

```bash
go build -o terraform-provider-gemctl
```

## Publishing to Terraform Registry

To publish this provider to the [Terraform Registry](https://registry.terraform.io), follow these steps:

### Prerequisites

1. Sign in to the [Terraform Registry](https://registry.terraform.io/sign-in)
2. Create a namespace `vb140772` or use an existing one
3. Ensure your GitHub repository is public

### Creating a Release

1. **Create a git tag:**
   ```bash
   git tag v0.1.0
   git push origin v0.1.0
   ```

2. **GitHub Actions will automatically:**
   - Build the provider
   - Generate documentation
   - Create a GitHub release
   - Publish to Terraform Registry (after linking)

3. **Link the repository to the registry:**
   - Go to [Terraform Registry](https://registry.terraform.io)
   - Sign in and navigate to your namespace
   - Click "Link a Provider" and select your repository
   - The registry will automatically import the provider from the GitHub release

### Workflows

The repository includes GitHub Actions workflows:

- **CI** (`.github/workflows/ci.yml`): Runs tests and checks on every push
- **Release** (`.github/workflows/release.yml`): Uses [GoReleaser](https://goreleaser.com/) to build and publish releases

The release workflow uses [GoReleaser](https://goreleaser.com/) for automated builds and releases, following the [Terraform Registry publishing guidelines](https://developer.hashicorp.com/terraform/registry/providers/publishing#using-goreleaser-locally).

### Local Release (Development)

To test releases locally with GoReleaser:

```bash
# Install GoReleaser
brew install goreleaser/tap/goreleaser

# Test the configuration
goreleaser check

# Build a snapshot release (no git tag required)
goreleaser release --snapshot

# Or build a full release
goreleaser release
```

This is useful for testing the release process before pushing tags to GitHub.

## License

MIT License - Copyright (c) 2024

## Support

For issues and questions:
- GitHub: https://github.com/vb140772/terraform-provider-gemctl
- Issues: https://github.com/vb140772/terraform-provider-gemctl/issues
