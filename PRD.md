# Product Requirements Document — phenotype-go-kit

## Overview

phenotype-go-kit is a Go module (`github.com/KooshaPari/phenotype-go-kit`) providing reusable, production-grade infrastructure packages extracted from the Phenotype ecosystem. Each package is independently importable with minimal cross-package dependencies. The module targets Go 1.25+ services that follow Hexagonal Architecture with SOLID, GRASP, and Law of Demeter principles.

---

## Epics and User Stories

### E1 — Authentication and Authorization

**E1.1** As a service developer, I want JWT generation and validation with both HMAC-SHA256 and RSA-256 signing so I can choose the appropriate algorithm per deployment context.

**E1.2** As a service developer, I want an HTTP middleware that validates Bearer tokens and injects user identity (ID, email, roles) into request context so auth state is available to all downstream handlers without repetition.

**E1.3** As a service developer, I want role-based access middleware (`RequireRole`) so I can protect individual endpoints with fine-grained permission checks.

**E1.4** As a service developer, I want an API key manager (generate, hash for storage, validate, revoke) so services can support machine-to-machine authentication without user sessions.

**E1.5** As a service developer, I want context helpers (`GetUserID`, `GetUserEmail`, `GetUserRoles`) so handler code can extract auth claims without type assertions against raw context keys.

---

### E2 — Observability

**E2.1** As a platform engineer, I want Prometheus metrics (HTTP request count, duration, response size; DB query duration and errors; job queue depth, processing time, retries) registered under the `phenotype` namespace so dashboards have a consistent label schema across all services.

**E2.2** As a platform engineer, I want an HTTP middleware (`MetricsMiddleware`) that records per-request metrics automatically, excluding health endpoints, so instrumentation requires no per-handler code.

**E2.3** As a platform engineer, I want OpenTelemetry tracing initialized via a single call (`tracing.Init`) with OTLP/gRPC export, so services share a consistent trace propagation setup.

**E2.4** As a service developer, I want structured logging with `log/slog` including log rotation support and a gRPC interceptor for request logging so all log output is machine-parseable and consistent.

**E2.5** As a platform engineer, I want a health check aggregator with liveness and readiness HTTP handlers that returns JSON status across all registered component checks, so Kubernetes probes work out of the box.

---

### E3 — Resilience Patterns

**E3.1** As a service developer, I want a circuit breaker with configurable failure threshold, success threshold, and timeout, and three states (Closed, Open, HalfOpen), so downstream failures stop cascading.

**E3.2** As a service developer, I want a retry utility with exponential backoff, configurable jitter, max attempts, and per-call deadline propagation, so transient errors are handled uniformly.

**E3.3** As a service developer, I want a token-bucket rate limiter keyed on API key, auth header, or IP, with block/unblock controls and HTTP middleware, so services can shed load without external dependencies.

**E3.4** As a service developer, I want a ring buffer with fixed capacity and thread safety so bounded in-memory queues can be built without generics complexity.

---

### E4 — Data Layer Abstractions

**E4.1** As a service developer, I want a database connection pool abstraction, query builder utilities, and index helpers so services share consistent DB setup patterns.

**E4.2** As a service developer, I want a Redis cache adapter with TTL-based invalidation and a cache service layer so read-heavy services can cache without writing custom Redis wrappers.

**E4.3** As a service developer, I want an object storage abstraction supporting AWS S3 and Google Cloud Storage behind a unified interface so storage backend can be swapped without service-layer changes.

**E4.4** As a service developer, I want a secrets provider port and adapter so credentials are fetched from secret stores (not hardcoded), satisfying the 12-factor configuration principle.

---

### E5 — Application Infrastructure

**E5.1** As a service developer, I want a configuration loader backed by Viper (file, env, flags) so service config is consistent and overrideable for all environments.

**E5.2** As a service developer, I want an OAuth2 provider integration so services can delegate identity to external providers (Google, GitHub, etc.) with minimal boilerplate.

**E5.3** As a service developer, I want a CORS middleware configurable per-route so API servers can safely expose endpoints to browser clients.

**E5.4** As a service developer, I want an event bus (local, in-process) for publishing and subscribing to domain events, so bounded contexts communicate without direct coupling.

**E5.5** As a service developer, I want a validation package so request DTOs can be validated with consistent error structures.

---

### E6 — Hexagonal Architecture Contracts

**E6.1** As an architect, I want port interfaces (inbound and outbound) defined in the `contracts/ports` package separate from implementations, so service cores depend only on abstractions.

**E6.2** As an architect, I want domain model types in `contracts/models` so services share canonical data structures without circular imports.

**E6.3** As an architect, I want a plugin system interface in `contracts/plugins` so runtime-loadable extensions conform to a known contract.

---

## Acceptance Criteria

| Epic | Criterion |
|------|-----------|
| E1 | `JWTValidator.GenerateTokenPair` returns HS256 or RS256 tokens depending on config; `ValidateAccessToken` rejects expired and misscoped tokens |
| E1 | `APIKeyManager` stores only the hash, never the raw key; `ValidateKey` returns the stored metadata |
| E2 | `NewMetrics()` registers all counters/histograms without panic; `MetricsMiddleware` skips `/health` and `/ready` |
| E2 | `tracing.Init` exports spans to a configurable OTLP endpoint |
| E3 | `Breaker.Execute` returns `ErrCircuitOpen` when the breaker is open; transitions to HalfOpen after timeout |
| E3 | `retry.Do` respects `context.Done` and returns early on cancellation |
| E3 | `RateLimiter.Allow` uses token bucket; blocked keys return false immediately |
| E4 | Cache adapter reads from Redis; on miss returns the appropriate sentinel error |
| E4 | Storage adapters (S3, GCS) implement the same file interface |
| E6 | No implementation code in `contracts/`; all types are interfaces or plain data structs |

---

## Non-Goals

- This module does not define application business logic or domain rules.
- It does not provide a full HTTP framework — it provides middleware composable with `go-chi/chi`.
- It does not implement distributed circuit breaking or rate limiting across pods (local only, Redis integration is a future extension).
- It does not include UI or CLI tooling.
