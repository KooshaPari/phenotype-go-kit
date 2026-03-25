# PRD - phenotype-go-kit

## E1: Context-Scoped Logging

### E1.1: Logger Injection
As an application developer, I inject a `slog.Logger` into `context.Context` and retrieve it downstream so logging is consistent across request boundaries.

**Acceptance**: `WithLogger` stores logger; `From` retrieves it; missing logger panics.

## E2: Ring Buffer

### E2.1: Fixed-Capacity Circular Buffer
As a systems developer, I use a generic ring buffer that overwrites oldest entries when full, providing O(1) push and oldest-first iteration.

**Acceptance**: `Push` overwrites oldest at capacity; `GetAll` returns oldest-first; `Len`/`Cap` accurate.

## E3: Polling with Backoff

### E3.1: WaitFor
As a service developer, I poll a condition with exponential backoff, configurable timeout, and testable clocks.

**Acceptance**: Polls until true/error/timeout; supports `quartz.Clock` for deterministic tests; returns `ErrTimedOut` on timeout.

## E4: Thread-Safe Registry

### E4.1: Key-Value Registry with Ref Counting
As a platform developer, I use a generic registry with owner-scoped lifecycle, ref counting, and change hooks.

**Acceptance**: Multiple owners per key; entry removed on last unregister; Hook interface for observability; thread-safe via `sync.RWMutex`.
