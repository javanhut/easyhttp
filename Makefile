# EasyHTTP Makefile

.PHONY: all build test bench clean lint fmt deps example

# Variables
BINARY_NAME=easyhttp-example
GO_VERSION=1.21

# Default target
all: fmt lint test

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, skipping lint"; \
		go vet ./...; \
	fi

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Build example
build:
	@echo "Building example..."
	go build -o $(BINARY_NAME) ./example

# Run example
example: build
	@echo "Running example..."
	./$(BINARY_NAME)

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	go clean
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

# Run tests in verbose mode
test-verbose:
	@echo "Running verbose tests..."
	go test -v -race ./...

# Check for security vulnerabilities
security:
	@echo "Checking for security vulnerabilities..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not installed, install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Check module dependencies
deps-check:
	@echo "Checking dependencies..."
	go list -u -m all
	go mod verify

# Update dependencies
deps-update:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

# Tag a new version (usage: make tag VERSION=v1.0.0)
tag:
	@if [ -z "$(VERSION)" ]; then \
		echo "Usage: make tag VERSION=v1.0.0"; \
		exit 1; \
	fi
	@echo "Tagging version $(VERSION)..."
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push origin $(VERSION)

# Development setup
dev-setup:
	@echo "Setting up development environment..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@echo "Development tools installed!"

# Help
help:
	@echo "Available targets:"
	@echo "  all          - Format, lint, and test"
	@echo "  deps         - Install dependencies"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  bench        - Run benchmarks"
	@echo "  build        - Build example"
	@echo "  example      - Build and run example"
	@echo "  clean        - Clean build artifacts"
	@echo "  security     - Check for security vulnerabilities"
	@echo "  deps-check   - Check dependencies"
	@echo "  deps-update  - Update dependencies"
	@echo "  tag          - Tag a new version (make tag VERSION=v1.0.0)"
	@echo "  dev-setup    - Install development tools"
	@echo "  help         - Show this help"
