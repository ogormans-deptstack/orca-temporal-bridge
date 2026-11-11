.PHONY: help build test test-fast test-bridge test-events test-coverage test-with-temporal \
         clean fmt lint vet run-temporal stop-temporal test-setup dev coverage docker-build docker-run

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
	@echo "  make test-setup        Verify Go + Docker setup"
	@echo ""
	@echo "Quality Checks:"
	@echo "  make fmt               Format code (gofmt)"
	@echo "  make lint              Lint code (golangci-lint)"
	@echo "  make vet               Run go vet"
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
test-with-temporal: run-temporal
	@echo "Running integration tests with Temporal..."
	@echo "Waiting for Temporal to be ready..."
	@sleep 3
	@go test -v -race -timeout 30s -tags=integration ./internal/bridge ./internal/events
	@echo "✓ Integration tests completed"
	@$(MAKE) stop-temporal

# Code quality
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@gofmt -s -w .

lint:
	@echo "Linting code..."
	@which golangci-lint > /dev/null || (echo "Install: brew install golangci-lint" && exit 1)
	@golangci-lint run ./... --timeout=5m

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
	@which go > /dev/null || (echo "✗ Go not installed" && exit 1)
	@echo "✓ Go installed: $$(go version)"
	@which docker > /dev/null || (echo "✗ Docker not installed" && exit 1)
	@echo "✓ Docker installed: $$(docker --version)"
	@which docker-compose > /dev/null || (echo "✗ Docker Compose not installed" && exit 1)
	@echo "✓ Docker Compose installed"
	@echo ""
	@echo "To run full integration tests:"
	@echo "  make test-with-temporal"

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
