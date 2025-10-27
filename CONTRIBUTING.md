# Contributing to terraform-provider-gemctl

Thank you for considering contributing to terraform-provider-gemctl! This document provides guidelines for contributing to the project.

## Getting Started

1. **Fork the repository**
   ```bash
   git clone https://github.com/vb140772/terraform-provider-gemctl.git
   cd terraform-provider-gemctl
   ```

2. **Set up development environment**
   - Go 1.21 or later
   - Terraform CLI
   - Google Cloud SDK
   - Terraform Plugin Docs (for generating documentation)

3. **Build the provider**
   ```bash
   go build .
   ```

4. **Run tests**
   ```bash
   go test ./...
   ```

5. **Configure dev overrides**
   Create `~/.terraformrc`:
   ```
   provider_installation {
     dev_overrides {
       "vb140772/gemctl" = "/path/to/terraform-provider-gemctl"
     }
   }
   ```

## Development Guidelines

### Code Style
- Follow Go standard conventions
- Use `gofmt` and `golint` for code formatting
- Add comments to exported functions and types
- Write descriptive commit messages

### Testing
- Write tests for new features
- Run `go test ./...` before submitting PR
- Test with real Google Cloud resources when possible

### Documentation
- Update README.md when adding features
- Run `go generate ./...` to regenerate docs
- Add examples for new resources/data sources
- Document breaking changes in CHANGELOG.md

### Commit Messages
Follow conventional commits format:
- `feat:` for new features
- `fix:` for bug fixes
- `docs:` for documentation changes
- `test:` for test changes
- `chore:` for maintenance tasks

## Submitting Changes

1. **Create a branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make changes**
   - Write code
   - Add tests
   - Update documentation
   - Run `go fmt ./...`

3. **Commit your changes**
   ```bash
   git commit -m "feat: add new feature"
   ```

4. **Push and create PR**
   ```bash
   git push origin feature/your-feature-name
   ```

## Project Structure

```
terraform-provider-gemctl/
├── internal/
│   ├── client/       # Google API client
│   └── provider/     # Terraform provider implementation
├── examples/         # Example configurations
├── docs/             # Generated documentation
├── main.go           # Provider entry point
└── README.md         # Project documentation
```

## Resources

- [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin)
- [Terraform Registry Documentation](https://developer.hashicorp.com/terraform/registry/providers/docs)
- [Google Discovery Engine API](https://cloud.google.com/generative-ai-app-builder/docs)

## Questions?

Feel free to open an issue for questions or clarifications.
