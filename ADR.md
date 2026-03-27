# Architecture Decision Records — phenotype-go-kit

---

## ADR-001 — Hexagonal Architecture as the Primary Structural Pattern

**Status:** Accepted
**Date:** 2025-Q4

### Context

phenotype-go-kit is consumed by multiple Phenotype-ecosystem services. Each service needs infrastructure utilities (auth, caching, metrics, etc.) without importing business logic from other services. Early versions placed implementation code directly in consumer services, causing duplication and tight coupling to specific vendors (e.g., Redis client calls scattered throughout handlers).

### Decision

The module is organized around Hexagonal Architecture (Ports and Adapters):

- `contracts/ports/inbound` — driving ports (use cases, command handlers)
- `contracts/ports/outbound` — driven ports (repositories, caches, event publishers)
- Each concrete package (`cache/`, `auth/`, `storage/`) contains secondary adapters implementing outbound ports
- Application-layer services in `cache/service/` implement inbound ports using outbound port interfaces

### Consequences

- Consumer services depend on `contracts/ports` interfaces, not concrete adapters; swapping backends requires zero application-layer changes.
- Each package is independently importable — a service that only needs `auth/` does not pull in `storage/`.
- New adapters (e.g., Memcached instead of Redis) are added without touching the contracts package.
- Increased initial directory depth compared to flat package layout.

---

## ADR-002 — xDD Methodologies Reference (TDD, BDD, DDD, ATDD)

**Status:** Reference
**Date:** 2025-Q4

### Context

The Phenotype organization mandates cross-methodology development discipline (xDD) across all modules.

### Decision

phenotype-go-kit adopts the following practices:

- **TDD**: All packages require `_test.go` files before implementation stabilizes. `ringbuffer/ringbuffer_test.go` is the canonical example.
- **DDD**: The `contracts/models` package captures domain language shared across bounded contexts.
- **BDD/ATDD**: Acceptance criteria in `PRD.md` drive integration tests. Test function names reference FR IDs.

### Consequences

- New packages require a test file to be created before or alongside the implementation.
- `contracts/models` must not import any adapter package to preserve the domain layer's independence.

---

## ADR-003 — Top-Level Directory Structure (Package-per-Feature)

**Status:** Accepted
**Date:** 2025-Q4

### Context

Go module organization choices have long-term import-path consequences. Two common patterns are flat (all code in one or few packages) and package-per-feature.

### Decision

Package-per-feature is used. Each infrastructure concern is a top-level directory with its own package:

```
auth/       circuit/    cache/      config/
db/         health/     logging/    metrics/
oauth2/     ratelimit/  retry/      ringbuffer/
storage/    tracing/    validation/ ...
```

Sub-packages within a feature directory (`cache/adapter/`, `cache/service/`) follow the hexagonal layering: adapter contains secondary adapters, service contains application-layer handlers.

### Consequences

- Import paths are descriptive: `github.com/KooshaPari/phenotype-go-kit/circuit`, not `github.com/KooshaPari/phenotype-go-kit/internal/resilience/circuit`.
- Consumers import exactly the packages they need; unused packages add no binary size.
- The `contracts/` directory is a special case — it contains only interfaces and models, never concrete implementations.

---

## ADR-004 — Plugin System and Extensibility

**Status:** Accepted
**Date:** 2025-Q4

### Context

Some capabilities (AI embeddings, external provider integrations) require runtime-selectable implementations. Hardcoding a single provider creates vendor lock-in.

### Decision

A plugin interface is defined in `contracts/plugins`. Concrete plugins live in the `plugins/` directory. The `embeddings` sub-package demonstrates the pattern with OpenAI and Ollama providers behind a common interface.

### Consequences

- New providers are added by implementing the plugin interface; no changes to consuming code.
- Plugin registration is done at service startup via dependency injection.

---

## ADR-005 — AI Embeddings Plugin System

**Status:** Accepted
**Date:** 2025-Q4

### Context

Vector embeddings are needed for semantic search features across Phenotype services. Multiple embedding providers exist (OpenAI, Ollama for local inference).

### Decision

An embeddings plugin under `plugins/embeddings/` implements the plugin interface from `contracts/plugins`. Both OpenAI and Ollama providers are supported. Selection is via configuration, not compile-time flags.

### Consequences

- Switching from OpenAI to Ollama (or adding a third provider) requires only a config change.
- The embeddings plugin is an optional import; services that do not need embeddings incur no dependency on OpenAI or Ollama SDKs.

---

## ADR-006 — Design Principles: SOLID, GRASP, Law of Demeter

**Status:** Accepted
**Date:** 2025-Q4

### Context

Without explicit guidance, infrastructure packages accumulate antipatterns: large structs with many responsibilities, deep method chains (train wrecks), and tight coupling between layers.

### Decision

All packages in this module SHALL adhere to:

| Principle | Application |
|-----------|-------------|
| SRP | Each struct/function has one reason to change. `JWTValidator` handles only token lifecycle; `APIKeyManager` handles only API key lifecycle. |
| OCP | New signing algorithms are added via config (`PrivateKey` nil/non-nil), not by modifying `signClaims`. |
| LSP | `Checker` implementations (`DatabaseChecker`, `RedisChecker`, `ComponentChecker`) are substitutable; `HealthChecker` only calls `Check(ctx)`. |
| ISP | Port interfaces in `contracts/ports` are small and focused. `cache.go`, `db.go`, `secrets.go` are separate files, not one giant `ports.go`. |
| DIP | Application services depend on port interfaces; adapters depend on external SDKs. The dependency arrow always points inward. |
| GRASP (Information Expert) | `TokenClaims` holds all JWT claim fields; validation of claim scope is in `JWTValidator`, not the caller. |
| Law of Demeter | Context accessors (`GetUserID`, `GetUserEmail`, `GetUserRoles`) eliminate train-wreck access to context values. |

### Consequences

- Code review must flag ISP violations (adding methods to an existing interface when a new interface would be more appropriate).
- Complexity ratchet enforcement limits function cyclomatic complexity to 10 and cognitive complexity to 15.

---

## ADR-007 — Standard Library Logging via `log/slog`

**Status:** Accepted
**Date:** 2025-Q4

### Context

Third-party logging libraries (zap, logrus) require consumers to adopt the same library. Go 1.21 introduced `log/slog` as the standard structured logging interface.

### Decision

All packages in phenotype-go-kit use `log/slog` exclusively for internal logging. The `logging/` package provides structured configuration, rotation wrappers, and a gRPC interceptor around `slog`.

### Consequences

- No consumer is forced to import a third-party logging library to work with this module.
- `slog.Default()` is used as the logger in packages that do not receive an injected logger; consumers can replace the default logger at program startup via `slog.SetDefault`.

---

## ADR-008 — Token Bucket for Rate Limiting (In-Process Only)

**Status:** Accepted
**Date:** 2025-Q4

### Context

Rate limiting can be implemented locally (in-process, per-pod) or distributed (shared Redis counter). Distributed rate limiting requires Redis availability as a hard dependency for the limiter itself.

### Decision

`ratelimit.RateLimiter` uses the token bucket algorithm, in-process only. A `DistributedRateLimiter` stub exists but delegates to the local implementation pending Redis integration. The `Config.RequestsPerSecond` and `Config.BurstSize` parameters control limits.

### Consequences

- Services running multiple pods will apply limits per-pod, not globally. For global limits, the Redis-backed variant must be completed (tracked as a future work item).
- No Redis dependency is introduced into the rate limiter package itself.

---

## ADR-009 — Prometheus as the Sole Metrics Backend

**Status:** Accepted
**Date:** 2025-Q4

### Context

Multiple metrics backends exist (Prometheus, Datadog, CloudWatch). Supporting all via abstraction adds complexity; the Phenotype platform standardizes on Prometheus for self-hosted observability.

### Decision

`metrics/collector.go` uses `github.com/prometheus/client_golang` directly. The `promauto` sub-package is used for registration at construction time.

### Consequences

- Services importing `phenotype-go-kit/metrics` pull in the Prometheus client library.
- Changing the metrics backend would require a new adapter package rather than modifying `collector.go`. This is acceptable given the platform-wide Prometheus standard.
- Datadog/CloudWatch configs present in `config/` are agent-side configuration files, not Go code — they do not create a Go dependency.
