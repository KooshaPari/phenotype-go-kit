# Functional Requirements — phenotype-go-kit

FR IDs follow the pattern `FR-{CAT}-{NNN}` where CAT is the package category abbreviation.

---

## FR-AUTH — Authentication and Authorization (`auth/`)

| ID | SHALL Statement | Traces To | Status |
|----|-----------------|-----------|--------|
| FR-AUTH-001 | The `JWTValidator` SHALL generate access tokens signed with HS256 (HMAC-SHA256) when `JWTConfig.PrivateKey` is nil | E1.1 | Implemented |
| FR-AUTH-002 | The `JWTValidator` SHALL generate access tokens signed with RS256 (RSA-256) when `JWTConfig.PrivateKey` is set | E1.1 | Implemented |
| FR-AUTH-003 | `GenerateTokenPair` SHALL return a `TokenPair` containing access token, refresh token, expiry seconds, and token type "Bearer" | E1.1 | Implemented |
| FR-AUTH-004 | `ValidateAccessToken` SHALL reject tokens with `Scope != "access"` with `ErrInvalidClaims` | E1.1 | Implemented |
| FR-AUTH-005 | `ValidateRefreshToken` SHALL reject tokens with `Scope != "refresh"` with `ErrInvalidClaims` | E1.1 | Implemented |
| FR-AUTH-006 | The JWT middleware SHALL extract the Bearer token from the `Authorization` header and return HTTP 401 if absent or malformed | E1.2 | Implemented |
| FR-AUTH-007 | The JWT middleware SHALL inject `user_id`, `user_email`, and `user_roles` values into `context.Context` on successful validation | E1.2 | Implemented |
| FR-AUTH-008 | `RequireRole` middleware SHALL return HTTP 403 when the authenticated user lacks any of the required roles | E1.3 | Implemented |
| FR-AUTH-009 | `GenerateAPIKey` SHALL produce a URL-safe base64 string of 32 random bytes, optionally prefixed | E1.4 | Implemented |
| FR-AUTH-010 | `HashAPIKey` SHALL produce a SHA-256 hash of the raw key, base64-encoded, suitable for storage | E1.4 | Implemented |
| FR-AUTH-011 | `APIKeyManager.CreateKey` SHALL store only the hash, not the raw key | E1.4 | Implemented |
| FR-AUTH-012 | `APIKeyManager.RevokeKey` SHALL delete the key by ID and return an error if the ID is not found | E1.4 | Implemented |
| FR-AUTH-013 | `GetUserID`, `GetUserEmail`, `GetUserRoles` SHALL return zero values (empty string / nil slice) when the corresponding context key is absent, not panic | E1.5 | Implemented |

---

## FR-METR — Metrics (`metrics/`)

| ID | SHALL Statement | Traces To | Status |
|----|-----------------|-----------|--------|
| FR-METR-001 | `NewMetrics` SHALL register HTTP request count, duration, and response size histograms under the `phenotype_http` namespace | E2.1 | Implemented |
| FR-METR-002 | `NewMetrics` SHALL register job queue depth, processing time, and retry counters under the `phenotype_jobs` namespace | E2.1 | Implemented |
| FR-METR-003 | `NewMetrics` SHALL register DB query duration and error counters under the `phenotype_db` namespace | E2.1 | Implemented |
| FR-METR-004 | `MetricsMiddleware` SHALL record method, path, and HTTP status code for every non-health request | E2.2 | Implemented |
| FR-METR-005 | `MetricsMiddleware` SHALL skip recording for paths `/health` and `/ready` | E2.2 | Implemented |
| FR-METR-006 | `RecordBusinessMetric` SHALL lazily register a new `CounterVec` for previously unseen metric names | E2.1 | Implemented |

---

## FR-TRAC — Tracing (`tracing/`)

| ID | SHALL Statement | Traces To | Status |
|----|-----------------|-----------|--------|
| FR-TRAC-001 | The tracing package SHALL initialize an OpenTelemetry `TracerProvider` with OTLP/gRPC export | E2.3 | Implemented |
| FR-TRAC-002 | The tracer endpoint SHALL be configurable at runtime (not hardcoded) | E2.3 | Implemented |

---

## FR-LOG — Logging (`logging/`)

| ID | SHALL Statement | Traces To | Status |
|----|-----------------|-----------|--------|
| FR-LOG-001 | The logging package SHALL provide structured output via `log/slog` | E2.4 | Implemented |
| FR-LOG-002 | The logging package SHALL provide log file rotation (size/time-based) | E2.4 | Implemented |
| FR-LOG-003 | A gRPC interceptor SHALL log method name, duration, and error for every gRPC call | E2.4 | Implemented |

---

## FR-HLTH — Health Checks (`health/`)

| ID | SHALL Statement | Traces To | Status |
|----|-----------------|-----------|--------|
| FR-HLTH-001 | `HealthChecker.Register` SHALL accept any type implementing the `Checker` interface | E2.5 | Implemented |
| FR-HLTH-002 | `RunAll` SHALL enforce the configured timeout per individual check | E2.5 | Implemented |
| FR-HLTH-003 | `ReadinessHandler` SHALL return HTTP 503 and JSON `{"status":"unhealthy"}` if any check status is `"unhealthy"` | E2.5 | Implemented |
| FR-HLTH-004 | `ReadinessHandler` SHALL return HTTP 200 and JSON `{"status":"healthy"}` when all checks pass | E2.5 | Implemented |
| FR-HLTH-005 | `LivenessHandler` SHALL always return HTTP 200 "OK" | E2.5 | Implemented |

---

## FR-CIRC — Circuit Breaker (`circuit/`)

| ID | SHALL Statement | Traces To | Status |
|----|-----------------|-----------|--------|
| FR-CIRC-001 | `Breaker.Execute` SHALL return `ErrCircuitOpen` immediately when the breaker is in the Open state and the timeout has not elapsed | E3.1 | Implemented |
| FR-CIRC-002 | `Breaker` SHALL transition from Open to HalfOpen after `Config.Timeout` has elapsed | E3.1 | Implemented |
| FR-CIRC-003 | `Breaker` SHALL transition from HalfOpen to Closed after `Config.SuccessThreshold` consecutive successes | E3.1 | Implemented |
| FR-CIRC-004 | `Breaker` SHALL transition from HalfOpen to Open on any single failure | E3.1 | Implemented |
| FR-CIRC-005 | `Breaker` SHALL transition from Closed to Open after `Config.FailureThreshold` consecutive failures | E3.1 | Implemented |
| FR-CIRC-006 | `Breaker.Execute` SHALL respect `Config.RequestTimeout` and count a timeout as a failure | E3.1 | Implemented |
| FR-CIRC-007 | `MultiBreaker.Get` SHALL be safe for concurrent callers and return the same `Breaker` for the same name | E3.1 | Implemented |

---

## FR-RETR — Retry (`retry/`)

| ID | SHALL Statement | Traces To | Status |
|----|-----------------|-----------|--------|
| FR-RETR-001 | `retry.Do` SHALL attempt the function up to `Config.MaxAttempts` times | E3.2 | Implemented |
| FR-RETR-002 | `retry.Do` SHALL apply exponential backoff: `delay *= Config.Multiplier` after each failure | E3.2 | Implemented |
| FR-RETR-003 | `retry.Do` SHALL cap the delay at `Config.MaxDelay` | E3.2 | Implemented |
| FR-RETR-004 | When `Config.Jitter` is true, `retry.Do` SHALL apply +/- 25% random jitter to each wait | E3.2 | Implemented |
| FR-RETR-005 | `retry.Do` SHALL return `ctx.Err()` immediately when `ctx.Done()` fires during a wait | E3.2 | Implemented |
| FR-RETR-006 | `PermanentError` SHALL mark errors as non-retryable; callers can detect via `IsPermanent` | E3.2 | Implemented |

---

## FR-RLIM — Rate Limiting (`ratelimit/`)

| ID | SHALL Statement | Traces To | Status |
|----|-----------------|-----------|--------|
| FR-RLIM-001 | `RateLimiter.Allow` SHALL implement the token bucket algorithm with configurable `RequestsPerSecond` and `BurstSize` | E3.3 | Implemented |
| FR-RLIM-002 | `RateLimiter.Allow` SHALL return false immediately for blocked keys before their block expiry | E3.3 | Implemented |
| FR-RLIM-003 | `RateLimiter.Middleware` SHALL extract the rate limit key in priority order: X-API-Key header, Authorization header, remote IP | E3.3 | Implemented |
| FR-RLIM-004 | `RateLimiter.Middleware` SHALL return HTTP 429 with `Retry-After: 1` when the limit is exceeded | E3.3 | Implemented |
| FR-RLIM-005 | The limiter SHALL run a background cleanup goroutine to evict stale client entries after `Config.CleanupInterval` | E3.3 | Implemented |

---

## FR-RING — Ring Buffer (`ringbuffer/`)

| ID | SHALL Statement | Traces To | Status |
|----|-----------------|-----------|--------|
| FR-RING-001 | The ring buffer SHALL have a fixed capacity set at construction time | E3.4 | Implemented |
| FR-RING-002 | Write operations to a full ring buffer SHALL overwrite the oldest entry | E3.4 | Implemented |
| FR-RING-003 | All ring buffer operations SHALL be safe for concurrent use | E3.4 | Implemented |

---

## FR-CACHE — Cache (`cache/`)

| ID | SHALL Statement | Traces To | Status |
|----|-----------------|-----------|--------|
| FR-CACHE-001 | The Redis cache adapter SHALL set values with a caller-provided TTL | E4.2 | Implemented |
| FR-CACHE-002 | `invalidation.go` SHALL provide TTL-based cache invalidation logic | E4.2 | Implemented |
| FR-CACHE-003 | The cache service layer SHALL expose CQRS-style handlers for cache read/write operations | E4.2 | Implemented |

---

## FR-STOR — Storage (`storage/`)

| ID | SHALL Statement | Traces To | Status |
|----|-----------------|-----------|--------|
| FR-STOR-001 | The S3 adapter SHALL support upload, download, and delete using AWS SDK v2 | E4.3 | Implemented |
| FR-STOR-002 | The GCS adapter SHALL support upload, download, and delete using the Google Cloud Storage client | E4.3 | Implemented |
| FR-STOR-003 | Both storage adapters SHALL implement the same file storage interface | E4.3 | Implemented |

---

## FR-CONF — Configuration (`config/`)

| ID | SHALL Statement | Traces To | Status |
|----|-----------------|-----------|--------|
| FR-CONF-001 | Configuration SHALL be loadable from YAML files, environment variables, and command-line flags via Viper | E5.1 | Implemented |

---

## FR-CONT — Contracts (`contracts/`)

| ID | SHALL Statement | Traces To | Status |
|----|-----------------|-----------|--------|
| FR-CONT-001 | `contracts/ports/inbound` SHALL define driving port interfaces (UseCase, CommandHandler, QueryHandler) | E6.1 | Implemented |
| FR-CONT-002 | `contracts/ports/outbound` SHALL define driven port interfaces (Repository, Cache, EventPublisher, Secrets) | E6.1 | Implemented |
| FR-CONT-003 | `contracts/models` SHALL define domain event and model types shared across adapters | E6.2 | Implemented |
| FR-CONT-004 | No implementation code (non-interface, non-struct) SHALL exist in the `contracts/` tree | E6.1 | Implemented |

---

## FR-SEC — Secrets Management (`secrets/`)

| ID | SHALL Statement | Traces To | Status |
|----|-----------------|-----------|--------|
| FR-SEC-001 | `SecretsManager.Get(ctx, key)` SHALL return the secret value or `ErrSecretNotFound` when the key does not exist; it SHALL NOT return an empty string for missing secrets | E4.3 | Implemented |
| FR-SEC-002 | The backend SHALL be selected from `Config{Backend}` accepting values `"vault"`, `"aws"`, and `"env"` | E4.3 | Implemented |
| FR-SEC-003 | The Vault adapter SHALL authenticate via the Vault HTTP API token and retrieve secrets from the KV v2 secrets engine | E4.3 | Implemented |
| FR-SEC-004 | The AWS adapter SHALL use `aws-sdk-go-v2/service/secretsmanager` and resolve credentials from the default credential chain | E4.3 | Implemented |
| FR-SEC-005 | The env adapter SHALL read secrets from process environment variables and return `ErrSecretNotFound` for unset variables | E4.3 | Implemented |

---

## FR-WEBH — Webhook Infrastructure (`webhook/`)

| ID | SHALL Statement | Traces To | Status |
|----|-----------------|-----------|--------|
| FR-WEBH-001 | `Delivery.Send(ctx, event)` SHALL POST a JSON-serialized `WebhookEvent` to the subscriber URL with `Content-Type: application/json` | E6.1 | Implemented |
| FR-WEBH-002 | Each outgoing request SHALL include an `X-Phenotype-Signature` header with value `sha256=<hex>` computed over the raw request body using HMAC-SHA256 with the subscriber secret | E6.1 | Implemented |
| FR-WEBH-003 | Failed deliveries SHALL be retried up to `Config.MaxAttempts` times with exponential backoff | E6.1 | Implemented |
| FR-WEBH-004 | `Verify(secret, payload []byte, signature string) bool` SHALL perform constant-time HMAC-SHA256 comparison | E6.2 | Implemented |
| FR-WEBH-005 | `Verify` SHALL return false for signatures that do not carry the `sha256=` prefix | E6.2 | Implemented |

---

## FR-MIG — Schema Migrations (`migrations/`)

| ID | SHALL Statement | Traces To | Status |
|----|-----------------|-----------|--------|
| FR-MIG-001 | `Runner.Up(ctx)` SHALL apply all pending migrations in ascending version order | E7.1 | Implemented |
| FR-MIG-002 | `Runner.Down(ctx, steps int)` SHALL roll back the specified number of applied migrations in descending order | E7.1 | Implemented |
| FR-MIG-003 | Applied migration versions SHALL be tracked in a `schema_migrations` table; already-applied versions SHALL be skipped | E7.1 | Implemented |
| FR-MIG-004 | Each migration SHALL execute inside a database transaction; a failure SHALL roll back the transaction and halt the run with the migration version in the error | E7.1 | Implemented |

---

## FR-REPO — Repository Pattern (`repository/`)

| ID | SHALL Statement | Traces To | Status |
|----|-----------------|-----------|--------|
| FR-REPO-001 | `Repository[T, ID].FindByID(ctx, id)` SHALL return `ErrNotFound` when no row matches the given ID | E8.1 | Implemented |
| FR-REPO-002 | `Repository[T, ID].Save(ctx, entity)` SHALL upsert via INSERT ... ON CONFLICT UPDATE | E8.1 | Implemented |
| FR-REPO-003 | `Repository[T, ID].Delete(ctx, id)` SHALL be idempotent — deleting a non-existent ID SHALL NOT return an error | E8.1 | Implemented |
| FR-REPO-004 | All repository methods SHALL accept `context.Context` as the first argument for query timeout and tracing propagation | E8.1 | Implemented |

---

## FR-ALRT — Alerting Rules (`alerting/`)

| ID | SHALL Statement | Traces To | Status |
|----|-----------------|-----------|--------|
| FR-ALRT-001 | `Rule{Name, Metric, Operator, Threshold, Severity}` SHALL compare a numeric fact against the threshold using one of the operators: `gt`, `lt`, `gte`, `lte`, `eq` | E9.1 | Implemented |
| FR-ALRT-002 | `RuleEngine.Evaluate(ctx, facts map[string]float64)` SHALL return one `Alert` per violated rule | E9.1 | Implemented |
| FR-ALRT-003 | `Alert{Rule, Value, Severity, FiredAt}` SHALL be serializable to JSON | E9.1 | Implemented |
| FR-ALRT-004 | Facts missing from the provided map SHALL NOT cause a panic; the rule SHALL be evaluated as non-violated when the fact is absent | E9.1 | Implemented |

---

## FR-EMBD — AI Embeddings Plugin (`plugins/embeddings/`)

| ID | SHALL Statement | Traces To | Status |
|----|-----------------|-----------|--------|
| FR-EMBD-001 | `Registry.Register(name string, provider Provider)` SHALL register an embedding provider under the given name | E10.1 | Implemented |
| FR-EMBD-002 | `Registry.Get(name string)` SHALL return `ErrProviderNotFound` for unknown provider names | E10.1 | Implemented |
| FR-EMBD-003 | `Provider.Embed(ctx, texts []string) ([][]float64, error)` SHALL return one embedding vector per input text | E10.1 | Implemented |
| FR-EMBD-004 | The OpenAI provider SHALL call the `/v1/embeddings` API endpoint with configurable model name and API key | E10.1 | Implemented |
| FR-EMBD-005 | The Ollama provider SHALL call the local Ollama HTTP API at a configurable base URL with configurable model name | E10.1 | Implemented |

---

## FR-CI — CI/CD Pipeline Utilities (`ci/`)

| ID | SHALL Statement | Traces To | Status |
|----|-----------------|-----------|--------|
| FR-CI-001 | `Pipeline.AddStage(stage *Stage)` SHALL append a stage with `Name`, `Commands []string`, `Env map[string]string`, and `Timeout time.Duration` | E11.1 | Implemented |
| FR-CI-002 | `Pipeline.Run(ctx)` SHALL execute stages sequentially; the first stage failure SHALL halt the run and return the stage name and underlying error | E11.1 | Implemented |
| FR-CI-003 | A per-stage `Timeout` SHALL be enforced via context deadline derived from the parent context | E11.1 | Implemented |
| FR-CI-004 | Environment variables from `Stage.Env` SHALL be merged into the subprocess environment without overwriting inherited variables unless keys collide | E11.1 | Implemented |
