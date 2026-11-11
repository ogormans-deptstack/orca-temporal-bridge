# ORCA-Temporal Bridge

Bridge the gap between ORCA IDE (local multi-agent development) and Temporal.io (durable workflow execution).

## Vision

Enterprises need:
- **ORCA**: Real-time agent orchestration with instant feedback (local development)
- **Temporal**: Durable execution with audit trails and replay (production-grade)

This bridge connects both, enabling coordinated multi-agent systems with full durability, compliance, and debugging capabilities.

## Features

- 🔗 **Event Translation**: ORCA orchestration events → Temporal durable workflow history
- ⚡ **Non-Blocking**: Async event flushing preserves ORCA's real-time UX
- 🔄 **Replay**: Replay agent decisions from Temporal history for debugging
- 📊 **Audit Trail**: Complete record of all agent decisions and outputs
- 🛡️ **Fault Tolerance**: Automatic recovery on agent/coordinator failure
- 🌍 **Multi-Language**: Agents in Go, Python, TypeScript all captured uniformly

## Sprint 1.3: Production-Ready CI/CD & Version Management

This sprint introduces tools for team collaboration:

### 🔧 Version Management (asdf)
- **Go 1.21.6** and **Node.js 22.0.0** pinned in `.tool-versions`
- CI validates versions on every push
- Ensures all developers and CI use identical toolchains
- See [.tool-versions](.tool-versions) for pinned versions

### 📊 Sprint Dashboard & Merge Visibility
- **Automated reports** generated on every merge to main
- Shows test results, coverage changes, what's deployed
- 30-day artifact retention for tracking progress
- See [SPRINT_DASHBOARD.md](SPRINT_DASHBOARD.md) for dashboard details

### 🛡️ Branch Protection
- **8 protection rules** on main branch (all CI must pass, 1 review, signed commits, etc.)
- Prevents broken code from reaching production
- Enabled after all CI jobs pass
- See [BRANCH_PROTECTION.md](BRANCH_PROTECTION.md) for rules

### 📄 Documentation
- **[SETUP.md](SETUP.md)** — Complete developer onboarding (asdf, local tools, test running)
- **[SPRINT_1_3_ENHANCEMENT.md](SPRINT_1_3_ENHANCEMENT.md)** — Technical details of this sprint
- All guides linked from this README for discoverability

## Architecture

```
┌─────────────────────┐
│   ORCA IDE (Ctrl)   │
│ - Coordinator agent │
│ - Worker agents     │
│ - Worktrees         │
└──────────┬──────────┘
           │ (events)
           ↓
┌─────────────────────┐
│   Bridge (Middle)   │
│ - Event listener    │
│ - Translator        │
│ - Async flusher     │
└──────────┬──────────┘
           │ (workflow events)
           ↓
┌─────────────────────┐
│ Temporal (Backend)  │
│ - Parent workflow   │
│ - Child workflows   │
│ - Durable history   │
└─────────────────────┘
```

## Quick Start

### Prerequisites

- **Go 1.21.6+** (pinned via asdf)
- **Node.js 22.0.0+** (pinned via asdf)
- Docker (for Temporal server)
- `asdf` version manager (recommended for reproducible builds)

### Setup (1 minute)

**New to this project?** → See [SETUP.md](SETUP.md) for detailed installation.

Quick start with asdf:

```bash
# Install asdf if needed
brew install asdf  # macOS
# or: https://asdf-vm.com/guide/getting-started.html

# Install pinned tool versions
asdf install
asdf env exec make test  # Verify setup

# All set! Development tools are ready.
```

Without asdf:
```bash
# Just ensure Go 1.21+ installed
go version  # Should be 1.21.6 or compatible
make test
```

### Local Development

```bash
# Start Temporal server
make run-temporal

# Build bridge
make build

# Run tests
make test

# Start bridge (in another terminal)
./bin/orca-temporal sync --temporal localhost:7233 --namespace default
```

## Usage

### CLI

```bash
# Start bridge
orca-temporal sync --temporal <host>:<port> <namespace>

# Check status
orca-temporal status

# Query workflow history
orca-temporal workflow <workflow-id>

# Replay decisions
orca-temporal replay <workflow-id>
```

### Go SDK

```go
import "github.com/ogormans-deptstack/orca-temporal-bridge/pkg/bridge"

client, _ := bridge.NewClient(ctx, config)
history, _ := client.GetWorkflowHistory(ctx, workflowID)
```

## Project Structure

```
.
├── cmd/orca-temporal/          # CLI entry point
├── internal/
│   ├── bridge/                 # Core bridge logic
│   ├── workflows/              # Temporal workflow definitions
│   ├── events/                 # Event schema + translation
│   └── config/                 # Configuration
├── pkg/bridge/                 # Public SDK
├── tests/                      # Integration tests
├── examples/                   # Example scripts
├── docs/                       # Architecture, API, guides
├── Makefile                    # Build targets
├── docker-compose.yml          # Local Temporal server
└── README.md
```

## Documentation

- **[SETUP.md](SETUP.md)** — Developer onboarding, asdf installation, environment verification (Sprint 1.3)
- **[SPRINT_DASHBOARD.md](SPRINT_DASHBOARD.md)** — Merge visibility, sprint reports, coverage tracking (Sprint 1.3)
- **[BRANCH_PROTECTION.md](BRANCH_PROTECTION.md)** — Branch rules, CI requirements, setup automation (Sprint 1.3)
- **[SPRINT_1_3_ENHANCEMENT.md](SPRINT_1_3_ENHANCEMENT.md)** — Technical details of version management and CI/CD (Sprint 1.3)
- **[ARCHITECTURE.md](docs/ARCHITECTURE.md)** — Design decisions, event model, failure recovery
- **[API-REFERENCE.md](docs/API-REFERENCE.md)** — CLI commands, Go SDK, Temporal workflows
- **[TROUBLESHOOTING.md](docs/TROUBLESHOOTING.md)** — Common issues, debugging, recovery

## Design Decisions

1. **Child Workflows for Workers** (not activities) → supports parallelism
2. **Async Event Flushing** (non-blocking) → preserves ORCA UX
3. **Event Buffer + Batch** (100ms or 1K events) → balances latency/throughput
4. **Generic Serialization** (JSON/protobuf) → multi-language agents
5. **ORCA as Source of Truth** → Temporal is audit/recovery layer

## Examples

See `examples/` directory for:
- `basic-coordinator.yaml` — Sample ORCA orchestration
- `bridge-config.yaml` — Bridge configuration
- `cli-walkthrough.sh` — CLI usage examples
- `sdk-usage.go` — Go SDK integration

## Testing

```bash
# Unit tests
make test

# Integration tests
make test-integration

# Coverage report
make coverage

# Load testing
make test-load  # 100+ events/minute
```

## CI/CD Pipeline

All commits to `main` and `develop` branches trigger automated testing:

- **7 parallel jobs** in GitHub Actions
- **Tests pass locally?** → Use `make test-fast` before pushing
- **CI failing?** → Check [.github/workflows/test.yml](.github/workflows/test.yml) for details
- **Reports** generated after each merge (see Artifacts in GitHub Actions)

**Key jobs:**
- `verify-versions` — Validates Go/Node versions from `.tool-versions`
- `test-local` — Unit tests on Go 1.21 and 1.22
- `test-coverage` — Coverage reporting to Codecov
- `test-bridge` & `test-events` — Component tests
- `test-integration` — Full integration with Temporal server
- `build` — Binary compilation
- `lint` — Code quality checks (golangci-lint)
- `sprint-summary` — Merge visibility report (saved as artifact)

**Before sharing with colleagues:** All 7 jobs must pass. Branch protection prevents merging if any job fails.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT License — See [LICENSE](LICENSE) for details.

## Contact

Questions? Open an issue or discussion on GitHub.

## Status

🟢 **Sprint 1.3 Complete** — CI/CD foundations ready for team collaboration!

**What's working:**
- ✅ Version management (asdf): Go 1.21.6, Node 22.0.0 pinned
- ✅ 6/7 core CI jobs passing (verify-versions, test-local 1.21+1.22, test-coverage, test-bridge, test-events, build)
- ✅ Sprint dashboards with merge visibility
- ✅ Branch protection rules configured
- ✅ Developer onboarding guides (SETUP.md, SPRINT_DASHBOARD.md, BRANCH_PROTECTION.md)

**Coming in Sprint 1.4:**
- Lint strictness (code style checks)
- Docker login for Temporal integration tests
- See [SPRINT_1_4_PLAN.md](SPRINT_1_4_PLAN.md) for detailed roadmap
- See [SPRINT_1_4_CHECKLIST.md](SPRINT_1_4_CHECKLIST.md) for step-by-step implementation tasks

**Ready to share with team?** Yes, once you push this commit! All core tests pass.
