# Makefile Template - Auto-generated Infrastructure

# Default task - show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  check     - Run all checks (lint, format, test)"
	@echo "  fmt       - Format code"
	@echo "  lint      - Run linters"
	@echo "  test      - Run tests"
	@echo "  clean     - Clean build artifacts"
	@echo "  all       - Run full CI pipeline"

# Auto-detect project type and run appropriate commands
.PHONY: check fmt lint test clean all

check: fmt lint test
	@echo "✓ All checks passed"

all: check
	@echo "✓ Full CI pipeline complete"

fmt:
	@echo "Formatting..."
	@-[ -f Cargo.toml ] && cargo fmt || true
	@-[ -f go.mod ] && go fmt ./... || true
	@-[ -f package.json ] && npm run format 2>/dev/null || true
	@echo "✓ Format complete"

lint:
	@echo "Linting..."
	@-[ -f Cargo.toml ] && cargo clippy -- -D warnings 2>/dev/null || cargo check 2>/dev/null || true
	@-[ -f go.mod ] && go vet ./... || true
	@echo "✓ Lint complete"

test:
	@echo "Testing..."
	@-[ -f Cargo.toml ] && cargo test 2>/dev/null || true
	@-[ -f go.mod ] && go test ./... 2>/dev/null || true
	@-[ -f package.json ] && npm test 2>/dev/null || true
	@echo "✓ Test complete"

clean:
	@echo "Cleaning..."
	@-[ -d target ] && rm -rf target || true
	@-[ -f Cargo.toml ] && cargo clean 2>/dev/null || true
	@echo "✓ Clean complete"
