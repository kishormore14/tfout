# Contributing to `tfout`

Thank you for your interest in contributing to `tfout`! We welcome community contributions, bug reports, feature requests, and pull requests.

## Development Requirements

- **Go**: 1.21 or higher.
- **HashiCorp Terraform**: Optional integration testing dependency (not required to work on the core formatter).

## Project Guidelines

1. **Language & Scope**:
   - Written strictly in **Go**.
   - Scope is restricted to HashiCorp Terraform integration and general terminal formatting.
   - Do not add dependencies or integrations for OpenTofu.

2. **Zero Dependency Principle**:
   - Keep external Go module dependencies minimal (preferably standard library only).

3. **Security Standards**:
   - All input text must be passed through terminal control sequence sanitizers (`pkg/sanitize`) to prevent ANSI escape code injection.
   - Outputs marked sensitive by Terraform must be redacted by default.
   - Subprocess execution must use direct parameter arrays without passing input to shell interpreters.

4. **Testing**:
   - All new features and bug fixes must include unit test coverage.
   - Run tests before opening a Pull Request:
     ```bash
     go test ./... -v
     go vet ./...
     ```

## Code Style

Follow standard Go idioms:
- `gofmt` code formatting.
- Clear error handling with descriptive error messages.
- Typed errors where appropriate.
