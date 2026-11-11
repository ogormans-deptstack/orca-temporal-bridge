#!/bin/bash

# Local test runner for ORCA-Temporal Bridge (Sprint 1.3)
# Usage: ./scripts/test-local.sh [options]
#
# Options:
#   (none)     Run all local tests
#   --fast     Quick smoke test only
#   --coverage Generate coverage report
#   --bridge   Bridge tests only
#   --events   Event tests only
#   --watch    Run tests on file changes (requires entr)
#   --help     Show this help

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;36m'
NC='\033[0m' # No Color

# Functions
print_header() {
    echo -e "${BLUE}==>${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}!${NC} $1"
}

check_go() {
    if ! command -v go &> /dev/null; then
        print_error "Go not installed"
        exit 1
    fi
    print_success "Go $(go version | awk '{print $3}')"
}

run_fast_tests() {
    print_header "Running fast smoke test (< 1s)"
    go test -timeout 10s ./internal/bridge ./internal/events
    print_success "Fast tests passed"
}

run_all_tests() {
    print_header "Running all tests (no Temporal required)"
    go test -v -race -timeout 10s ./internal/...
    print_success "All tests passed"
}

run_bridge_tests() {
    print_header "Running Bridge E2E tests"
    go test -v -timeout 10s ./internal/bridge/...
    print_success "Bridge tests passed"
}

run_event_tests() {
    print_header "Running Event translator tests"
    go test -v -timeout 10s ./internal/events/...
    print_success "Event tests passed"
}

run_with_coverage() {
    print_header "Running tests with coverage"
    go test -v -race -timeout 10s -coverprofile=coverage.out ./internal/...
    
    print_header "Coverage summary"
    go tool cover -func=coverage.out | grep total
    
    print_header "Generating HTML report"
    go tool cover -html=coverage.out -o coverage.html
    
    print_success "Coverage report: coverage.html"
    print_warning "Open in browser: open coverage.html"
}

run_quality_checks() {
    print_header "Running quality checks"
    
    # Format check
    print_warning "Checking code format..."
    if ! go fmt ./... > /dev/null 2>&1; then
        print_error "Code formatting issues found"
        go fmt ./...
        return 1
    fi
    print_success "Code format OK"
    
    # Vet
    print_warning "Running go vet..."
    if ! go vet ./... > /dev/null 2>&1; then
        print_error "Vet issues found"
        go vet ./...
        return 1
    fi
    print_success "Go vet OK"
    
    return 0
}

run_watch_mode() {
    if ! command -v entr &> /dev/null; then
        print_error "entr not installed. Install with: brew install entr"
        exit 1
    fi
    
    print_header "Running tests in watch mode"
    print_warning "Press Ctrl+C to stop"
    
    find "$PROJECT_ROOT"/internal -name "*.go" | entr -c make test
}

show_help() {
    cat << EOF
ORCA-Temporal Bridge — Local Test Runner

Usage: ./scripts/test-local.sh [options]

Options:
  (none)         Run all local tests
  --fast         Quick smoke test only (< 1s)
  --coverage     Generate HTML coverage report
  --bridge       Bridge E2E tests only
  --events       Event translator tests only
  --quality      Run fmt + vet checks only
  --watch        Run tests on file changes (requires entr)
  --help         Show this help message

Examples:
  ./scripts/test-local.sh                    # Run all tests
  ./scripts/test-local.sh --fast             # Quick smoke test
  ./scripts/test-local.sh --coverage         # Coverage report
  ./scripts/test-local.sh --watch            # Watch mode

Environment:
  VERBOSE=1     Show more output

EOF
}

# Main
main() {
    cd "$PROJECT_ROOT"
    
    print_header "ORCA-Temporal Bridge — Test Runner (Sprint 1.3)"
    echo ""
    
    check_go
    echo ""
    
    # Handle arguments
    case "${1:-}" in
        --fast)
            run_fast_tests
            ;;
        --coverage)
            run_all_tests
            echo ""
            run_with_coverage
            ;;
        --bridge)
            run_bridge_tests
            ;;
        --events)
            run_event_tests
            ;;
        --quality)
            run_quality_checks
            ;;
        --watch)
            run_watch_mode
            ;;
        --help)
            show_help
            ;;
        "")
            # Default: run all tests + quality checks
            run_quality_checks
            echo ""
            run_all_tests
            echo ""
            print_success "All checks passed! Ready to commit."
            ;;
        *)
            print_error "Unknown option: $1"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

main "$@"
