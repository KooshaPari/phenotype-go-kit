# phenotype-go-kit

## Project

Go module containing generic infrastructure packages extracted from the Phenotype ecosystem.
Each package is independent with minimal dependencies.

## Stack

- **Language**: Go 1.24+
- **Test**: `go test ./...`
- **Lint**: `go vet ./...`
- **Format**: `gofumpt`

## Structure

```
logctx/      # Context-scoped slog logger
ringbuffer/  # Generic ring buffer
waitfor/     # Polling with timeout (uses github.com/coder/quartz)
registry/    # Thread-safe registry with ref counting & quota
```

## Conventions

- Each package is self-contained with `_test.go` alongside
- No inter-package dependencies
- External deps minimized (only `quartz` for waitfor)
- Thread-safe by default where applicable (registry uses sync.RWMutex)
- Generic type parameters where appropriate (ringbuffer, registry)
