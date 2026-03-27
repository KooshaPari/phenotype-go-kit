# phenotype-go-kit — Functional Requirements

## FR-LOGCTX-001: Context-Scoped Logger Injection

| Requirement ID | Description | Verification |
|---|---|---|
| FR-LOGCTX-001.1 | Package `logctx` SHALL export function `WithLogger(ctx context.Context, logger *slog.Logger) context.Context` | Function signature matches; returns new context with logger attached |
| FR-LOGCTX-001.2 | Function `From(ctx context.Context) *slog.Logger` SHALL retrieve logger from context | Function retrieves value for loggerKey; type asserts to *slog.Logger |
| FR-LOGCTX-001.3 | `From()` SHALL panic with descriptive message if no logger found in context | Test: call From() with clean context; verify panic occurs |
| FR-LOGCTX-001.4 | Panic message SHALL include "no logger found in context" for debugging | Panic recovery test; verify error string |
| FR-LOGCTX-001.5 | Function SHALL work with nested contexts (parent→child chains) | Test: create parent ctx with logger; derive child; child retrieves parent's logger |
| FR-LOGCTX-001.6 | Logger retrieval SHALL be O(1) with zero allocations | Benchmark WithLogger/From; verify < 10ns per call |
| FR-LOGCTX-001.7 | Logger type used SHALL be `*slog.Logger` from Go 1.21+ stdlib | Import from log/slog; no wrapper types |

---

## FR-RINGBUFFER-001: Generic Circular Buffer

| Requirement ID | Description | Verification |
|---|---|---|
| FR-RINGBUFFER-001.1 | Package `ringbuffer` SHALL export generic type `RingBuffer[T any]` | Type definition uses `[T any]` syntax |
| FR-RINGBUFFER-001.2 | Function `New[T](capacity int) *RingBuffer[T]` SHALL create buffer with fixed capacity | Constructor allocates slice of size capacity; capacity immutable |
| FR-RINGBUFFER-001.3 | Method `Push(item T)` SHALL add item to buffer | Item added to next position in circular array |
| FR-RINGBUFFER-001.4 | On overflow (buffer full), `Push()` SHALL overwrite oldest entry without error | Test: push 4 items to capacity 3; oldest is overwritten |
| FR-RINGBUFFER-001.5 | Method `GetAll() []T` SHALL return all items in FIFO order (oldest first) | Test: push 3 items; GetAll() returns them in order |
| FR-RINGBUFFER-001.6 | Method `Len() int` SHALL return current count of items in buffer | Count <= capacity; count == 0 after New() |
| FR-RINGBUFFER-001.7 | Method `Cap() int` SHALL return fixed capacity set at creation | Cap() == capacity arg to New() |
| FR-RINGBUFFER-001.8 | `GetAll()` SHALL return a copy (not reference to internal slice) | Modify returned slice; verify internal buffer unchanged |
| FR-RINGBUFFER-001.9 | Push/GetAll/Len/Cap SHALL be thread-safe with sync.RWMutex | Concurrent push + read; verify no data corruption or race conditions |
| FR-RINGBUFFER-001.10 | Overflow behavior (round-robin overwrite) SHALL use modulo arithmetic | Index wraps: (current + 1) % capacity |

---

## FR-WAITFOR-001: Polling with Exponential Backoff

| Requirement ID | Description | Verification |
|---|---|---|
| FR-WAITFOR-001.1 | Package `waitfor` SHALL export type `WaitTimeout` struct with fields: `Timeout`, `MinInterval`, `MaxInterval`, `InitialWait` | Struct definition with correct types: time.Duration for intervals, bool for InitialWait |
| FR-WAITFOR-001.2 | Function `WaitFor(ctx, timeout, condition) error` SHALL poll condition until true, timeout, or error | Implementation: loop calling condition; sleep between attempts |
| FR-WAITFOR-001.3 | Condition function signature: `func() (bool, error)` | Condition returns (satisfied bool, err error) |
| FR-WAITFOR-001.4 | If condition returns error, WaitFor SHALL propagate error immediately without retry | Test: condition returns error; WaitFor returns that error |
| FR-WAITFOR-001.5 | If condition returns false, WaitFor SHALL sleep and retry | Loop continues until condition true, timeout, or context deadline |
| FR-WAITFOR-001.6 | Backoff algorithm: sleep MinInterval * 2^(attempt-1), capped at MaxInterval | Attempt 1: sleep MinInterval * 2^0 = MinInterval; Attempt 2: sleep MinInterval * 2^1; capped at MaxInterval |
| FR-WAITFOR-001.7 | `InitialWait=false` SHALL check condition before first sleep | Attempt 0: condition check; if true, return immediately |
| FR-WAITFOR-001.8 | `InitialWait=true` SHALL sleep before first condition check | Attempt 0: sleep MinInterval; then check condition |
| FR-WAITFOR-001.9 | WaitFor SHALL return `ErrTimedOut` if timeout expires before condition true | Export ErrTimedOut sentinel; test verifies error equality |
| FR-WAITFOR-001.10 | WaitFor SHALL respect context cancellation (ctx.Done()) | Test: cancel context during poll; WaitFor returns context.Cause(ctx) |
| FR-WAITFOR-001.11 | Function `After(clock quartz.Clock, duration time.Duration) <-chan time.Time` SHALL return channel that fires after duration | If clock nil, use time.After(); if clock provided, use clock.After() |
| FR-WAITFOR-001.12 | After() SHALL support github.com/coder/quartz for testable (fake) clocks | Test: pass quartz.NewMock(); advance time; verify channel fires |

---

## FR-REGISTRY-001: Thread-Safe Generic Registry with Ref Counting

| Requirement ID | Description | Verification |
|---|---|---|
| FR-REGISTRY-001.1 | Package `registry` SHALL export generic type `Registry[K comparable, V any]` | Type definition uses [K comparable, V any] |
| FR-REGISTRY-001.2 | Function `New[K, V]() *Registry[K, V]` SHALL create empty registry | Constructor initializes map and mutex |
| FR-REGISTRY-001.3 | Method `Register(ownerID K, key K, value V)` SHALL store value under key claimed by owner | Entry created if not exists; ref count incremented if exists |
| FR-REGISTRY-001.4 | Multiple owners MAY register same key with same/different values | Test: owner-a and owner-b both register "svc"; both entries live as long as both own it |
| FR-REGISTRY-001.5 | Method `Unregister(ownerID K)` SHALL remove all entries owned by ownerID | All keys registered by ownerID are unregistered |
| FR-REGISTRY-001.6 | Unregister() SHALL decrement ref count; remove entry only if count reaches 0 | Test: 2 owners of same key; unregister 1; entry still exists; unregister 2; entry removed |
| FR-REGISTRY-001.7 | Method `Get(key K) (V, bool)` SHALL retrieve value by key | Returns (value, true) if exists; (zero, false) if not |
| FR-REGISTRY-001.8 | Method `Count(key K) int` SHALL return number of owners holding key | Count = len(owners) for key |
| FR-REGISTRY-001.9 | Method `List() map[K]V` SHALL return snapshot of all live entries | Returns copy of entries map; mutations don't affect registry |
| FR-REGISTRY-001.10 | Method `SetHook(hook Hook[K, V])` SHALL register observer for changes | Hook interface: OnRegister(ownerID K, key K, value V), OnUnregister(ownerID K) |
| FR-REGISTRY-001.11 | Hook.OnRegister() SHALL fire synchronously during Register() call | Observe every Register call; hook must complete before Register returns |
| FR-REGISTRY-001.12 | Hook.OnUnregister() SHALL fire synchronously during Unregister() call | Observe every Unregister call |
| FR-REGISTRY-001.13 | If no hook set, registry SHALL operate with zero observability overhead | No-op hook code path optimized away; same performance as unhooks version |
| FR-REGISTRY-001.14 | Registry operations SHALL be thread-safe | sync.RWMutex protects all reads and writes; concurrent Register/Get/Unregister safe |

---

## FR-CROSS-001: Zero Cross-Package Dependencies

| Requirement ID | Description | Verification |
|---|---|---|
| FR-CROSS-001.1 | logctx package SHALL have zero imports of ringbuffer, waitfor, registry | Grep imports in logctx/logctx.go; no github.com/KooshaPari/phenotype-go-kit imports |
| FR-CROSS-001.2 | ringbuffer package SHALL have zero imports of logctx, waitfor, registry | Same grep check |
| FR-CROSS-001.3 | waitfor package SHALL have zero imports of logctx, ringbuffer, registry | Same grep check |
| FR-CROSS-001.4 | registry package SHALL have zero imports of logctx, ringbuffer, waitfor | Same grep check |
| FR-CROSS-001.5 | No circular imports between any packages | `go build ./...` succeeds; `go mod graph` shows no cycles |

---

## FR-TEST-001: Test Coverage & Quality

| Requirement ID | Description | Verification |
|---|---|---|
| FR-TEST-001.1 | All packages SHALL have unit tests with >= 95% code coverage | `go test -cover ./...` reports coverage >= 95% for each package |
| FR-TEST-001.2 | logctx tests SHALL verify panic on From() with no logger | Test case: from(clean context) panics with "no logger found" |
| FR-TEST-001.3 | logctx tests SHALL verify logger retrieval in nested contexts | Parent context → child context → retrieve parent's logger |
| FR-TEST-001.4 | ringbuffer tests SHALL verify FIFO ordering on push/getall | Test: push [1,2,3]; getall returns [1,2,3] |
| FR-TEST-001.5 | ringbuffer tests SHALL verify overflow (oldest overwritten) | Test: capacity 3, push [1,2,3,4]; getall returns [2,3,4] |
| FR-TEST-001.6 | ringbuffer tests SHALL verify concurrent push/read (race detector) | `go test -race ./ringbuffer` passes |
| FR-TEST-001.7 | waitfor tests SHALL use quartz.NewMock() for deterministic timing | Test advances fake clock; verifies polling without real sleep |
| FR-TEST-001.8 | waitfor tests SHALL verify backoff algorithm: 50ms, 100ms, 200ms, capped at 2s | Test condition that fails N times; verify total sleep time matches backoff |
| FR-TEST-001.9 | waitfor tests SHALL verify ErrTimedOut on timeout | Test: timeout 100ms; condition never true; error == ErrTimedOut |
| FR-TEST-001.10 | waitfor tests SHALL verify immediate return if condition true at first check | InitialWait=false; condition true immediately; no sleep |
| FR-TEST-001.11 | registry tests SHALL verify ref counting (2 owners, unregister 1, entry still exists) | Test: Register(owner1, "svc"), Register(owner2, "svc"), Unregister(owner1), Get("svc") == ok |
| FR-TEST-001.12 | registry tests SHALL verify hook invocation on register/unregister | Test: SetHook(observer); Register/Unregister; verify OnRegister/OnUnregister called |
| FR-TEST-001.13 | registry tests SHALL verify concurrent register/unregister (race detector) | `go test -race ./registry` passes |
| FR-TEST-001.14 | All tests SHALL pass with `go test -race ./...` (race detector enabled) | Command succeeds; no race conditions detected |

---

## FR-BUILD-001: Build & Quality Gates

| Requirement ID | Description | Verification |
|---|---|---|
| FR-BUILD-001.1 | Go version requirement: >= 1.22 (generics, slog) | go.mod: `go 1.22` or later |
| FR-BUILD-001.2 | `go test ./...` SHALL pass all unit tests | All tests pass; no failures or panics |
| FR-BUILD-001.3 | `go test -race ./...` SHALL pass with no race conditions | Race detector finds zero data races |
| FR-BUILD-001.4 | `go vet ./...` SHALL find zero issues | Static analysis clean |
| FR-BUILD-001.5 | `gofumpt -l .` SHALL report zero formatting issues | Code formatted with gofumpt |
| FR-BUILD-001.6 | `golangci-lint run` SHALL report zero lint violations | Linter config in .golangci.yml; all rules green |
| FR-BUILD-001.7 | go.mod dependencies SHALL be minimal | No unnecessary vendoring; only coder/quartz for tests |
| FR-BUILD-001.8 | Module name: `github.com/KooshaPari/phenotype-go-kit` | go.mod contains correct module |

---

## FR-DOC-001: Documentation

| Requirement ID | Description | Verification |
|---|---|---|
| FR-DOC-001.1 | README.md SHALL include quick-start examples for each package | Examples: logctx WithLogger/From, ringbuffer Push/GetAll, waitfor WaitFor, registry Register/Get |
| FR-DOC-001.2 | All public types and functions SHALL have Godoc comments | `go doc ./...` displays comments for all exported symbols |
| FR-DOC-001.3 | Godoc comments SHALL include usage examples where applicable | E.g., `Example: rb.Push(1); rb.GetAll() // [1]` |
| FR-DOC-001.4 | README.md SHALL include sections: Overview, Packages, Install, Development | Structure matches standard Go project layout |
| FR-DOC-001.5 | Development section SHALL document: `go test -race ./...`, `go vet ./...`, `gofumpt`, `golangci-lint` | Clear instructions for contributors |

---

## Traceability Matrix

| FR ID | Package | Type | Priority |
|-------|---------|------|----------|
| FR-LOGCTX-001 | logctx | Core API | P1 |
| FR-RINGBUFFER-001 | ringbuffer | Core API | P1 |
| FR-WAITFOR-001 | waitfor | Core API | P1 |
| FR-REGISTRY-001 | registry | Core API | P1 |
| FR-CROSS-001 | All | Architecture | P1 |
| FR-TEST-001 | All | Quality | P1 |
| FR-BUILD-001 | All | Build | P1 |
| FR-DOC-001 | All | Documentation | P2 |

---

## Test Scenarios by Package

### logctx Test Scenarios
- T-LOGCTX-01: WithLogger stores logger in context
- T-LOGCTX-02: From retrieves logger from context
- T-LOGCTX-03: From panics if no logger (programmer error)
- T-LOGCTX-04: Nested contexts inherit parent logger

### ringbuffer Test Scenarios
- T-RINGBUFFER-01: Push/GetAll maintains FIFO order
- T-RINGBUFFER-02: Overflow overwrites oldest entry
- T-RINGBUFFER-03: Len/Cap return correct values
- T-RINGBUFFER-04: Concurrent push/read (race detector)

### waitfor Test Scenarios
- T-WAITFOR-01: Exponential backoff: 50ms → 100ms → 200ms → capped at 2s
- T-WAITFOR-02: InitialWait=false checks condition immediately
- T-WAITFOR-03: InitialWait=true sleeps before first check
- T-WAITFOR-04: ErrTimedOut on timeout
- T-WAITFOR-05: Condition error propagates immediately
- T-WAITFOR-06: Context cancellation interrupts polling
- T-WAITFOR-07: Fake clock (quartz) enables deterministic tests

### registry Test Scenarios
- T-REGISTRY-01: Single owner register/get
- T-REGISTRY-02: Multiple owners same key (ref counting)
- T-REGISTRY-03: Unregister decrements ref count
- T-REGISTRY-04: Last unregister removes entry
- T-REGISTRY-05: Hook.OnRegister fires on Register
- T-REGISTRY-06: Hook.OnUnregister fires on Unregister
- T-REGISTRY-07: List returns snapshot (safe iteration)
- T-REGISTRY-08: Concurrent register/unregister (race detector)
