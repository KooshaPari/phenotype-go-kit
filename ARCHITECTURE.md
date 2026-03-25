# phenotype-go-kit Architecture

## Overview

Go module containing generic infrastructure packages for the Phenotype ecosystem.
Each package is independent with minimal dependencies, designed for composability.

## Design Philosophy

- **Plug-and-Play**: Each package is self-contained and swappable
- **Interface-First**: Define contracts before implementations (CDD)
- **Testable**: Every package has companion `_test.go` (TDD)
- **Documented**: Public APIs have doc comments (RDD)
- **Minimal Dependencies**: Core packages have zero external deps

## Applied Patterns

### Clean Architecture Layers

```
┌─────────────────────────────────────────────────────────────┐
│                      PRESENTATION                            │
│   CLI, API Clients, Web UIs, Dashboards                    │
├─────────────────────────────────────────────────────────────┤
│                    APPLICATION                               │
│   Use Cases, Commands, Queries, DTOs                        │
├─────────────────────────────────────────────────────────────┤
│                       DOMAIN                                │
│   Entities, Value Objects, Services, Interfaces            │
├─────────────────────────────────────────────────────────────┤
│                   INFRASTRUCTURE                            │
│   Adapters: Redis, Postgres, S3, HTTP Clients              │
└─────────────────────────────────────────────────────────────┘
```

### Hexagonal Architecture (Ports & Adapters)

```
┌──────────────────────────────────────────────────────────────┐
│                        CORE (Domain)                          │
│  ┌────────────┐  ┌────────────┐  ┌────────────────────┐   │
│  │  Entities  │  │  Services  │  │  Value Objects     │   │
│  └────────────┘  └────────────┘  └────────────────────┘   │
│                              │                              │
│  ┌───────────────────────────┴───────────────────────────┐  │
│  │              Ports (Interfaces)                       │  │
│  │  Repository │ Cache │ Logger │ Metrics │ Notifier     │  │
│  └───────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────┘
              │                           ▲
              ▼                           │
┌─────────────────────────────────────────────────────────────┐
│                    Adapters (Implementations)                │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────────┐  │
│  │  Redis   │ │ Postgres │ │  Slog    │ │ Prometheus   │  │
│  └──────────┘ └──────────┘ └──────────┘ └──────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## Directory Structure

### Core Infrastructure (Domain)

| Package | Purpose | Dependencies |
|---------|---------|--------------|
| `logctx/` | Context-scoped structured logging | None |
| `ringbuffer/` | Generic thread-safe ring buffer | None |
| `waitfor/` | Polling with exponential backoff | quartz |
| `registry/` | Thread-safe registry with ref counting | None |
| `versioning/` | Semantic versioning utilities | None |
| `validation/` | Input validation helpers | None |

### Application Services

| Package | Purpose | Dependencies |
|---------|---------|--------------|
| `cache/` | Caching with invalidation | redis |
| `metrics/` | Metrics collection & export | prometheus |
| `tracing/` | Distributed tracing (OTEL) | otel |
| `alerting/` | Alert rules and notifications | None |

### Infrastructure Adapters

| Package | Purpose | Dependencies |
|---------|---------|--------------|
| `auth/` | JWT authentication | None |
| `storage/` | File/S3 storage | AWS SDK |
| `db/` | Database utilities | postgres |
| `bus/` | Message bus | kafka/nats |
| `repository/` | Repository pattern impl | db |

### Presentation

| Package | Purpose | Dependencies |
|---------|---------|--------------|
| `cli/` | CLI utilities | cobra |
| `frontend/` | Web UI helpers | None |
| `health/` | Health check endpoints | None |
| `dashboards/` | Dashboard configs | None |

## SOLID Principles Applied

### Single Responsibility
Each package has one concern:
- `logctx` → Logging only
- `cache` → Caching only
- `metrics` → Metrics only

### Open/Closed
Extensions via interfaces:
```go
// Open for extension
type Cache interface {
    Get(ctx context.Context, key string) ([]byte, error)
    Set(ctx context.Context, key string, val []byte, ttl time.Duration) error
}

// Closed for modification - add new adapters without changing callers
type RedisCache struct { /* implements Cache */ }
type MemoryCache struct { /* implements Cache */ }
```

### Liskov Substitution
All cache implementations are interchangeable:
```go
func NewService(cache cache.Cache) *Service {
    // Works with RedisCache, MemoryCache, or any cache.Cache impl
}
```

### Interface Segregation
Small, focused interfaces:
```go
// Rather than one large interface:
type Logger interface { Debug, Info, Warn, Error }
type MetricsCollector interface { Counter, Gauge, Histogram }
```

### Dependency Inversion
High-level modules depend on abstractions:
```go
// Domain defines the interface
type Repository interface { Save, Find, Delete }

// Infrastructure implements it
type PostgresRepository struct { /* ... */ }
```

## DDD Bounded Contexts

### Infrastructure Context
Core utilities used across all services:
- `logctx`, `ringbuffer`, `waitfor`, `registry`

### Observability Context
Monitoring and tracing:
- `metrics`, `tracing`, `alerting`, `health`

### Data Context
Persistence and caching:
- `cache`, `db`, `repository`, `storage`

### Security Context
Authentication and authorization:
- `auth`, `oauth2`, `secrets`, `cors`

## Testing Strategy (TDD + BDD)

### Unit Tests
- All packages have `_test.go` files
- Test naming: `Test<Subject>_<Method>_<Condition>_<Expected>`
- 80%+ coverage target

### Integration Tests
- Adapters tested against real services (Redis, Postgres)
- Use testcontainers for isolation

### Property-Based Tests (proptest)
- Generic packages tested with random inputs
- Invariant testing for `ringbuffer`, `registry`

## Contract Testing

External integrations verified with contracts:
- API clients have contract tests
- Use `pact` or similar for consumer-driven contracts

## ADRs

| Number | Title | Status |
|--------|-------|--------|
| ADR-001 | Package Independence Policy | Accepted |
| ADR-002 | Interface-First Design | Accepted |
| ADR-003 | Zero Dependency Core | Accepted |
| ADR-004 | Context-Propagation Logging | Accepted |
| ADR-005 | Structured Metrics with OTEL | Accepted |

See `docs/adr/` for full decision records.

## Migration Path

### Phase 1: Document
- [x] Create ARCHITECTURE.md (this file)
- [x] Add ADRs for existing decisions
- [ ] Document all public APIs

### Phase 2: Refactor
- [ ] Extract interfaces for all adapters
- [ ] Move toward port/adapter isolation
- [ ] Add explicit domain models

### Phase 3: Quality
- [ ] Increase test coverage to 80%+
- [ ] Add property-based tests
- [ ] Add mutation testing

## References

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [DDD Reference](https://domainlanguage.com/ddd/reference/)
- See `xDD_METHODOLOGIES.md` for full methodology reference
