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

- Go 1.21+
- Docker (for Temporal server)
- `go.temporal.io/sdk` installed

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

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT License — See [LICENSE](LICENSE) for details.

## Contact

Questions? Open an issue or discussion on GitHub.

## Status

🟡 **In Development** — Phase 2 (Generate) underway. v0.1.0 coming soon.
