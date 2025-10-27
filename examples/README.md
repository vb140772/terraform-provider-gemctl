# Terraform Examples

This directory contains example Terraform configurations demonstrating how to use the gemctl provider.

## Examples

### 1. Simple Engine (`simple-engine/`)

Creates a basic search engine without any data stores.

**Usage:**
```bash
cd examples/simple-engine
terraform init
terraform plan
terraform apply
```

### 2. Simple Data Store (`simple-datastore/`)

Creates a data store that imports data from a GCS bucket.

**Usage:**
```bash
cd examples/simple-datastore
# Update the gcs_uri in main.tf with your bucket
terraform init
terraform plan
terraform apply
```

### 3. Engine with Data Store (`engine-with-datastore/`)

Creates a data store and an engine connected to it.

**Usage:**
```bash
cd examples/engine-with-datastore
# Update the gcs_uri in main.tf with your bucket
terraform init
terraform plan
terraform apply
```

### 4. Data Sources (`data-sources/`)

Demonstrates how to use data sources to look up existing engines and data stores.

**Usage:**
```bash
cd examples/data-sources
# Update engine_id and data_store_id with your existing resources
terraform init
terraform plan
```

### 5. Use Existing Resources (`use-existing-resources/`)

Shows how to reference existing engines and data stores using data sources, and how to use them in new resources.

**Usage:**
```bash
cd examples/use-existing-resources
# Update the engine_id and data_store_id with your existing resources
terraform init
terraform plan
```

### 6. Multiple Data Stores (`multiple-stores/`)

Creates multiple data stores and engines that connect to them, demonstrating organization of different content types.

**Usage:**
```bash
cd examples/multiple-stores
# Update the GCS URIs in main.tf
terraform init
terraform plan
terraform apply
```

### 7. Production Ready (`production-ready/`)

Demonstrates a production-ready setup with separate environments, multiple data stores, and comprehensive outputs.

**Usage:**
```bash
cd examples/production-ready
# Update project_id and GCS buckets
terraform init
terraform plan
terraform apply
```

### 8. Modular Configuration (`modular/`)

Demonstrates a modular approach using variables and separate tfvars files for different environments.

**Usage:**
```bash
cd examples/modular
# Apply to development
terraform apply -var-file=dev.tfvars

# Apply to production
terraform apply -var-file=prod.tfvars
```

## Configuration

Before running any example, update the provider configuration in `main.tf`:

```hcl
provider "gemctl" {
  project_id = "your-project-id"  # Update this
  location   = "us"
}
```

## Authentication

Make sure you're authenticated with Google Cloud:

```bash
gcloud auth application-default login
```

Or use a service account:

```bash
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account.json
```

## Local Development

For local development, configure dev overrides in `~/.terraformrc`:

```
provider_installation {
  dev_overrides {
    "vb140772/gemctl" = "/path/to/terraform-provider-gemctl"
  }
}
```

## Clean Up

To remove all resources:

```bash
terraform destroy
```
