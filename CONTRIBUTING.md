# Contributing to ORCA-Temporal Bridge

Thank you for your interest in contributing! This document provides guidelines for contributing to the project.

## Getting Started

### Prerequisites

- Go 1.26 (see `.tool-versions`)
- Docker (for local Temporal server and linting with Go 1.26)
- Git

### Setup Development Environment

```bash
# Clone the repository
git clone https://github.com/ogormans-deptstack/orca-temporal-bridge.git
cd orca-temporal-bridge

# Install dependencies
go mod download

# (Optional) Install git pre-commit hooks for automatic checks
make hooks-install

# Start local Temporal server
make run-temporal

# Run tests to verify setup
make test
```

## Development Workflow

### Code Organization

- `cmd/` — Command-line interface entry points
- `internal/bridge/` — Core bridge logic
- `internal/workflows/` — Temporal workflow definitions
- `internal/events/` — Event schema and translation
- `pkg/bridge/` — Public SDK
- `tests/` — Integration and load tests
- `examples/` — Example scripts and configurations
- `docs/` — Architecture and API documentation

### Writing Code

1. **Follow Go conventions**:
   - Use `gofmt` for formatting
   - Use `go vet` to check for common errors
   - Write clear, idiomatic Go

2. **Write tests**:
   - TDD: Write tests before implementation
   - Target 80%+ code coverage
   - Test both happy path and error cases

3. **Document your work**:
   - Add comments to exported functions and types
   - Update relevant documentation files
   - Include examples for new features

### Testing

```bash
# Run all tests
make test

# Run specific test
go test -v -run TestName ./...

# Check coverage
make coverage

# Run integration tests (requires Temporal)
make test-with-temporal

# Quality checks before committing
make fmt              # Format code
make vet              # Run go vet
make lint             # Lint code (Docker-based, Go 1.26 compatible)
make lint-local       # Lint locally if compatible version installed
```

### Code Review

1. Create a feature branch: `git checkout -b feature/my-feature`
2. Make your changes and write tests
3. **Before committing**: Run quality checks:
   ```bash
   make fmt            # Format code
   make vet            # Check with go vet
   make lint           # Lint with Docker (works with Go 1.26)
   make test           # Run tests
   ```
4. Push and open a pull request
5. Address review feedback

**Note on Go 1.26 Linting**: The linting toolchain for Go 1.26 is still being adopted by the ecosystem. 
We use Docker-based golangci-lint (`make lint`) for linting to ensure compatibility. 
You can optionally install pre-commit hooks to automate these checks:
```bash
make hooks-install   # Install git pre-commit hooks
make hooks-uninstall # Remove git pre-commit hooks
```

## Architecture Decisions

Before making significant architectural changes:

1. **Read the design documentation**: See `docs/ARCHITECTURE.md`
2. **Understand the 5 key design decisions**:
   - Child workflows for workers (not activities)
   - Async event flushing (non-blocking)
   - Event buffer + batch strategy
   - Generic serialization (multi-language)
   - ORCA as source of truth
3. **Propose changes via issue** before implementing

## Go 1.26 Linting Strategy

**Context**: Go 1.26 is bleeding edge. The linting toolchain (golangci-lint) is still being built/tested against Go 1.26.

**Solution**: We use Docker-based linting to ensure compatibility:
- **`make lint`** — Runs golangci-lint in Docker (works with any Go version)
- **`make lint-local`** — Runs locally if you have a compatible version installed
- **`.githooks/pre-commit`** — Optional pre-commit hook that runs linting before commits

**For developers**:
```bash
# Run linting before pushing
make lint

# Or install pre-commit hooks (automatic checks)
make hooks-install

# Bypass hooks if needed (not recommended)
git commit --no-verify
```

**For team**: All linting currently runs locally. CI does not lint (due to Go 1.26 tooling gap). 
Re-enable CI linting when Go 1.26 golangci-lint is available.

## Reporting Issues

- Use GitHub Issues to report bugs and request features
- Provide reproducible examples
- Include relevant logs and environment details

## Style Guide

### Naming Conventions

- Packages: lowercase, concise (`events`, `workflows`, `bridge`)
- Functions: CamelCase, descriptive (`TranslateOrcaEvent`, `GetWorkflowHistory`)
- Constants: SCREAMING_SNAKE_CASE (`EventCoordinatorStarted`)
- Errors: Prefix with "Err" (`ErrBridgeStopped`, `ErrNotImplemented`)

### Comments

- Public functions and types must have comments
- Explain the "why", not just the "what"
- Include examples for complex logic

### Error Handling

- Never ignore errors silently
- Use `fmt.Errorf` with context
- Create named error variables for common errors

## Pull Request Process

1. Update documentation if needed
2. Add/update tests
3. Ensure `make test` and `make lint` pass
4. Keep commits atomic and well-messaged
5. Keep PRs focused on a single concern

## Release Process

This project follows semantic versioning (MAJOR.MINOR.PATCH):

- **MAJOR**: Breaking API changes
- **MINOR**: New features (backward-compatible)
- **PATCH**: Bug fixes

Releases are tagged on main branch with version number (e.g., `v0.1.0`).

## Questions?

- Open a GitHub Discussion
- Check existing documentation in `docs/`
- Review examples in `examples/`

Thank you for contributing!
