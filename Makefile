.PHONY: test build lint clean tidy install-tools benchmark security

# Install development tools
install-tools:
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest

# Test all packages
test:
	go test ./...

# Test with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run benchmarks
benchmark:
	go test -bench=. -benchmem ./...

# Test with race detection
test-race:
	go test -race ./...

# Build all packages
build:
	go build ./...

# Run linting
lint:
	go vet ./...
	go fmt ./...
	@which goimports > /dev/null 2>&1 && goimports -w . || echo "goimports not found, skipping import formatting"

# Run advanced linting (requires golangci-lint)
lint-advanced:
	golangci-lint run

# Security scan
security:
	gosec ./...

# Clean up
clean:
	go clean ./...
	rm -f coverage.out coverage.html

# Tidy dependencies
tidy:
	go mod tidy
	go mod verify

# Run all checks
check: lint test build

# Full CI pipeline locally
ci: install-tools tidy lint-advanced test-race test-coverage benchmark build

# Release preparation
release-prep: ci
	@echo "Package is ready for release!"

# Help target
help:
	@echo "Available targets:"
	@echo "  install-tools    - Install development tools"
	@echo "  test            - Run tests"
	@echo "  test-coverage   - Run tests with coverage report"
	@echo "  test-race       - Run tests with race detection"
	@echo "  benchmark       - Run benchmark tests"
	@echo "  build           - Build all packages"
	@echo "  lint            - Run basic linting and formatting"
	@echo "  lint-advanced   - Run advanced linting (requires golangci-lint)"
	@echo "  security        - Run security scan"
	@echo "  clean           - Clean build artifacts"
	@echo "  tidy            - Tidy go modules"
	@echo "  check           - Run basic checks (lint, test, build)"
	@echo "  ci              - Run full CI pipeline locally"
	@echo "  release-prep    - Prepare package for release"