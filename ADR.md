# Architecture Decision Records - phenotype-go-kit

## ADR-001: Independent Packages with Zero Inter-Dependencies

**Status**: Accepted
**Context**: Infrastructure utilities must be consumable individually without pulling in unrelated code.
**Decision**: Each package (`logctx`, `ringbuffer`, `waitfor`, `registry`) is self-contained with no imports from sibling packages.
**Consequences**: Consumers add only what they need; no transitive dependency bloat.

## ADR-002: Panic on Missing Logger

**Status**: Accepted
**Context**: A missing logger is a programmer error, not a runtime condition.
**Decision**: `logctx.From` panics if no logger is in context rather than returning a no-op logger.
**Consequences**: Bugs surface immediately in development; no silent log loss in production.

## ADR-003: Quartz Clock for Deterministic Testing

**Status**: Accepted
**Context**: Time-dependent code (`waitfor`) is hard to test with real clocks.
**Decision**: Accept `quartz.Clock` interface; `nil` defaults to real clock.
**Consequences**: Tests are fast and deterministic; single external dependency (`github.com/coder/quartz`).

## ADR-004: Generics for Type Safety

**Status**: Accepted
**Context**: Go 1.22+ supports type parameters.
**Decision**: Use generics for `ringbuffer[T]` and `registry[K, V]` to avoid `interface{}` casts.
**Consequences**: Compile-time type safety; clear API contracts.
