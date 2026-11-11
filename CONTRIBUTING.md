# Contributing to ORCA-Temporal Bridge

Thank you for your interest in contributing! This document provides guidelines for contributing to the project.

## Getting Started

### Prerequisites

- Go 1.21 or later
- Docker (for local Temporal server)
- Git

### Setup Development Environment

```bash
# Clone the repository
git clone https://github.com/ogormans-deptstack/orca-temporal-bridge.git
cd orca-temporal-bridge

# Install dependencies
go mod download

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
make test-integration

# Run load tests
make test-load
```

### Code Review

1. Create a feature branch: `git checkout -b feature/my-feature`
2. Make your changes and write tests
3. Format and lint: `make fmt && make lint`
4. Push and open a pull request
5. Address review feedback

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
