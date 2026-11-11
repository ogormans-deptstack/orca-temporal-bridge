# Setup Guide — ORCA-Temporal Bridge

**For new team members** — Getting started with the project.

## Quick Start (5 minutes)

### 1. Clone the Repository
```bash
git clone https://github.com/ogormans-deptstack/orca-temporal-bridge.git
cd orca-temporal-bridge
```

### 2. Check Requirements
```bash
make test-setup
```

**Expected output:**
```
=== Tool Versions (from .tool-versions) ===
golang 1.21.6
nodejs 22.0.0

=== Installed Versions ===
✓ Go: go1.21.6
✓ Docker: 27.0.0
✓ Docker Compose: v2.0.0

=== asdf Status ===
✓ asdf installed: v0.14.0
  Run: asdf install
  Then: asdf local golang <version>
```

If you see warnings, follow the install instructions.

### 3. Install Go Version (if using asdf)
```bash
make asdf-install
make asdf-verify
```

Or just use your system Go (1.21 or 1.22).

### 4. Run Tests
```bash
make test-fast    # Quick check (< 1s)
make test         # Full local tests (59 tests)
```

**Expected:** All tests pass ✅

### 5. Build the Binary
```bash
make build
./bin/orca-temporal --help
```

**Done!** You're ready to develop.

---

## Full Setup (15 minutes)

### Prerequisites

**Required:**
- macOS 12+ or Linux (Ubuntu 20.04+)
- Git
- Go 1.21 or 1.22 (managed by asdf, or install manually)
- Docker and Docker Compose

**Optional but recommended:**
- asdf (version manager)
- Temporal CLI (for standalone testing)

### Step 1: Install asdf

**macOS:**
```bash
brew install asdf
asdf plugin add golang https://github.com/asdf-community/asdf-golang.git
asdf plugin add nodejs https://github.com/asdf-vm/asdf-nodejs.git
```

**Linux (Ubuntu):**
```bash
git clone https://github.com/asdf-vm/asdf.git ~/.asdf --branch v0.14.0
# Add to ~/.bashrc or ~/.zshrc:
# . "$HOME/.asdf/asdf.sh"
# . "$HOME/.asdf/completions/asdf.bash"
source ~/.bashrc

asdf plugin add golang https://github.com/asdf-community/asdf-golang.git
asdf plugin add nodejs https://github.com/asdf-vm/asdf-nodejs.git
```

**Verify:**
```bash
asdf --version  # Should print v0.14.0+
```

### Step 2: Clone and Install Tools
```bash
git clone https://github.com/ogormans-deptstack/orca-temporal-bridge.git
cd orca-temporal-bridge
git maintenance start  # Enable background optimizations

make asdf-install     # Install Go + Node versions
make asdf-verify      # Verify versions match .tool-versions
```

### Step 3: Install Development Dependencies
```bash
# Install linter (for local linting, optional in CI)
brew install golangci-lint

# Install Temporal CLI (optional, for advanced testing)
brew tap temporalio/temporal
brew install temporal-cli
```

### Step 4: Verify Setup
```bash
make test-setup       # Full environment check
make test-fast        # Quick smoke test
make test             # All tests (59/59 should pass)
```

### Step 5: Configure Git Hooks (Optional)
```bash
# Auto-format on commit
cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash
go fmt ./...
git add -u
EOF
chmod +x .git/hooks/pre-commit
```

---

## Development Workflow

### Day-to-Day

**Before you start coding:**
```bash
make test-setup          # Check environment
git pull origin main     # Get latest
```

**While coding (use one of these):**
```bash
task test:watch          # Auto-test on file changes (if using Taskfile)
make test-fast           # Quick tests
make test                # All tests
```

**Before committing:**
```bash
make fmt                 # Format code
make vet                 # Check for issues
make test                # Full test suite
```

**Before pushing:**
```bash
make build               # Ensure binary builds
git push origin <branch> # Push to PR
```

### Running Integration Tests

**With Temporal:**
```bash
# Option 1: Full integration (starts/stops Temporal)
make test-with-temporal

# Option 2: Manual (more control)
make run-temporal
make test                # Run tests against Temporal at localhost:7233
make stop-temporal
```

**View Temporal UI:**
```bash
# After running make run-temporal
open http://localhost:8080
```

### Common Tasks

| Task | Command | Purpose |
|------|---------|---------|
| Run tests | `make test` | Verify code works |
| Generate coverage report | `make coverage` | See test coverage (opens HTML) |
| Format code | `make fmt` | Ensure consistent style |
| Check for issues | `make vet` | Find potential bugs |
| Run linter | `make lint` | Check code quality |
| Build binary | `make build` | Create executable |
| Clean build artifacts | `make clean` | Reset to clean state |
| View help | `make help` | List all commands |

---

## Working with Taskfile (Modern Alternative)

Taskfile provides a more modern, cross-platform task runner. Optional but recommended for better UX.

### Install Taskfile
```bash
brew install go-task/tap/go-task
```

### Common Tasks
```bash
task --list                # See all available tasks
task dev                   # Full dev workflow (fmt → build → test)
task test:watch            # Auto-test on changes
task test:fast             # Quick smoke test
task temporal:start        # Start Temporal server
task temporal:stop         # Stop Temporal server
task ci:local              # Simulate GitHub Actions locally
```

**Example workflow:**
```bash
task dev                   # Full validation before commit
git commit -m "feature: add new functionality"
git push origin feature-branch
# GitHub Actions CI runs automatically
# Once all checks pass and PR is reviewed, you can merge
```

---

## Understanding the CI/CD Pipeline

### What Runs on Every Push/PR

**Location:** `.github/workflows/test.yml`

**Jobs (7 parallel, ~2 minutes total):**

1. **verify-versions** — Checks `.tool-versions` consistency
2. **test-local** — Tests on Go 1.21 and 1.22
3. **test-coverage** — Full coverage report (uploads to Codecov)
4. **test-bridge** — Bridge E2E tests
5. **test-events** — Event translator tests
6. **test-integration** — Tests with real Temporal
7. **build** — Compiles binary

**What happens if a job fails:**
- Merge button is disabled ❌
- You see ✗ on your PR
- Click "Details" to see error logs
- Fix the issue, push again
- CI re-runs automatically

### Branch Protection

**Main branch is protected:**
- All CI checks must pass
- At least 1 review required
- Conversation must be resolved
- Branches must be up-to-date

### Sprint-End Merge Report

After each merge to main:
- GitHub Actions generates a summary report
- Shows what was deployed, test results, coverage
- Available as artifact in Actions tab
- Shared with team via PR

---

## Troubleshooting

### "Go version mismatch"
```bash
# Verify your Go version
go version

# Should match .tool-versions
cat .tool-versions | grep golang

# If using asdf, fix it:
make asdf-verify
# or
asdf global golang 1.21.6
```

### "Temporal won't start"
```bash
# Check Docker
docker ps

# If Docker is not running:
# macOS: open -a Docker
# Linux: sudo systemctl start docker

# Try again:
make run-temporal
```

### "Tests fail locally but pass on CI"
```bash
# Run with race detector (what CI uses):
go test -race ./...

# If that fails, you have a concurrency issue
# Check the specific failing test:
make test-bridge   # Or test-events, etc.
```

### "Import package not found"
```bash
# Fetch dependencies
go mod download
go mod tidy

# Then retry:
make test
```

### "golangci-lint not installed"
```bash
# For local linting (optional)
brew install golangci-lint

# Or skip it in CI (automatic fallback)
make test  # Uses built-in go vet instead
```

---

## IDE/Editor Setup

### VS Code

1. Install Go extension: `golang.go`
2. Open the project folder
3. Prompt to install Go tools → Click "Install All"

**.vscode/settings.json** (optional):
```json
{
  "go.lintOnSave": "package",
  "go.lintTool": "golangci-lint",
  "go.formatOnSave": true,
  "go.useLanguageServer": true,
  "[go]": {
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
      "source.organizeImports": "explicit"
    }
  }
}
```

### GoLand / IntelliJ IDEA

1. Open project
2. Go → Enable Go Modules (if not automatic)
3. Settings → Go → Go Modules → Check "Enable Go Modules integration"
4. File → Invalidate Caches → Restart

### Vim/Neovim

Use your favorite LSP client with gopls:
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

---

## Resources

| Topic | Link |
|-------|------|
| **Testing** | See `TESTING.md` (complete testing guide) |
| **Taskfile** | See `TASKFILE.md` (tasks reference) |
| **Sprint Dashboard** | See `SPRINT_DASHBOARD.md` (merge visibility) |
| **Branch Protection** | See `BRANCH_PROTECTION.md` (safety rules) |
| **CI/CD** | See `.github/workflows/test.yml` (GitHub Actions) |
| **Makefile** | Run `make help` (all build targets) |

---

## Getting Help

### Stuck?

1. **Check the docs** — Especially `TESTING.md` and `SPRINT_DASHBOARD.md`
2. **Run make help** — See all available commands
3. **Check error logs** — Full logs in GitHub Actions (if CI failure)
4. **Ask on Slack** — #engineering-bridge channel

### Reporting Issues

**Found a bug?**
1. Create a GitHub issue
2. Include: what you did, what you expected, what happened
3. Attach error logs (from `make test` output)
4. Create a PR with a fix (or comment "help wanted")

---

## Next Steps

1. ✅ Run `make test-fast` to verify setup
2. ✅ Read `TESTING.md` to understand testing
3. ✅ Create a small test PR to try the workflow
4. ✅ Read `SPRINT_DASHBOARD.md` to see what gets deployed

**Welcome to the team!** 🎉
