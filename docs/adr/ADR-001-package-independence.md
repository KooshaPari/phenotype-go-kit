# ADR-001: Package Independence Policy

## Status
**Accepted** | 2024-01-15

## Context

phenotype-go-kit serves as a shared infrastructure library across multiple Phenotype services.
Each service may need different subsets of functionality. We need to avoid monolithic dependencies.

**Problem:** Adding a dependency to one package shouldn't force all consumers to update their dependency tree.

## Decision

### Rule 1: Zero Dependency Core
Core packages (`logctx`, `ringbuffer`, `registry`, `versioning`, `validation`) have **zero external dependencies**.

### Rule 2: Layered Dependencies
```
presentation → application → domain → (no external deps)
     ↓            ↓           ↓
infrastructure (can have external deps)
```

### Rule 3: Interface-First
Core interfaces are defined in domain packages, implementations in infrastructure.

### Rule 4: No Circular Dependencies
Package imports form a DAG (directed acyclic graph).

## Implementation

```go
// domain/logctx/logctx.go - Zero dependencies
package logctx

func FromContext(ctx context.Context) *Logger  // Standard library only
```

```go
// infrastructure/logctx/slog.go - External dep allowed
package slog

type SlogAdapter struct { /* implements logctx.Logger */ }
```

## Consequences

### Positive
- Services can import only what they need
- Reduced dependency conflicts
- Faster builds
- Easier testing

### Negative
- More package boundaries to maintain
- Slightly more boilerplate for adapters

## Alternatives Considered

1. **Mono-package with build tags** - rejected, adds complexity
2. **Multiple modules** - rejected, complicates versioning
3. **Full dependency injection** - rejected, too much boilerplate

## References
- [Go Wiki: How to Write Go Code](https://github.com/golang/go/wiki/CodeTools)
- Hexagonal Architecture: Ports & Adapters
