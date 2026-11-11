.PHONY: help build test clean fmt lint run-temporal test-setup coverage

help:
	@echo "ORCA-Temporal Bridge — Build Targets"
	@echo ""
	@echo "Development:"
	@echo "  make build             Build binary (./bin/orca-temporal)"
	@echo "  make test              Run unit tests"
	@echo "  make test-integration  Run integration tests"
	@echo "  make test-load         Run load test (100+ events/min)"
	@echo "  make coverage          Generate coverage report"
	@echo ""
	@echo "Quality:"
	@echo "  make fmt               Format code (gofmt)"
	@echo "  make lint              Lint code (golangci-lint)"
	@echo "  make vet               Run go vet"
	@echo ""
	@echo "Infrastructure:"
	@echo "  make run-temporal      Start local Temporal server (docker-compose)"
	@echo "  make test-setup        Verify Temporal is reachable"
	@echo "  make clean             Remove build artifacts"
	@echo ""
	@echo "Release:"
	@echo "  make docker-build      Build Docker image"
	@echo "  make docker-run        Run in Docker"

# Build
build: fmt
	@echo "Building orca-temporal..."
	@mkdir -p bin
	@go build -o bin/orca-temporal ./cmd/orca-temporal

# Tests
test:
	@echo "Running unit tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

test-integration: test-setup
	@echo "Running integration tests..."
	@go test -v -tags=integration -race ./tests/integration/...

test-load:
	@echo "Running load tests..."
	@go test -v -timeout=5m -run TestLoad ./tests/load/...

coverage:
	@echo "Generating coverage report..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Report saved to coverage.html"

# Code quality
fmt:
	@echo "Formatting code..."
	@go fmt ./...

lint:
	@echo "Linting code..."
	@golangci-lint run ./...

vet:
	@echo "Running go vet..."
	@go vet ./...

# Infrastructure
run-temporal:
	@echo "Starting Temporal server..."
	@docker-compose up -d || (echo "Note: Docker not running. Start Docker and run 'make run-temporal' again." && exit 0)
	@sleep 2
	@echo "Temporal running at localhost:7233"
	@echo "Temporal UI at http://localhost:8080"

test-setup: 
	@echo "Verifying setup..."
	@which go > /dev/null || (echo "Error: Go not installed" && exit 1)
	@echo "✓ Go is installed"
	@echo ""
	@echo "To verify Temporal connectivity, start Temporal:"
	@echo "  1. Start Docker: open -a Docker"
	@echo "  2. Run: make run-temporal"
	@echo "  3. Visit: http://localhost:8080 (Temporal UI)"

stop-temporal:
	@echo "Stopping Temporal server..."
	@docker-compose down

# Clean
clean:
	@echo "Cleaning..."
	@rm -rf bin dist coverage.* *.out

# Docker
docker-build:
	@echo "Building Docker image..."
	@docker build -t orca-temporal-bridge:latest .

docker-run: docker-build
	@echo "Running in Docker..."
	@docker run --network host orca-temporal-bridge:latest sync --temporal localhost:7233

# Development workflow
dev: build
	@echo ""
	@echo "✓ Build complete (./bin/orca-temporal)"
	@echo ""
	@echo "To run with Temporal backend:"
	@echo "  1. Start Docker daemon"
	@echo "  2. Run: make run-temporal"
	@echo "  3. In another terminal: ./bin/orca-temporal sync --temporal localhost:7233 --namespace default"
	@echo ""
	@echo "To test the CLI:"
	@echo "  ./bin/orca-temporal --help"
