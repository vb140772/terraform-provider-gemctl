# Modular Infrastructure Example

This example demonstrates a modular approach to organizing Gemini Enterprise resources across different environments.

## Structure

```
environments/
├── dev.tfvars        # Development environment variables
├── staging.tfvars    # Staging environment variables
├── prod.tfvars       # Production environment variables
└── main.tf           # Main configuration
```

## Usage

```bash
# Apply to development
terraform apply -var-file=dev.tfvars

# Apply to staging
terraform apply -var-file=staging.tfvars

# Apply to production
terraform apply -var-file=prod.tfvars
```

## Module Variables

- `environment`: Environment name (dev, staging, prod)
- `project_id`: Google Cloud project ID
- `gcs_bucket`: GCS bucket name for data
- `engine_id`: Unique identifier for the engine
- `display_name`: Display name for the engine
