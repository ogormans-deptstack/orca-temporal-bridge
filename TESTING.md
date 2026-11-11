# Testing Guide (Sprint 1.3)

This guide covers local testing and GitHub CI/CD for the ORCA-Temporal bridge.

---

## Quick Start

### Run All Tests Locally (No Temporal Required)
```bash
make test
# Output: All tests passed ✓
```

Takes ~1 second. No external dependencies needed.

### Run With Coverage
```bash
make test-coverage
# Opens coverage.html in your browser
```

### Run Specific Test Suites
```bash
make test-bridge    # Bridge E2E tests only
make test-events    # Event translator tests only
make test-fast      # Quick smoke test
```

---

## Local Testing

### Prerequisites

**Minimum** (for unit tests):
- Go 1.21+
- Make

**Full** (for integration tests):
- Go 1.21+
- Docker
- Docker Compose
- Make

### Unit Tests (No Docker)

```bash
# Run all unit tests
make test

# Run specific package
go test -v ./internal/bridge/...
go test -v ./internal/events/...

# Run single test
go test -v ./internal/bridge -run TestBridge_Lifecycle

# Run with race detector
go test -race ./...
```

**Output**:
```
=== RUN   TestBridge_Lifecycle
--- PASS: TestBridge_Lifecycle (0.00s)
...
ok  	github.com/ogormans-deptstack/orca-temporal-bridge/internal/bridge	0.260s
ok  	github.com/ogormans-deptstack/orca-temporal-bridge/internal/events	0.201s
✓ All tests passed
```

**Speed**: <1 second for all tests

**Coverage**:
- Events: 98.6%
- Bridge: 40.0% (workflows stubs bring down average)
- Total: 64.4%

### Integration Tests (With Temporal)

```bash
# Start Temporal in background
make run-temporal
# Temporal running at localhost:7233
# Temporal UI at http://localhost:8080

# In another terminal, run integration tests
make test-with-temporal

# When done, stop Temporal
make stop-temporal
```

**Manual integration testing**:
```bash
# Terminal 1: Start Temporal
make run-temporal

# Terminal 2: Run tests
go test -v -tags=integration -race ./internal/bridge ./internal/events

# Terminal 1: Stop when done
make stop-temporal
```

### Coverage Report

```bash
# Generate HTML coverage report
make test-coverage

# View in browser
open coverage.html

# View summary in terminal
go tool cover -func=coverage.out | tail -20
```

**Example output**:
```
github.com/ogormans-deptstack/orca-temporal-bridge/internal/bridge/bridge.go:14:		NewBridge			100.0%
github.com/ogormans-deptstack/orca-temporal-bridge/internal/bridge/bridge.go:40:		Start				85.0%
github.com/ogormans-deptstack/orca-temporal-bridge/internal/bridge/listener.go:11:		NewEventListener		100.0%
...
total:												(statements)		64.4%
```

---

## Quality Checks

### Format Code
```bash
make fmt
# Formats all Go files to standard

# Check if formatted (without modifying)
go fmt ./...
```

### Run Vet
```bash
make vet
# Checks for suspicious constructs

# Run vet + fmt before committing
make build
```

### Run Lint
```bash
make lint
# Requires: brew install golangci-lint

# Install if missing:
brew install golangci-lint

# Lint specific package:
golangci-lint run ./internal/bridge/...
```

---

## GitHub CI/CD

### Overview

The CI pipeline (`.github/workflows/test.yml`) runs on every push and PR:

| Job | Purpose | Requirements | Duration |
|-----|---------|--------------|----------|
| **test-local** | Fast unit tests, multiple Go versions | None | ~10s each |
| **test-coverage** | Full coverage report with upload to Codecov | None | ~20s |
| **test-bridge** | Bridge E2E tests | None | ~10s |
| **test-events** | Event translator tests | None | ~10s |
| **test-integration** | Integration tests with real Temporal | Docker service | ~30s |
| **build** | Build binary for Linux | None | ~10s |
| **lint** | golangci-lint check | None | ~30s |
| **test-summary** | Aggregate all job results | All above | ~5s |

### Workflow Trigger

Runs automatically on:
- Push to `main` or `develop` branches
- Pull requests against `main` or `develop`

### View Results

**On GitHub**:
1. Go to repository → Actions
2. Click workflow run
3. See job results and logs

**Coverage**:
1. Job: "Tests + Coverage"
2. Codecov badge shows coverage %
3. Click badge to see detailed breakdown

### Required Checks

All PR merges require:
- ✅ test-local
- ✅ test-coverage
- ✅ test-bridge
- ✅ test-events
- ✅ test-integration
- ✅ build
- ✅ test-summary

Lint is **optional** (doesn't block merge).

---

## Test Examples

### Bridge Configuration Test
```go
func TestBridge_Configuration(t *testing.T) {
    ctx := context.Background()
    mockProvider := NewMockEventProvider()
    defer mockProvider.Close()

    config := Config{
        TemporalHost:  "localhost",
        TemporalPort:  7233,
        Namespace:     "default",
        TaskQueue:     "test-queue",
        EventProvider: mockProvider,
    }

    bridge, err := NewBridge(ctx, config)
    require.NoError(t, err)
    assert.NotNil(t, bridge)
}
```

### Event Emission Test
```go
func TestMockEventProvider_CoordinatorStarted(t *testing.T) {
    provider := NewMockEventProvider()
    defer provider.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()

    eventStream, err := provider.Subscribe(ctx)
    require.NoError(t, err)

    // Emit test event
    err = provider.EmitCoordinatorStarted("wt-1", "coord-1", map[string]string{"version": "1.0"})
    require.NoError(t, err)

    // Verify received
    select {
    case event := <-eventStream:
        require.NotNil(t, event)
        assert.Equal(t, EventCoordinatorStarted, event.EventType)
    case <-ctx.Done():
        t.Fatal("timeout waiting for event")
    }
}
```

---

## Debugging Tests

### Run Single Test
```bash
go test -v ./internal/bridge -run TestBridge_Lifecycle
```

### Run With Verbose Output
```bash
go test -v -race ./internal/bridge/...
```

### Run With Debugging
```bash
# Print log statements
go test -v ./internal/bridge -run TestName -logger=verbose

# Keep temp files
go test -v ./internal/bridge -work
```

### Integration Test Debugging

```bash
# Start Temporal
make run-temporal

# Run single integration test
TEMPORAL_HOST=localhost:7233 \
TEMPORAL_NAMESPACE=default \
TEMPORAL_TASK_QUEUE=orca-bridge \
go test -v ./internal/bridge -run TestEventListener_Subscription

# View Temporal UI at http://localhost:8080

# Stop Temporal
make stop-temporal
```

---

## Continuous Integration Best Practices

### Before Committing
```bash
make build   # Formats + vets + builds
make test    # All tests pass
```

### Before Pushing
```bash
git status                 # Check what you're pushing
git diff                   # Review changes
make test-coverage         # Final coverage check
```

### On PR Creation
1. GitHub Actions starts automatically
2. All CI jobs run in parallel
3. Results visible in PR checks
4. Click "Details" to see logs

### If Tests Fail

**Locally:**
```bash
# Reproduce the failure
make test

# Debug with verbose output
go test -v ./internal/bridge -run FailingTest

# Check if it's a formatting issue
make fmt

# Check for race conditions
go test -race ./...
```

**On GitHub:**
1. Go to Actions → failing workflow
2. Click job → see error log
3. Copy error to local terminal
4. Reproduce and fix locally
5. Push fix (CI runs automatically again)

---

## Test Organization

### Test Packages

```
internal/
├── bridge/
│   ├── bridge.go              # Core bridge
│   ├── listener.go            # Event listener
│   ├── flusher.go             # Event flusher
│   ├── mock_provider.go       # Mock event source
│   └── bridge_test.go         # 13 E2E tests
│
└── events/
    ├── translator.go          # Event translation
    ├── types.go               # Event types
    └── translator_test.go     # 46 unit tests
```

### Test Naming

- **TestPackage_Feature** — high-level test (e.g., `TestBridge_Lifecycle`)
- **TestPackage_Feature_Variant** — variant test (e.g., `TestMockEventProvider_MultipleEvents`)

### Running Test Groups

```bash
# Run all Bridge tests
go test ./internal/bridge/... -v

# Run all Event tests
go test ./internal/events/... -v

# Run specific test
go test ./internal/bridge -run TestBridge_Configuration

# Run tests matching pattern
go test ./... -run "Bridge.*" -v
```

---

## Performance Benchmarks

### Current Metrics

| Component | Metric | Value |
|-----------|--------|-------|
| Unit tests | Time | <1s |
| Test count | Total | 59 |
| Coverage | Overall | 64.4% |
| Coverage | Events | 98.6% |
| Build time | Binary | ~2s |
| Integration tests | Time (with Temporal) | ~30s |

### Expected Performance After Sprint 1.4

- Add 10-15 workflow tests → total ~75 tests
- Add 5-10 SDK client tests → total ~90 tests
- Coverage target: 80%+ (currently 64.4% due to workflow stubs)

---

## Troubleshooting

### Tests Pass Locally But Fail on GitHub

**Cause**: Usually environment or Go version difference.

**Solution**:
```bash
# Test with same Go version as CI (1.21, 1.22)
go version

# Upgrade if needed
brew upgrade golang

# Run with race detector (CI uses -race)
go test -race ./...
```

### Temporal Service Won't Start

```bash
# Check if Docker is running
docker ps

# Check Temporal logs
docker-compose logs temporal

# Restart Docker
open -a Docker

# Clean up old containers
docker-compose down
docker system prune -a

# Try again
make run-temporal
```

### Coverage Report Won't Open

```bash
# Check if file exists
ls -la coverage.html

# Open manually
open coverage.html

# Or view in terminal
go tool cover -func=coverage.out
```

### golangci-lint Errors

```bash
# Install latest version
brew install golangci-lint

# Update
brew upgrade golangci-lint

# Run with verbose output
golangci-lint run ./... -v
```

---

## Next Steps (Sprint 1.4)

- Add workflow implementation tests (CoordinatorWorkflow, WorkerWorkflow)
- Add integration tests for full event pipeline
- Target 80%+ coverage
- Add performance benchmarks (`BenchmarkBridge*`)
- Add load testing scenarios

---

## Resources

- **Makefile**: `./Makefile` — Build targets
- **GitHub Actions**: `./.github/workflows/test.yml` — CI configuration
- **Project README**: `./README.md` — Architecture overview
- **Sprint Docs**: `./SPRINT_1_3.md` — Sprint 1.3 summary
- **Temporal Docs**: https://docs.temporal.io/
- **Go Testing**: https://golang.org/pkg/testing/
