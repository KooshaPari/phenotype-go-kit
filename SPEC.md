# Phenotype Go-Kit Specification

> Go Infrastructure Toolkit

**Version**: 1.0.0 | **Status**: Active | **Last Updated**: 2026-04-02

## Overview

Go infrastructure toolkit providing reusable, production-grade packages extracted from the Phenotype ecosystem. Each package is independently importable with minimal cross-package dependencies.

**Language**: Go 1.24+

**Key Features**:
- Context-scoped logging
- Generic ring buffer
- Polling with backoff
- Thread-safe registry
- Hexagonal architecture support

## Architecture

```
phenotype-go-kit/
├── contracts/              # Port interfaces (hexagonal)
│   ├── ports/
│   │   ├── inbound/        # Driving ports
│   │   └── outbound/       # Driven ports
│   ├── models/             # Domain models
│   └── plugins/            # Plugin interfaces
├── plugins/                # Plugin implementations
│   └── embeddings/         # AI embeddings providers
├── cache/                  # Redis cache implementation
├── logctx/                 # Context-scoped logging
├── ringbuffer/             # Generic ring buffer
├── waitfor/                # Polling with backoff
├── registry/               # Thread-safe registry
└── tests/                  # Integration tests
```

## Packages

| Package | Description | Status |
|---------|-------------|--------|
| `logctx` | Context-scoped `slog.Logger` injection | Stable |
| `ringbuffer` | Generic fixed-capacity circular buffer | Stable |
| `waitfor` | Polling with exponential backoff | Stable |
| `registry` | Thread-safe generic registry | Stable |
| `contracts` | Hexagonal architecture ports | In Development |
| `cache` | Redis cache adapter | In Development |
| `plugins/embeddings` | AI embeddings providers | In Development |

## Quick Start

```bash
go get github.com/KooshaPari/phenotype-go-kit
```

### logctx

```go
import "github.com/KooshaPari/phenotype-go-kit/logctx"

ctx := logctx.WithLogger(context.Background(), slog.Default())
logger := logctx.From(ctx)
logger.Info("hello from context logger")
```

### ringbuffer

```go
import "github.com/KooshaPari/phenotype-go-kit/ringbuffer"

rb := ringbuffer.New[int](100)
rb.Push(1)
val, ok := rb.Pop()
```

### waitfor

```go
import "github.com/KooshaPari/phenotype-go-kit/waitfor"

err := waitfor.Condition(func() bool {
    return service.IsReady()
}, 30*time.Second)
```

### registry

```go
import "github.com/KooshaPari/phenotype-go-kit/registry"

reg := registry.New[string, Service]()
reg.Register("svc1", svc)
svc, found := reg.Get("svc1")
```

## Dependencies

- Go 1.24+
- Minimal external dependencies per package
- Standard library preferred

## Performance Targets

| Metric | Target |
|--------|--------|
| Allocation per log call | 0 (hot path) |
| Ring buffer ops | O(1) |
| Registry lookup | O(1) |

## License

MIT
