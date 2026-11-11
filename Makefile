.PHONY: help build test test-fast test-bridge test-events test-coverage test-with-temporal \
         clean fmt lint lint-local vet run-temporal stop-temporal test-setup dev coverage docker-build docker-run \
         hooks-install hooks-uninstall

help:
	@echo "ORCA-Temporal Bridge — Build Targets"
	@echo ""
	@echo "Local Development (Sprint 1.3):"
	@echo "  make test              Run all tests (no Temporal required)"
	@echo "  make test-fast         Quick unit tests (no coverage)"
	@echo "  make test-bridge       Run bridge E2E tests only"
	@echo "  make test-events       Run event translator tests"
	@echo "  make test-coverage     Run tests + coverage report (HTML)"
	@echo ""
	@echo "Integration Testing (requires Docker):"
	@echo "  make test-with-temporal Run full E2E with real Temporal"
	@echo "  make run-temporal      Start local Temporal server (docker-compose)"
	@echo "  make stop-temporal     Stop Temporal server"
	@echo ""
	@echo "Code Quality (Go 1.26 compatible):"
	@echo "  make fmt               Format code (gofmt)"
	@echo "  make lint              Lint code (golangci-lint via Docker)"
	@echo "  make lint-local        Lint code locally (if compatible)"
	@echo "  make vet               Run go vet"
	@echo ""
	@echo "Git Hooks (optional):"
	@echo "  make hooks-install     Install pre-commit hooks"
	@echo "  make hooks-uninstall   Remove pre-commit hooks"
	@echo ""
	@echo "Version Management (asdf):"
	@echo "  make test-setup        Verify environment + asdf status"
	@echo "  make asdf-install      Install versions from .tool-versions"
	@echo "  make asdf-verify       Verify installed versions match .tool-versions"
	@echo ""
	@echo "Build & Release:"
	@echo "  make build             Build binary (./bin/orca-temporal)"
	@echo "  make docker-build      Build Docker image"
	@echo "  make clean             Remove build artifacts"

# Build
build: fmt vet
	@echo "Building orca-temporal..."
	@mkdir -p bin
	@go build -o bin/orca-temporal ./cmd/orca-temporal
	@echo "✓ Binary: ./bin/orca-temporal"

# Local Tests (no Temporal required)
test:
	@echo "Running all tests (no Temporal required)..."
	@go test -v -race -timeout 10s ./internal/...
	@echo "✓ All tests passed"

test-fast:
	@echo "Running fast unit tests..."
	@go test -timeout 10s ./internal/bridge ./internal/events
	@echo "✓ Fast tests completed"

test-bridge:
	@echo "Running Bridge E2E tests..."
	@go test -v -timeout 10s ./internal/bridge/...
	@echo "✓ Bridge tests completed"

test-events:
	@echo "Running Event translator tests..."
	@go test -v -timeout 10s ./internal/events/...
	@echo "✓ Event tests completed"

test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -race -timeout 10s -coverprofile=coverage.out ./internal/...
	@go tool cover -func=coverage.out | grep total
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report: coverage.html"

# Integration Tests (requires Temporal)
test-with-temporal:
	@echo "Running integration tests with Temporal..."
	@echo "Checking if Temporal is already running..."
	@if ! nc -z localhost 7233 2>/dev/null; then \
		echo "Temporal not running, starting..."; \
		$(MAKE) run-temporal; \
	else \
		echo "✓ Temporal already running"; \
	fi
	@echo "Waiting for Temporal to be ready..."
	@sleep 2
	@go test -v -race -timeout 30s -tags=integration ./internal/bridge ./internal/events
	@echo "✓ Integration tests completed"

# Code quality
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@gofmt -s -w .

lint:
	@echo "Linting code with golangci-lint (Docker)..."
	@echo "Note: Using Docker to support Go 1.26 (linter tooling not yet available for Go 1.26)"
	@if command -v docker &> /dev/null; then \
		docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:latest golangci-lint run --config=.golangci.yml --timeout=5m; \
	else \
		echo "⚠ Docker not available. Install via: brew install docker"; \
		exit 1; \
	fi

lint-local:
	@echo "Linting code locally (if compatible version available)..."
	@which golangci-lint > /dev/null || (echo "⚠ golangci-lint not installed. Run: brew install golangci-lint" && exit 1)
	@echo "Note: May fail with Go 1.26 if golangci-lint was built with older Go"
	@golangci-lint run ./... --config=.golangci.yml --timeout=5m || echo "ℹ Try 'make lint' (Docker-based) for Go 1.26 support"

vet:
	@echo "Running go vet..."
	@go vet ./...

# Temporal infrastructure
run-temporal:
	@echo "Starting Temporal server..."
	@docker-compose up -d
	@echo "Waiting for Temporal to be ready..."
	@sleep 3
	@if docker-compose ps temporal | grep -q Up; then \
		echo "✓ Temporal running at localhost:7233"; \
		echo "✓ Temporal UI at http://localhost:8080"; \
	else \
		echo "✗ Failed to start Temporal"; \
		exit 1; \
	fi

stop-temporal:
	@echo "Stopping Temporal server..."
	@docker-compose down

test-setup:
	@echo "Verifying local setup..."
	@echo ""
	@echo "=== Tool Versions (from .tool-versions) ==="
	@if [ -f .tool-versions ]; then \
		cat .tool-versions; \
		echo ""; \
	else \
		echo "⚠ .tool-versions not found"; \
	fi
	@echo ""
	@echo "=== Installed Versions ==="
	@which go > /dev/null || (echo "✗ Go not installed" && exit 1)
	@echo "✓ Go: $$(go version | awk '{print $$3}')"
	@which docker > /dev/null || (echo "✗ Docker not installed" && exit 1)
	@echo "✓ Docker: $$(docker --version | cut -d' ' -f3 | tr -d ',')"
	@which docker-compose > /dev/null || (echo "✗ Docker Compose not installed" && exit 1)
	@echo "✓ Docker Compose: $$(docker-compose --version | awk '{print $$3}' | tr -d ',')"
	@echo ""
	@echo "=== asdf Status ==="
	@if which asdf > /dev/null; then \
		echo "✓ asdf installed: $$(asdf --version)"; \
		echo "  Run: asdf install"; \
		echo "  Then: asdf local golang <version>"; \
	else \
		echo "⚠ asdf not installed (optional)"; \
		echo "  Install: brew install asdf"; \
		echo "  Setup: asdf plugin add golang"; \
	fi
	@echo ""
	@echo "To run full integration tests:"
	@echo "  make test-with-temporal"

# asdf version management
asdf-install:
	@echo "Installing tools via asdf..."
	@which asdf > /dev/null || (echo "✗ asdf not installed. Install: brew install asdf" && exit 1)
	@echo "Setting up asdf plugins..."
	@asdf plugin add golang 2>/dev/null || echo "  (golang plugin already added)"
	@asdf plugin add nodejs 2>/dev/null || echo "  (nodejs plugin already added)"
	@echo ""
	@echo "Installing versions from .tool-versions..."
	@asdf install
	@echo "✓ asdf install complete"
	@echo ""
	@echo "Verify:"
	@echo "  asdf versions"
	@echo "  go version"
	@echo "  node --version"

asdf-verify:
	@echo "Verifying asdf versions..."
	@which asdf > /dev/null || (echo "✗ asdf not installed" && exit 1)
	@echo ""
	@GOLANG_EXPECTED=$$(grep "^golang" .tool-versions | awk '{print $$2}'); \
	GOLANG_ACTUAL=$$(go version | awk '{print $$3}' | sed 's/go//'); \
	echo "Go expected: $$GOLANG_EXPECTED"; \
	echo "Go actual: $$GOLANG_ACTUAL"; \
	if [[ "$$GOLANG_ACTUAL" =~ ^$$GOLANG_EXPECTED ]]; then \
		echo "✓ Go version matches"; \
	else \
		echo "✗ Go version mismatch!"; \
		echo "  Run: asdf global golang $$GOLANG_EXPECTED"; \
		exit 1; \
	fi
	@echo ""
	@echo "✓ All versions verified"

coverage:
	@echo "Generating coverage report..."
	@go test -coverprofile=coverage.out ./internal/...
	@go tool cover -func=coverage.out | tail -20
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Report saved to coverage.html"

# Clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin dist coverage.* *.out *.html
	@go clean ./...
	@echo "✓ Clean complete"

# Docker
docker-build:
	@echo "Building Docker image..."
	@docker build -t orca-temporal-bridge:latest .

docker-run: docker-build run-temporal
	@echo "Running in Docker..."
	@docker run --network host orca-temporal-bridge:latest sync --temporal localhost:7233
	@$(MAKE) stop-temporal

# Development workflow
dev: build test
	@echo ""
	@echo "✓ Build complete (./bin/orca-temporal)"
	@echo ""
	@echo "To run with Temporal backend:"
	@echo "  1. make run-temporal"
	@echo "  2. ./bin/orca-temporal sync --temporal localhost:7233 --namespace default"
	@echo ""
	@echo "To test the CLI:"
	@echo "  ./bin/orca-temporal --help"
	@echo ""
	@echo "To view coverage:"
	@echo "  make test-coverage  # opens coverage.html"

# Git hooks
hooks-install:
	@echo "Installing git hooks..."
	@mkdir -p .git/hooks
	@cp .githooks/pre-commit .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@git config core.hooksPath .githooks
	@echo "✓ Pre-commit hook installed"
	@echo "  Hook will run: fmt, vet, lint before commits"
	@echo "  Bypass with: git commit --no-verify"

hooks-uninstall:
	@echo "Removing git hooks..."
	@rm -f .git/hooks/pre-commit
	@echo "✓ Pre-commit hook removed"
