# Taskfile Guide (Sprint 1.3)

Modern task runner for ORCA-Temporal Bridge using [Taskfile](https://taskfile.dev/).

## Installation

### macOS
```bash
brew install go-task/tap/go-task
```

### Linux
```bash
sh -c '$(curl --location https://taskfile.dev/install.sh)' -- -d -b /usr/local/bin
```

### Or download from https://taskfile.dev/installation/

## Quick Start

### Show all tasks
```bash
task              # Lists all available tasks
task --list       # Same, with descriptions
```

### Common workflows

```bash
# Full dev cycle: format → build → test
task dev

# Run tests
task test                 # All tests (no Temporal)
task test:fast            # Quick smoke test
task test:coverage        # Coverage report
task test:watch           # Watch mode (requires entr)

# Build
task build                # Build binary

# Code quality
task fmt                  # Format code
task vet                  # Run go vet
task lint                 # Run linter

# Temporal
task temporal:start       # Start Temporal server
task temporal:stop        # Stop Temporal server
task temporal:ui          # Open Temporal UI
```

---

## Task Organization

Taskfile is split into 5 files for clarity:

```
Taskfile.yml          Main file (orchestration + quick shortcuts)
└── tasks/
    ├── test.yml      Testing tasks (local, integration, watch, coverage)
    ├── build.yml     Build tasks (binary, docker, clean)
    ├── quality.yml   Code quality (fmt, vet, lint, imports, tidy)
    ├── temporal.yml  Infrastructure (start, stop, logs, UI, health)
    └── ci.yml        CI/CD pipeline (local CI, PR checks, release)
```

---

## Usage Examples

### Development Workflow

#### Start of day: Verify environment
```bash
task setup
# Shows: Go version, Docker version, test setup
```

#### During development: Quick feedback
```bash
task test:fast
# Runs in < 1 second, no coverage overhead
```

#### Before commit: Full validation
```bash
task dev
# Formats code → builds binary → runs all tests
```

### Testing

#### Run specific test
```bash
task test:run TEST=TestBridge_Lifecycle
# Run single test by name
```

#### Coverage report
```bash
task coverage
# Generates coverage.html and opens in browser
```

#### Watch mode (auto-test on changes)
```bash
task test:watch
# Re-runs tests when .go files change
# Requires: brew install entr
```

#### Integration tests with Temporal
```bash
task temporal:start
task test:with-temporal
task temporal:stop

# Or in one command with automatic cleanup:
task ci:local   # Full pipeline including integration tests
```

### Build & Release

#### Build binary
```bash
task build
# Output: ./bin/orca-temporal
```

#### Build Docker image
```bash
task docker
# Image: orca-temporal-bridge:latest
```

#### Prepare release
```bash
task ci:release VERSION=v0.1.0
# Runs all checks + creates tag
```

### Code Quality

#### Format code
```bash
task fmt
# Apply gofmt to all files
```

#### Check formatting (no modifications)
```bash
task fmt:check
# Exits with error if changes needed (for CI)
```

#### Run all quality checks
```bash
task quality:all
# fmt:check + vet + lint
```

#### Quick quality check (no lint)
```bash
task quality:quick
# fmt:check + vet only (faster)
```

### Infrastructure

#### Manage Temporal
```bash
task temporal:start      # Start server
task temporal:stop       # Stop server
task temporal:restart    # Restart
task temporal:status     # Check status
task temporal:logs       # View logs
task temporal:ui         # Open UI
task temporal:health     # Health check
```

### CI/CD

#### Run full local CI pipeline
```bash
task ci:local
# Everything GitHub Actions runs:
# - Quality checks
# - All tests
# - Coverage report
# - Build
```

#### PR validation
```bash
task ci:pr-checks
# Checks that must pass for PR:
# - Format validation
# - Vet analysis
# - All tests
# - Coverage check
# - Build verification
```

#### CI status dashboard
```bash
task ci:status
# Shows project health overview
```

---

## Task Reference

### Test Tasks (`tasks/test.yml`)

| Task | Description | Time |
|------|-------------|------|
| `test:all` | All tests (no Temporal) | ~1s |
| `test:fast` | Smoke test (bridge + events) | <1s |
| `test:bridge` | Bridge E2E tests only | ~0.2s |
| `test:events` | Event translator tests | ~0.2s |
| `test:coverage` | Full coverage report | ~2s + generates HTML |
| `test:with-temporal` | Integration tests (requires Temporal) | ~10s |
| `test:watch` | Auto-test on file changes | continuous |
| `test:run TEST=Name` | Run specific test | varies |
| `test:bench` | Run benchmarks | varies |

### Build Tasks (`tasks/build.yml`)

| Task | Description |
|------|-------------|
| `build:binary` | Build binary (./bin/orca-temporal) |
| `build:docker` | Build Docker image |
| `build:docker:run` | Build and run image |
| `build:clean` | Remove build artifacts |
| `build:info` | Show build information |

### Quality Tasks (`tasks/quality.yml`)

| Task | Description | Time |
|------|-------------|------|
| `quality:fmt` | Format code | <1s |
| `quality:fmt:check` | Check formatting (no changes) | <1s |
| `quality:vet` | Run go vet | <1s |
| `quality:lint` | Run golangci-lint | ~30s |
| `quality:all` | fmt:check + vet + lint | ~31s |
| `quality:quick` | fmt:check + vet (no lint) | ~1s |
| `quality:imports` | Organize imports | <1s |
| `quality:tidy` | Tidy go.mod/go.sum | <1s |
| `quality:deps` | Check outdated dependencies | <1s |
| `quality:report` | Code quality report | ~1s |

### Temporal Tasks (`tasks/temporal.yml`)

| Task | Description |
|------|-------------|
| `temporal:start` | Start Temporal server |
| `temporal:stop` | Stop Temporal server |
| `temporal:restart` | Restart Temporal |
| `temporal:status` | Check status |
| `temporal:logs` | View logs (tail -f) |
| `temporal:logs:last` | Last 50 lines |
| `temporal:ui` | Open UI in browser |
| `temporal:clean` | Delete volumes (⚠️ data loss) |
| `temporal:health` | Health check |
| `temporal:setup` | Complete setup |
| `temporal:version` | Show Temporal version |

### CI Tasks (`tasks/ci.yml`)

| Task | Description |
|------|-------------|
| `ci:local` | Full CI pipeline (GitHub Actions equivalent) |
| `ci:pr-checks` | PR validation checks |
| `ci:release VERSION=v0.1.0` | Prepare release |
| `ci:github-actions` | Simulate GitHub Actions |
| `ci:status` | Status dashboard |
| `ci:report:html` | HTML test report |
| `ci:report:json` | JSON test report |
| `ci:clean` | Clean CI artifacts |

### Shortcuts (in main Taskfile.yml)

| Task | Shortcut For |
|------|--------------|
| `task dev` | `quality:fmt` → `build:binary` → `test:all` |
| `task test` | `test:all` |
| `task test:fast` | `test:fast` |
| `task coverage` | `test:coverage` |
| `task build` | `build:binary` |
| `task fmt` | `quality:fmt` |
| `task lint` | `quality:lint` |
| `task vet` | `quality:vet` |
| `task clean` | `build:clean` |

---

## Environment Variables

### Global (set in Taskfile.yml)
```yaml
GO_VERSION: '1.21'
COVERAGE_THRESHOLD: '80'
TEMPORAL_HOST: localhost
TEMPORAL_PORT: 7233
TEMPORAL_NAMESPACE: default
TEMPORAL_TASK_QUEUE: orca-bridge
```

### Override in shell
```bash
# Use different Go version (if installed)
GO_VERSION=1.22 task test

# Use different Temporal host
TEMPORAL_HOST=remote.example.com task test:with-temporal

# Skip linter in quality check
task quality:quick  # Only fmt + vet, no lint
```

---

## Advanced Usage

### Task Dependencies

Tasks can depend on other tasks:

```yaml
task:
  deps:
    - quality:fmt
    - quality:vet
  cmds:
    - go build ./...
```

Taskfile automatically runs dependencies first.

### Conditional Tasks

Use preconditions to skip tasks:

```yaml
task:
  preconditions:
    - sh: docker-compose ps temporal | grep -q Up
      msg: "Temporal not running. Run: task temporal:start"
  cmds:
    - go test -tags=integration ./...
```

### Task Parameters

Pass variables to tasks:

```bash
task test:run TEST=TestName
task ci:release VERSION=v0.1.0
```

### Parallel Execution

Run independent tasks in parallel:

```bash
task -p build docker   # Run build:binary and build:docker in parallel
```

---

## Troubleshooting

### Task not found
```bash
task --list   # Show all available tasks
task -a       # Alias for --list
```

### Precondition failed
```bash
# Example: Temporal not running
task test:with-temporal
# Error: Temporal not running. Run: task temporal:start

task temporal:start
task test:with-temporal  # Now works
```

### Verbose output
```bash
task -v test           # Verbose mode
task --verbose test    # Same
```

### Dry run (preview without executing)
```bash
task -d build          # Shows what would run
task --dry build       # Same
```

### Continue on error
```bash
task -i test           # Continue even if test fails
task --ignore-status test
```

---

## Comparison: Taskfile vs Make

| Feature | Taskfile | Make |
|---------|----------|------|
| **Language** | YAML | Makefile syntax |
| **Readability** | Very clear | Can be cryptic |
| **Cross-platform** | ✓ Windows/Mac/Linux | ✗ Needs GNU make |
| **Task dependencies** | ✓ Built-in | ✓ Built-in |
| **Preconditions** | ✓ Native | ✗ Workaround with shell |
| **Parallel tasks** | ✓ `-p` flag | ✓ `-j` flag |
| **Variables** | ✓ YAML-based | ~ Make variables |
| **Output formatting** | ✓ Pretty | ~ Basic |
| **Composition** | ✓ includes | ✗ Limited |
| **Error handling** | ✓ Built-in | ~ Manual |

**Both are supported** — use whichever you prefer!

```bash
# Using Make
make test
make build
make fmt

# Using Taskfile (equivalent)
task test
task build
task fmt
```

---

## Migration Guide (from Make to Taskfile)

### Make targets → Taskfile tasks

```bash
# Make
make build        →  task build        # or task build:binary
make test         →  task test         # or task test:all
make fmt          →  task fmt          # or task quality:fmt
make vet          →  task vet          # or task quality:vet
make lint         →  task lint         # or task quality:lint
make run-temporal →  task temporal:up  # or task temporal:start
make clean        →  task clean        # or task build:clean
make coverage     →  task coverage     # or task test:coverage
```

### Both work! Use what you prefer

```bash
# Mix and match
task test && make build     # Task for testing, make for build
make fmt && task test       # Make for formatting, task for testing
```

---

## Next Steps

1. **Install Taskfile**: `brew install go-task/tap/go-task`
2. **Explore tasks**: `task --list`
3. **Try it out**: `task dev` (full workflow)
4. **Read docs**: https://taskfile.dev/usage/

---

## Tips & Tricks

### Alias for less typing
```bash
alias t=task

# Then use:
t test
t build
t dev
```

### View specific task details
```bash
task -p test:all       # Shows task definition
```

### Watch mode for development
```bash
task test:watch        # Auto-test on file changes
```

### Temporal + tests in one command
```bash
task temporal:start && task test:with-temporal && task temporal:stop
```

### Generate artifacts for CI
```bash
task ci:report:json    # JSON for artifact collection
task ci:report:html    # HTML for reports
```

---

**Enjoy the improved developer experience with Taskfile!** 🚀
