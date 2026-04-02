.PHONY: all fmt lint test vet build clean install-deps

# Default target
all: fmt lint test

# Format check
fmt:
	go fmt ./...

# Lint
lint:
	golangci-lint run --timeout=5m

# Run go vet
vet:
	go vet ./...

# Run tests
test:
	go test -v -race -count=1 ./...

# Run tests with coverage
test-cover:
	go test -v -race -count=1 -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Build
build:
	go build ./...

# Build release
build-release:
	go build -ldflags="-s -w" ./...

# Tidy dependencies
tidy:
	go mod tidy

# Verify dependencies
verify:
	go mod verify

# Clean build artifacts
clean:
	go clean

# Install pre-commit hooks
install-deps:
	which golangci-lint || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	pre-commit install || echo "Install pre-commit: pip install pre-commit"

# Run all quality gates
qa: fmt vet lint test

# Help
help:
	@echo "Available targets:"
	@echo "  all         - Run fmt, lint, and test (default)"
	@echo "  fmt         - Format code"
	@echo "  lint        - Run golangci-lint"
	@echo "  vet         - Run go vet"
	@echo "  test        - Run tests"
	@echo "  test-cover  - Run tests with coverage"
	@echo "  build       - Build"
	@echo "  build-release - Build release"
	@echo "  tidy        - Tidy dependencies"
	@echo "  verify      - Verify dependencies"
	@echo "  clean       - Clean build artifacts"
	@echo "  qa          - Run all quality gates"
	@echo "  help        - Show this help"
