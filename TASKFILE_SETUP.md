# Taskfile Setup (Sprint 1.3 Enhancement)

This is a **separate piece of work** that complements the existing Makefile with a modern task runner.

## What's New

### Files Added

```
Taskfile.yml              Main orchestration file (100 lines)
tasks/
├── test.yml             Testing tasks (150 lines, 12 tasks)
├── build.yml            Build tasks (70 lines, 5 tasks)
├── quality.yml          Code quality (170 lines, 10 tasks)
├── temporal.yml         Infrastructure (180 lines, 11 tasks)
└── ci.yml               CI/CD pipeline (200 lines, 8 tasks)

TASKFILE.md              Complete guide (400+ lines)
```

### Why Taskfile + Make?

**Make** (`Makefile`):
- Simple, Unix-standard
- No dependency (included on macOS/Linux)
- Works for basic scripts
- Good for CI/CD heredocs

**Taskfile** (`Taskfile.yml`):
- Much cleaner YAML syntax
- Better task composition
- Native preconditions
- Parallel execution
- Cross-platform (Windows/Mac/Linux)
- Better error handling
- More readable for complex workflows

### Use Both!

```bash
# Either works — no conflict
make test              # Make target
task test              # Taskfile task
task dev               # Taskfile workflow
make test-coverage     # Make target
```

---

## Installation

### Option 1: Homebrew (macOS/Linux)
```bash
brew install go-task/tap/go-task
task --version
```

### Option 2: Direct download
```bash
sh -c '$(curl --location https://taskfile.dev/install.sh)' -- -d -b /usr/local/bin
```

### Option 3: Go install
```bash
go install github.com/go-task/task/v3/cmd/task@latest
```

See https://taskfile.dev/installation/ for more options.

---

## Quick Start

### Show all available tasks
```bash
task              # Lists with descriptions
task --list
task --list-all   # Includes subtasks
```

### Common workflows

```bash
# 1. Full dev cycle (format → build → test)
task dev

# 2. Quick test (< 1s)
task test:fast

# 3. Coverage report
task coverage

# 4. Code quality checks
task fmt           # Format
task vet           # Analyze
task lint          # Lint

# 5. Temporal operations
task temporal:start
task temporal:ui    # Open in browser
task temporal:logs
task temporal:stop

# 6. CI pipeline locally
task ci            # Full pipeline
```

---

## Task Structure

### Main Tasks (Taskfile.yml)

**Quick shortcuts** for common operations:

```yaml
task dev           # Format → Build → Test
task test          # All tests
task build         # Binary
task fmt           # Format code
task lint          # Linter
task ci            # CI pipeline
task temporal:up   # Start Temporal
```

### Subtasks (tasks/*.yml)

**Organized by concern**:

```
test.yml       → test:all, test:fast, test:bridge, test:events, test:coverage, test:with-temporal, test:watch
build.yml      → build:binary, build:docker, build:clean, build:info
quality.yml    → quality:fmt, quality:vet, quality:lint, quality:all, quality:quick, quality:tidy, quality:report
temporal.yml   → temporal:start, temporal:stop, temporal:logs, temporal:ui, temporal:status, temporal:health
ci.yml         → ci:local, ci:pr-checks, ci:release, ci:status, ci:report:html, ci:report:json
```

---

## Example Workflows

### Daily Development

```bash
# Morning: Setup
task setup              # Verify environment

# During coding: Quick feedback
task test:fast          # < 1 second

# Before commit
task dev                # Full validation
```

### Before Push

```bash
# Check everything
task ci:pr-checks       # All PR validation
```

### Release Preparation

```bash
# Prepare release
task ci:release VERSION=v0.1.0
```

### Testing with Temporal

```bash
# Terminal 1: Start Temporal
task temporal:start

# Terminal 2: Run tests
task test:with-temporal

# Terminal 1: Stop when done
task temporal:stop
```

---

## Advantages over Make

### 1. **Cleaner Syntax**

**Make** (cryptic):
```makefile
test:
	@echo "Running tests..."
	@go test -v -race ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Tests passed"
```

**Taskfile** (readable):
```yaml
test:
  desc: Run all tests
  cmds:
    - echo "Running tests..."
    - go test -v -race ./...
    - go tool cover -html=coverage.out -o coverage.html
    - echo "✓ Tests passed"
```

### 2. **Task Composition**

```yaml
# Dependencies run automatically
dev:
  desc: Full dev workflow
  deps:
    - quality:fmt
    - build:binary
    - test:all
```

### 3. **Preconditions**

```yaml
test:
  preconditions:
    - sh: docker-compose ps temporal | grep -q Up
      msg: "Temporal not running. Run: task temporal:start"
  cmds:
    - go test ./...
```

### 4. **Parallel Execution**

```bash
task -p build docker   # Runs both in parallel
```

### 5. **Better Output Formatting**

- Colored output
- Task progress
- Error messages
- Status indicators

---

## Comparison Table

| Feature | Make | Taskfile | Winner |
|---------|------|----------|--------|
| Syntax clarity | 2/10 | 9/10 | Taskfile |
| Task dependencies | 6/10 | 9/10 | Taskfile |
| Preconditions | 2/10 | 9/10 | Taskfile |
| Parallel execution | 6/10 | 9/10 | Taskfile |
| Cross-platform | 4/10 | 10/10 | Taskfile |
| Installation ease | 10/10 | 6/10 | Make |
| No external deps | 10/10 | 3/10 | Make |
| CI/CD support | 7/10 | 8/10 | Taskfile |

**Verdict**: Taskfile for local development, Make for CI/CD (no install needed)

---

## Both Systems Work

### Make (in CI/CD, no install needed)
```bash
make test
make build
make fmt
```

### Taskfile (local development, better UX)
```bash
task dev          # Full workflow
task test:watch   # Auto-test on changes
task ci:local     # Local CI simulation
```

### Mix and match
```bash
task test && make build     # Task for testing, make for build
make fmt && task test       # Make for CI, task for dev
```

---

## Testing the Taskfile

Without installing task, you can still validate the YAML:

```bash
# Check syntax
cd /Users/ogormans/workspace/orca-temporal-bridge
cat Taskfile.yml | head -20    # Preview syntax
ls -la tasks/                   # See all task files
```

After installing `task`:

```bash
# See all tasks
task --list

# Try a task
task quality:fmt

# See specific task
task -p test:all
```

---

## Local Testing vs GitHub CI

### Local (with Taskfile)
```bash
task dev                    # Quick
task test:watch             # Auto-test
task ci:local               # Full pipeline
task temporal:start && ... # Interactive
```

### GitHub Actions (Make + Shell)
```bash
make test                   # In CI, no task install needed
make build
go test -race ./...
```

**Why both?**
- Taskfile makes local development faster and cleaner
- Make keeps CI/CD simple (no external dependencies)
- They don't conflict — use what fits best

---

## Next Steps

### 1. Install Taskfile (optional, for local dev)
```bash
brew install go-task/tap/go-task
```

### 2. Explore tasks
```bash
task --list           # See all available tasks
task --help           # Show help
```

### 3. Try a task
```bash
task test:fast        # Quick smoke test
task fmt              # Format code
task dev              # Full dev workflow
```

### 4. Read the guide
```bash
# Comprehensive guide with examples
cat TASKFILE.md
```

---

## Files Added in Sprint 1.3

| File | Size | Purpose |
|------|------|---------|
| `Taskfile.yml` | ~100 lines | Main orchestration + shortcuts |
| `tasks/test.yml` | ~150 lines | 12 testing tasks |
| `tasks/build.yml` | ~70 lines | 5 build tasks |
| `tasks/quality.yml` | ~170 lines | 10 code quality tasks |
| `tasks/temporal.yml` | ~180 lines | 11 infrastructure tasks |
| `tasks/ci.yml` | ~200 lines | 8 CI/CD pipeline tasks |
| `TASKFILE.md` | ~400 lines | Complete guide + examples |

**Total**: ~1,270 lines of task definitions + documentation

---

## Architecture

### Task Organization

```
Taskfile.yml (main orchestrator)
│
├── Quick shortcuts (task dev, task test, task build)
│
└── Includes:
    ├── tasks/test.yml       → test:* (local, integration, watch, coverage)
    ├── tasks/build.yml      → build:* (binary, docker, clean)
    ├── tasks/quality.yml    → quality:* (fmt, vet, lint, tidy, report)
    ├── tasks/temporal.yml   → temporal:* (start, stop, logs, ui, health)
    └── tasks/ci.yml         → ci:* (local, pr-checks, release, report)
```

### Task Dependencies

```
task dev
├── quality:fmt
├── build:binary
│   ├── quality:fmt
│   └── quality:vet
└── test:all
    └── (internal/... tests)
```

### Parallel Execution

```bash
task -p build docker    # Both run in parallel
```

---

## Compatibility

### Existing Systems

- ✅ Makefile still works (100% compatible)
- ✅ GitHub Actions still works (uses Make)
- ✅ Docker Compose still works (unchanged)
- ✅ Test scripts still work (./scripts/test-local.sh)

### New Systems

- ✅ Taskfile adds modern alternatives
- ✅ No breaking changes
- ✅ Optional (not required)
- ✅ Can be added/removed anytime

---

## Troubleshooting

### Task not found
```bash
task --list           # See all tasks
task -a               # Show everything
```

### Precondition error
```bash
# Example: Temporal not running
task test:with-temporal
# Error: Temporal not running. Run: task temporal:start

task temporal:start
task test:with-temporal  # Try again
```

### No task executable
```bash
# Install task first
brew install go-task/tap/go-task

# Or use Make instead
make test
```

---

## Recommendations

### For Local Development (Interactive)
Use **Taskfile** for:
- Better UX (colored output, progress)
- `task test:watch` (auto-test on changes)
- `task dev` (one-command workflows)
- `task ci:local` (test like GitHub Actions)

### For CI/CD (GitHub Actions)
Use **Make** for:
- No external dependencies
- Works in any environment
- Simpler GitHub Actions configuration

### Rule of Thumb
- **Local**: `task` (if installed)
- **CI/CD**: `make` (no dependencies)
- **Both available**: Use what you prefer!

---

## Learning Resources

- **Taskfile Docs**: https://taskfile.dev/
- **Tutorial**: https://taskfile.dev/usage/
- **YAML syntax**: https://yaml.org/ (basic knowledge sufficient)

---

## Summary

**Taskfile is a separate enhancement to Sprint 1.3** — not a replacement, but an improvement on local development workflows.

- ✅ Modern YAML-based task runner
- ✅ Better UX than Make
- ✅ Complements existing Makefile
- ✅ Optional to install
- ✅ No breaking changes
- ✅ CI/CD still uses Make (no install needed)

**Use both!** They work side-by-side.

---

**Next**: Try `task --list` after installing taskfile, then experiment with `task dev` workflow!
