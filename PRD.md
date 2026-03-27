# Product Requirements Document — phenotype-go-kit

## Product Vision

`phenotype-go-kit` (`github.com/KooshaPari/phenotype-go-kit`) is a Go infrastructure toolkit providing reusable, production-grade packages extracted from the Phenotype ecosystem. Each package is independently importable with minimal cross-package dependencies. The module targets Go 1.25+ services following Hexagonal Architecture with SOLID, GRASP, and Law of Demeter principles.

---

## Epics

### E1 — Authentication and Authorization

**Goal**: Give service developers a complete, drop-in JWT and API-key authentication layer with HTTP middleware that works with any `net/http`-compatible router.

**User Stories**:

- E1.1: As a service developer, I want JWT generation and validation with both HMAC-SHA256 and RSA-256 signing so I can choose the appropriate algorithm per deployment context.
  - AC: `auth/jwt.go` exports `JWTValidator` with `GenerateTokenPair(userID, email string, roles []string) (access, refresh string, err error)`.
  - AC: `ValidateAccessToken(tokenString string) (*Claims, error)` rejects expired, malformed, and misscoped tokens with typed errors.
  - AC: Algorithm is selected from config; missing config key fails loudly.
  - Source: `auth/jwt.go`, `go.mod` — `github.com/golang-jwt/jwt/v5 v5.3.1`

- E1.2: As a service developer, I want an HTTP middleware that validates Bearer tokens and injects user identity (ID, email, roles) into request context so auth state is available to all downstream handlers without repetition.
  - AC: Middleware returns HTTP 401 with JSON error body when token is missing or invalid.
  - AC: Claims are accessible via `GetUserID(ctx)`, `GetUserEmail(ctx)`, `GetUserRoles(ctx)` helpers.
  - Source: `auth/adapter/`

- E1.3: As a service developer, I want role-based access middleware (`RequireRole`) so I can protect individual endpoints with fine-grained permission checks.
  - AC: `RequireRole(roles ...string)` returns HTTP 403 when authenticated user lacks all specified roles.
  - Source: `auth/adapter/`

- E1.4: As a service developer, I want an API key manager (generate, hash for storage, validate, revoke) so services can support machine-to-machine authentication.
  - AC: Generated key is shown once; only the bcrypt/SHA256 hash is persisted.
  - AC: `ValidateKey(rawKey string) (*APIKeyMetadata, bool)` returns false for revoked keys without exposing timing.
  - Source: `auth/`

- E1.5: As a service developer, I want context helpers (`GetUserID`, `GetUserEmail`, `GetUserRoles`) so handler code can extract auth claims without type assertions against raw context keys.
  - AC: Helpers panic with a descriptive message when called without prior auth middleware (fail loudly per governance).
  - Source: `auth/adapter/`

---

### E2 — Observability

**Goal**: Provide a consistent, Prometheus-native and OpenTelemetry-native observability stack composable with any HTTP service.

**User Stories**:

- E2.1: As a platform engineer, I want Prometheus metrics (HTTP request count, duration, response size; DB query duration and errors; job queue depth, processing time, retries) registered under the `phenotype` namespace so dashboards have a consistent label schema across all services.
  - AC: All metrics registered at package init under `phenotype_` prefix.
  - AC: `MetricsMiddleware` skips `/health` and `/ready` endpoints from recording.
  - Source: `metrics/`, `go.mod` — `github.com/prometheus/client_golang v1.22.0`

- E2.2: As a platform engineer, I want OpenTelemetry tracing initialized via a single call (`tracing.Init`) with OTLP/gRPC export so services share consistent trace propagation.
  - AC: `tracing.Init(ctx, serviceName, otlpEndpoint string) (func(), error)` returns a shutdown function.
  - AC: Missing endpoint returns error immediately (no silent no-op).
  - Source: `tracing/`, `go.mod` — `go.opentelemetry.io/otel v1.40.0`, `otlptracegrpc v1.24.0`

- E2.3: As a service developer, I want structured logging with `log/slog` including log rotation support and a gRPC interceptor so all log output is machine-parseable and consistent.
  - Source: `logging/`, `logctx/logctx.go`

- E2.4: As a platform engineer, I want a health check aggregator with liveness and readiness HTTP handlers returning JSON so Kubernetes probes work out of the box.
  - AC: `GET /health/live` returns 200 always (liveness). `GET /health/ready` returns 503 if any registered check fails.
  - AC: Response body is `{"status":"ok","checks":{"db":"ok","redis":"degraded",...}}`.
  - Source: `health/`

- E2.5: As a service developer, I want `logctx.WithLogger(ctx, logger)` and `logctx.From(ctx)` so logger propagates through context without global state.
  - AC: `logctx.From(ctx)` panics if no logger in context — programmer error, fail loudly.
  - Source: `logctx/logctx.go`

---

### E3 — Resilience Patterns

**Goal**: Provide circuit breaker, retry, rate limiter, and bounded buffer primitives so services handle downstream failures and load spikes without custom implementations.

**User Stories**:

- E3.1: As a service developer, I want a circuit breaker with configurable failure threshold, success threshold, and timeout and three states (Closed, Open, HalfOpen) so downstream failures stop cascading.
  - AC: `Breaker.Execute(fn func() error) error` returns `ErrCircuitOpen` when breaker is open.
  - AC: State transitions: Closed → Open after N failures; Open → HalfOpen after timeout; HalfOpen → Closed after M successes.
  - Source: `circuit/breaker.go`

- E3.2: As a service developer, I want a retry utility with exponential backoff, configurable jitter, max attempts, and context deadline propagation so transient errors are handled uniformly.
  - AC: `retry.Do(ctx, fn, opts...)` respects `context.Done()` and returns the context error on cancellation.
  - AC: Maximum retry interval is configurable; default is capped at 30s.
  - Source: `retry/`

- E3.3: As a service developer, I want a token-bucket rate limiter keyed on API key, auth header, or IP, with HTTP middleware, so services can shed load without external Redis dependency.
  - AC: `RateLimiter.Allow(key string) bool` uses in-memory token bucket.
  - AC: HTTP middleware returns 429 with `Retry-After` header when rate limit is exceeded.
  - Source: `ratelimit/`

- E3.4: As a service developer, I want a generic ring buffer with fixed capacity and thread safety so bounded in-memory queues can be built without external deps.
  - AC: `ringbuffer.New[T](capacity int)` creates a buffer; `Push(v T)` overwrites oldest on full.
  - AC: `GetAll()` returns items in insertion order, oldest first.
  - AC: `Len()` and `Cap()` are O(1).
  - Source: `ringbuffer/ringbuffer.go` (confirmed real implementation)

- E3.5: As a service developer, I want `waitfor.Poll(ctx, condition func() bool, opts...)` with exponential backoff and testable clock injection so polling loops work in tests without `time.Sleep`.
  - AC: Uses `github.com/coder/quartz` for clock abstraction enabling deterministic test control.
  - AC: Returns `context.DeadlineExceeded` when deadline is reached.
  - Source: `waitfor/waitfor.go`, `go.mod` — `github.com/coder/quartz v0.1.2`

---

### E4 — Data Layer Abstractions

**Goal**: Provide storage-backend-agnostic adapters for cache, object storage, and secrets so services swap backends via config without code changes.

**User Stories**:

- E4.1: As a service developer, I want a Redis cache adapter with TTL-based invalidation and a service layer so read-heavy services can cache without writing custom Redis wrappers.
  - AC: `cache.NewRedisAdapter(client *redis.Client)` returns a `CacheAdapter` interface implementation.
  - AC: Cache miss returns a sentinel error (`cache.ErrCacheMiss`), not `nil`.
  - Source: `cache/redis.go`, `cache/adapter/`, `go.mod` — `github.com/redis/go-redis/v9 v9.18.0`

- E4.2: As a service developer, I want an object storage abstraction supporting AWS S3 and Google Cloud Storage behind a unified interface so the storage backend is swappable via config.
  - AC: `storage.FileStore` interface with `Put(ctx, key, reader)`, `Get(ctx, key)`, `Delete(ctx, key)`.
  - AC: `storage.NewS3Store(cfg)` and `storage.NewGCSStore(cfg)` both implement `FileStore`.
  - Source: `storage/`, `go.mod` — AWS SDK v2, `cloud.google.com/go/storage v1.43.0`

- E4.3: As a service developer, I want a secrets provider port so credentials are fetched from secret stores and never hardcoded.
  - AC: `secrets.Provider` interface with `GetSecret(ctx, name string) (string, error)`.
  - Source: `secrets/`

- E4.4: As a service developer, I want a generic thread-safe key-value registry with ref counting, owner tracking, and change hooks so service registries can be built without reinventing locking logic.
  - AC: `registry.New[K comparable, V any]()` creates typed registry.
  - AC: `Register(key K, value V, owner string)` is goroutine-safe.
  - AC: `OnChange(hook func(key K, old, new V))` notifies on every update.
  - Source: `registry/registry.go` (confirmed real implementation)

---

### E5 — Application Infrastructure

**Goal**: Provide config, OAuth2, CORS, event bus, and validation utilities that cover the common application infrastructure needs of Phenotype services.

**User Stories**:

- E5.1: As a service developer, I want a configuration loader backed by Viper (file, env, flags) so service config is consistent and overrideable for all environments.
  - AC: `config.Load(path string) (*Config, error)` loads YAML config; env vars override file with `PHENOTYPE_` prefix.
  - Source: `config/`, `go.mod` — `github.com/spf13/viper v1.21.0`

- E5.2: As a service developer, I want an OAuth2 provider integration so services can delegate identity to external providers with minimal boilerplate.
  - AC: `oauth2.NewProvider(cfg)` returns a handler for authorization code flow.
  - Source: `oauth2/`

- E5.3: As a service developer, I want a CORS middleware configurable per-route so API servers can safely expose endpoints to browser clients.
  - AC: `cors.New(opts)` returns `net/http`-compatible middleware.
  - Source: `cors/`, `go.mod` — `github.com/go-chi/cors v1.2.1` (via go-chi)

- E5.4: As a service developer, I want an in-process event bus for publishing and subscribing to domain events so bounded contexts communicate without direct coupling.
  - AC: `bus.Publish(ctx, event)` and `bus.Subscribe(topic, handler)` with typed generics.
  - Source: `bus/`

- E5.5: As a service developer, I want a validation package so request DTOs can be validated with consistent, structured error responses.
  - AC: Validation errors return a map of field name to error message for HTTP 422 responses.
  - Source: `validation/`

---

### E6 — Hexagonal Architecture Contracts

**Goal**: Define port interfaces and domain model types in a `contracts/` package separate from implementations so service cores depend only on abstractions.

**User Stories**:

- E6.1: As an architect, I want port interfaces (inbound and outbound) in `contracts/ports` so service cores import no implementation packages.
  - AC: `contracts/ports` contains only `interface` definitions and no `struct` implementations.
  - Source: `contracts/`

- E6.2: As an architect, I want domain model types in `contracts/models` so services share canonical data structures without circular imports.
  - Source: `contracts/`

- E6.3: As an architect, I want a plugin system interface in `contracts/plugins` so runtime-loadable extensions conform to a known contract.
  - Source: `contracts/`, `plugins/`

---

## Acceptance Criteria Matrix

| Epic | Verifiable Criterion |
|------|---------------------|
| E1 | `JWTValidator.GenerateTokenPair` returns signed token; `ValidateAccessToken` rejects expired tokens |
| E1 | API key store persists hash only; validate returns false for revoked keys |
| E2 | `MetricsMiddleware` skips `/health`; all `phenotype_*` counters exist after first request |
| E2 | `tracing.Init` with invalid endpoint returns error (not no-op) |
| E2 | `logctx.From(ctx)` panics without prior `WithLogger` call |
| E3 | `Breaker.Execute` returns `ErrCircuitOpen` when open |
| E3 | `waitfor.Poll` uses quartz clock; test advances clock deterministically |
| E3 | `ringbuffer.Push` on full buffer evicts oldest; `GetAll` returns oldest-first |
| E4 | `cache.Get` on miss returns `ErrCacheMiss`, not nil |
| E4 | S3 and GCS adapters implement identical `FileStore` interface |
| E6 | `contracts/` directory contains only interfaces and plain structs |

---

## Non-Goals

- This module does not define application business logic or domain rules specific to any Phenotype service.
- It does not provide a full HTTP framework — it provides middleware composable with `go-chi/chi v5`.
- It does not implement distributed circuit breaking or rate limiting across pods (in-process only).
- It does not include UI or CLI tooling.
