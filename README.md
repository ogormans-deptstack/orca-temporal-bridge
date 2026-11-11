# ORCA-Temporal Bridge

Bridge the gap between ORCA IDE (local multi-agent development) and Temporal.io (durable workflow execution).

## Vision

Enterprises need coordinated multi-agent systems with:
- **Real-time development**: ORCA IDE with instant feedback
- **Production durability**: Temporal.io with full audit trails, replay, and fault tolerance

This bridge connects both, enabling teams to develop locally and deploy durably.

## Features

- 🔗 **Event Translation**: ORCA events → Temporal durable workflow history
- ⚡ **Non-Blocking**: Async event flushing preserves ORCA's real-time UX
- 🔄 **Replay**: Debug agent decisions from Temporal history
- 📊 **Audit Trail**: Complete record of all agent decisions
- 🛡️ **Fault Tolerance**: Automatic recovery on agent/coordinator failure
- 🌍 **Multi-Language**: Agents in Go, Python, TypeScript captured uniformly

## Status: Ready for Team Collaboration ✅

### Latest: Sprint 1.4 Complete
**Go 1.26 Foundation + Temporal Integration + Developer Workflow**

| Component | Status | Details |
|-----------|--------|---------|
| **Go Version** | ✅ 1.26 | Modern, future-proof toolchain in `.tool-versions` |
| **Temporal Integration** | ✅ Working | `temporal server start-dev`, no Docker credentials needed |
| **Code Quality** | ✅ Clean | 0 lint issues, linting via Docker (make lint) |
| **Local Dev** | ✅ Ready | Pre-commit hooks available (make hooks-install) |
| **Integration Tests** | ✅ Passing | All tests passing locally and in CI |
| **CI/CD** | ✅ 7/7 jobs | All core jobs passing, branch protection ready |

### What Sprints 1.3 & 1.4 Delivered

**Sprint 1.3: Foundation**
- Version management (Go 1.26, Node.js 22 pinned via asdf)
- 7-job CI/CD pipeline with automated testing
- Sprint dashboards with merge visibility
- Complete developer onboarding guides

**Sprint 1.4: Developer Experience**
- Go 1.26 adoption (ecosystem tooling gaps managed)
- Docker-based linting (works with Go 1.26)
- Pre-commit hook infrastructure
- Temporal integration (simpler, no secrets needed)
- Comprehensive team documentation

**See [SPRINT_1_4_RETROSPECTIVE.md](SPRINT_1_4_RETROSPECTIVE.md)** for detailed plan vs reality analysis.

## Quick Start

### Prerequisites
- Go 1.26 (pinned via asdf in `.tool-versions`)
- Node.js 22.0.0 (pinned via asdf in `.tool-versions`)
- Docker (optional, for local Temporal UI)

### Setup (< 1 minute)

```bash
# Install asdf if needed
brew install asdf  # macOS
# Or: https://asdf-vm.com/guide/getting-started.html

# Install pinned versions
asdf install

# Verify setup
make test

# Enable auto-linting (recommended)
make hooks-install
```

### Local Development

```bash
# Start Temporal (temporal server start-dev)
make run-temporal

# Run tests
make test

# Lint code
make lint

# Build binary
make build
```

See [SETUP.md](SETUP.md) for detailed onboarding.

## Architecture

```
ORCA (Local Dev)           Bridge (Sync)           Temporal (Durable)
├─ Coordinator agent       ├─ Event listener       ├─ Parent workflow
├─ Worker agents      ───► ├─ Translator      ───► ├─ Child workflows
└─ Real-time events        └─ Async flusher        └─ Audit history
```

## Documentation

- **[SETUP.md](SETUP.md)** — Developer onboarding (asdf, tools, environment)
- **[CONTRIBUTING.md](CONTRIBUTING.md)** — Development workflow, Go 1.26 linting guide
- **[LINTING_GO_1_26.md](LINTING_GO_1_26.md)** — Docker-based linting for Go 1.26
- **[SPRINT_1_4_RETROSPECTIVE.md](SPRINT_1_4_RETROSPECTIVE.md)** — Sprint planning vs execution
- **[BRANCH_PROTECTION.md](BRANCH_PROTECTION.md)** — Branch rules and CI requirements
- **[SPRINT_DASHBOARD.md](SPRINT_DASHBOARD.md)** — Merge visibility and sprint reports
- **[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)** — Design decisions and technical details
- **[docs/API-REFERENCE.md](docs/API-REFERENCE.md)** — CLI and Go SDK reference
- **[docs/TROUBLESHOOTING.md](docs/TROUBLESHOOTING.md)** — Common issues and solutions

## Testing

```bash
make test              # All unit tests
make test-with-temporal # Integration tests (starts Temporal)
make coverage          # Coverage report
make lint              # Code quality (Docker-based, Go 1.26 compatible)
make hooks-install     # Enable auto-linting on commits
```

## CI/CD Pipeline

All commits to `main` trigger automated testing:

- **verify-versions** — Go/Node versions from `.tool-versions`
- **test-local** — Unit tests (Go 1.26)
- **test-coverage** — Coverage reporting
- **test-bridge & test-events** — Component tests
- **test-integration** — Full integration with Temporal
- **build** — Binary compilation
- **sprint-summary** — Merge visibility report

All jobs must pass before team sharing (branch protection enforced).

## Project Structure

```
├── cmd/orca-temporal/          # CLI entry point
├── internal/
│   ├── bridge/                 # Core bridge logic
│   ├── workflows/              # Temporal workflow definitions
│   ├── events/                 # Event schema + translation
│   └── config/                 # Configuration
├── pkg/bridge/                 # Public SDK
├── .github/workflows/          # CI/CD pipeline
├── .githooks/                  # Git pre-commit hooks
├── Makefile                    # Build targets
├── docker-compose.yml          # Local Temporal (optional)
└── .tool-versions              # asdf pinned versions
```

## Design Decisions

1. **Child Workflows for Workers** → Supports parallelism and isolation
2. **Async Event Flushing** → Preserves ORCA's real-time UX
3. **Event Buffer + Batch** → Balances latency/throughput (100ms or 1K events)
4. **Generic Serialization** → Multi-language agent support
5. **ORCA as Source of Truth** → Temporal is audit/recovery layer

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

Key workflow:
```bash
git checkout -b feature/my-feature
# Make changes
make fmt && make lint && make test
git commit -m "feature: ..."  # Pre-commit hooks run automatically
git push origin feature/my-feature
# Open PR → Team review → Merge
```

## Team Onboarding

**For new team members:**

```bash
# Clone and setup (30 seconds)
git clone https://github.com/ogormans-deptstack/orca-temporal-bridge.git
cd orca-temporal-bridge
make hooks-install

# Now you're ready to develop!
# Every commit auto-checks: fmt, vet, lint, test
git commit -m "feature: ..."
```

**For team leads:**
- See [CONTRIBUTING.md](CONTRIBUTING.md) for Go 1.26 workflow
- See [LINTING_GO_1_26.md](LINTING_GO_1_26.md) for linting strategy
- Share [SETUP.md](SETUP.md) with onboarding engineers

## Future Roadmap

### Sprint 1.5
- [ ] Re-enable CI linting (when Go 1.26 tooling stabilizes)
- [ ] Enable full branch protection rules
- [ ] Workflow implementations and examples

### Sprint 1.6+
- [ ] Production deployment guide
- [ ] Multi-region Temporal support
- [ ] Advanced replay and debugging tools
- [ ] Performance benchmarking and optimization

## License

MIT License — See [LICENSE](LICENSE) for details.

## Questions?

Open an issue or discussion on GitHub.

---

**Status**: ✅ **Ready for team collaboration**
- All core infrastructure complete
- Clear developer path established
- Comprehensive documentation available
- No blockers for production use

**Last Updated**: Sprint 1.4 Complete (November 2026)
