# Go 1.26 Linting Guide

Go 1.26 is bleeding edge. The linting ecosystem is still catching up. This guide explains how we handle linting in our Go 1.26 codebase.

## The Problem

- **golangci-lint v1.64.8** (latest v1.x): Built with Go 1.24, can't lint Go 1.26 code
- **golangci-lint v2.x**: Supports Go 1.26 but has different CLI flags
- **CI**: Disabled linting until Go 1.26 tooling stabilizes
- **Local dev**: Need a practical way to lint without breaking workflow

## Our Solution: Docker-Based Linting

We use **Docker** to run golangci-lint. This isolates the linter from the Go version mismatch.

### Quick Start

```bash
# Lint your code
make lint

# If you want local linting (may fail if golangci-lint is Go 1.24)
make lint-local
```

## Installation for Developers

### Option 1: Manual Linting (Recommended)

Run before each commit:
```bash
make fmt        # Format code
make vet        # Check with go vet
make lint       # Lint with Docker
make test       # Run tests
```

### Option 2: Automatic Pre-Commit Hooks

Install hooks to run linting automatically:
```bash
make hooks-install
```

Uninstall later with:
```bash
make hooks-uninstall
```

Bypass hooks if needed:
```bash
git commit --no-verify
```

## How It Works

### Docker-Based Linting (`make lint`)

```bash
docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:latest golangci-lint run
```

This:
1. Runs golangci-lint inside a Docker container
2. Mounts your repo as a volume
3. Uses the container's Go version (compatible with our code)
4. Returns results to your terminal

**Pros**:
- ✅ Works with Go 1.26
- ✅ No local installation needed
- ✅ Consistent across machines
- ✅ No version conflicts

**Cons**:
- ⚠️ Requires Docker running
- ⚠️ Slightly slower than local lint

### Local Linting (`make lint-local`)

Tries to run locally-installed golangci-lint:
```bash
golangci-lint run ./... --config=.golangci.yml
```

**Pros**:
- ✅ Fast
- ✅ No Docker dependency

**Cons**:
- ❌ May fail if golangci-lint was built with Go 1.24
- ❌ May show spurious errors with Go 1.26

## Troubleshooting

### "Docker not found"
Install Docker: `brew install docker`

### "golangci-lint not found" (local)
That's okay! Use Docker instead: `make lint`

### Linting fails locally but works with Docker
Your local golangci-lint is incompatible with Go 1.26. Use `make lint` instead.

### Pre-commit hook not running
Reinstall hooks:
```bash
make hooks-uninstall
make hooks-install
```

## Configuration

Linting configuration is in `.golangci.yml`:
- **Enabled linters**: errcheck, govet, staticcheck, misspell, unused, ineffassign, gocritic
- **Excluded**: vendor/ directory, test files for certain linters
- **Timeout**: 5 minutes

Edit `.golangci.yml` to adjust linting rules for the team.

## Future: When Go 1.26 Tooling Matures

Once golangci-lint releases a v1.x version built with Go 1.26+:
1. Update `.github/workflows/test.yml` to enable linting
2. Specify the new version in the workflow
3. Developers can optionally switch back to local linting
4. Keep Docker approach as fallback

## References

- [golangci-lint documentation](https://golangci-lint.run/)
- [.golangci.yml configuration](../.golangci.yml)
- [CONTRIBUTING.md](../CONTRIBUTING.md) - Full development guide
