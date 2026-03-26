# phenotype-go-kit — Product Requirements Document

## Executive Summary

**phenotype-go-kit** is a focused, minimal Go infrastructure toolkit extracted from the Phenotype ecosystem. It provides four independent, composable packages for common infrastructure concerns: context-scoped logging, fixed-capacity circular buffers, polling with exponential backoff, and thread-safe generic registries.

### Product Vision

Enable Go developers to build robust, observable, and composable microservices and agents by providing small, focused, well-tested packages that follow SOLID principles and composable architecture patterns. Each package has zero cross-package dependencies, allowing developers to use only what they need without bloat.

---

## User Personas & Use Cases

### Persona 1: Distributed Service Developer
- **Goal**: Build multi-service systems where context-aware logging, service discovery, and resilient polling are critical
- **Use Case**: Inject slog.Logger into context; use registry for service discovery; poll for health checks with backoff
- **Key Needs**: Type-safe, testable, minimal dependencies

### Persona 2: Systems Engineer
- **Goal**: Instrument observability with context-scoped logging across call chains
- **Use Case**: Retrieve logger from context in any downstream function; logs include request context and traces
- **Key Needs**: Zero allocation overhead, context-first design, slog integration

### Persona 3: Agent & AI Framework Developer
- **Goal**: Build robust polling, state management, and service registration patterns for long-running agents
- **Use Case**: Poll agent state with exponential backoff; use registry to track agent instances; buffer recent events
- **Key Needs**: Configurable timeouts, testable with fake clocks, ref-counted lifecycle management

---

## Product Architecture

### Four Independent Packages

```
github.com/KooshaPari/phenotype-go-kit/

├── logctx/        # Context-scoped slog.Logger injection
├── ringbuffer/    # Generic fixed-capacity circular buffer
├── waitfor/       # Polling with exponential backoff
└── registry/      # Thread-safe generic key-value registry with ref counting
```

**Design Principle**: Each package is independent with zero cross-dependencies. Developers can import only the packages they need.

---

## Package Specifications

### PKG-1: logctx — Context-Scoped Logging

**Purpose**: Inject and retrieve `*slog.Logger` from `context.Context` anywhere in the call chain.

#### Core Abstractions

```
Package: github.com/KooshaPari/phenotype-go-kit/logctx

Functions:
- WithLogger(ctx context.Context, logger *slog.Logger) context.Context
  └─ Returns new context with logger attached
  
- From(ctx context.Context) *slog.Logger
  └─ Retrieves logger; panics if missing (intentional: programmer error)
```

#### Design Rationale

- **Panic on Missing Logger**: Missing logger is a programmer error, not a runtime concern. Fail loudly.
- **No Wrapper**: Return `*slog.Logger` directly; leverage Go 1.21+ stdlib slog.
- **Zero Overhead**: Context values are just pointers; no serialization or overhead.
- **Thread-Safe**: slog.Logger is thread-safe; context.WithValue is thread-safe.

#### User Stories

**Story PKG-1.1: Inject Logger at Request Entry**
```
func handleRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) {
  logger := slog.Default().With("request_id", requestID)
  ctx = logctx.WithLogger(ctx, logger)
  
  // Logger is now available to all downstream functions
  processRequest(ctx)
}
```

**Story PKG-1.2: Retrieve Logger Anywhere Downstream**
```
func processRequest(ctx context.Context) {
  logger := logctx.From(ctx)
  logger.Info("processing request")
}
```

**Story PKG-1.3: Testing with Slog Handler Inspection**
```
func TestLogging(t *testing.T) {
  var buf bytes.Buffer
  handler := slog.NewTextHandler(&buf, nil)
  logger := slog.New(handler)
  
  ctx := logctx.WithLogger(context.Background(), logger)
  processRequest(ctx)
  
  assert.Contains(t, buf.String(), "processing request")
}
```

---

### PKG-2: ringbuffer — Generic Fixed-Capacity Circular Buffer

**Purpose**: Store up to N items (most recent) and retrieve them in FIFO order, discarding oldest entries on overflow.

#### Core Abstractions

```
Package: github.com/KooshaPari/phenotype-go-kit/ringbuffer

Type: RingBuffer[T any]

Methods:
- New[T](capacity int) *RingBuffer[T]
  └─ Create circular buffer with fixed capacity
  
- Push(item T)
  └─ Add item; overwrite oldest if full
  
- GetAll() []T
  └─ Return all items in FIFO order (oldest first)
  
- Len() int
  └─ Current count
  
- Cap() int
  └─ Maximum capacity
```

#### Implementation Details

- **Generic**: Supports any type T
- **Fixed Capacity**: Capacity set at creation; immutable
- **Circular**: Uses modulo arithmetic to wrap around
- **Thread-Safe**: Protected by sync.RWMutex for concurrent reads
- **GetAll() Ordering**: Returns items in FIFO order (oldest first)

#### User Stories

**Story PKG-2.1: Recent Event Buffer**
```
rb := ringbuffer.New[Event](100)
rb.Push(event1)
rb.Push(event2)
rb.Push(event3)

// Later, retrieve recent events
events := rb.GetAll()  // [event1, event2, event3]
```

**Story PKG-2.2: Circular Behavior on Overflow**
```
rb := ringbuffer.New[int](3)
rb.Push(1)
rb.Push(2)
rb.Push(3)
rb.Push(4)  // Overwrites 1

items := rb.GetAll()  // [2, 3, 4]
```

**Story PKG-2.3: Agent State History**
```
type AgentState struct {
  Timestamp time.Time
  Status    string
  Token     int
}

history := ringbuffer.New[AgentState](50)
for event := range eventStream {
  history.Push(AgentState{
    Timestamp: event.At,
    Status:    event.State,
    Token:     event.Tokens,
  })
}

// Inspect recent state for diagnostics
recentStates := history.GetAll()
```

---

### PKG-3: waitfor — Polling with Exponential Backoff

**Purpose**: Poll a condition with exponential backoff, timeout, and testable clock support.

#### Core Abstractions

```
Package: github.com/KooshaPari/phenotype-go-kit/waitfor

Type: WaitTimeout struct {
  Timeout     time.Duration
  MinInterval time.Duration  // Base interval (e.g., 50ms)
  MaxInterval time.Duration  // Cap on exponential growth
  InitialWait bool           // Check immediately or wait first
}

Functions:
- WaitFor(ctx context.Context, timeout WaitTimeout, 
          condition func() (bool, error)) error
  └─ Poll until condition true, timeout, or error
  └─ Returns: ErrTimedOut or condition error
  
- After(clock quartz.Clock, duration time.Duration) <-chan time.Time
  └─ Sleep helper; uses real clock if nil
```

#### Backoff Algorithm

```
Attempt 1: wait InitialWait ? MinInterval : 0
Attempt 2: wait MinInterval * 2^1 (capped at MaxInterval)
Attempt 3: wait MinInterval * 2^2 (capped at MaxInterval)
...
```

#### Design Rationale

- **quartz.Clock Integration**: Testable; allows fake time advancement in tests
- **Configurable Intervals**: MinInterval, MaxInterval, InitialWait for flexibility
- **Error Propagation**: If condition returns error, propagate immediately (don't retry)
- **Timeout Semantics**: Context deadline + explicit timeout for fine-grained control

#### User Stories

**Story PKG-3.1: Poll for Resource Readiness**
```
err := waitfor.WaitFor(ctx, waitfor.WaitTimeout{
  Timeout:     10 * time.Second,
  MinInterval: 100 * time.Millisecond,
  MaxInterval: 2 * time.Second,
  InitialWait: false,  // Check immediately
}, func() (bool, error) {
  resp, err := http.Get("http://service:8080/health")
  if err != nil {
    return false, nil  // Network error; retry
  }
  return resp.StatusCode == 200, nil
})

if err == waitfor.ErrTimedOut {
  log.Fatal("service never became healthy")
}
```

**Story PKG-3.2: Testable Polling with Fake Clock**
```
func TestPolling(t *testing.T) {
  clock := quartz.NewMock()  // Fake clock from github.com/coder/quartz
  
  err := waitfor.WaitFor(ctx, waitfor.WaitTimeout{
    Timeout:     1 * time.Second,
    MinInterval: 100 * time.Millisecond,
    InitialWait: true,
  }, func() (bool, error) {
    return readyFlag, nil
  })
  
  // Advance fake clock
  clock.Advance(500 * time.Millisecond)
  
  // Test completes instantly (fake time)
  assert.NoError(t, err)
}
```

**Story PKG-3.3: Agent State Polling**
```
// Agent framework: poll for task completion
err := waitfor.WaitFor(ctx, waitfor.WaitTimeout{
  Timeout:     5 * time.Minute,
  MinInterval: 50 * time.Millisecond,
  MaxInterval: 5 * time.Second,
  InitialWait: false,
}, func() (bool, error) {
  state, err := agent.GetState()
  if err != nil {
    return false, err  // Propagate error; don't retry
  }
  return state.Complete, nil
})
```

---

### PKG-4: registry — Thread-Safe Generic Registry with Ref Counting

**Purpose**: Manage lifecycle of named entities with owner-scoped registration. Multiple owners can hold the same key; entry removed only when last owner unregisters.

#### Core Abstractions

```
Package: github.com/KooshaPari/phenotype-go-kit/registry

Type: Registry[K comparable, V any]

Methods:
- New[K, V]() *Registry[K, V]
  └─ Create new registry
  
- Register(ownerID K, key K, value V)
  └─ Register value under key by owner
  └─ Increments ref count if key already exists
  
- Unregister(ownerID K)
  └─ Remove all entries owned by ownerID
  └─ Decrements ref count; removes entry if count == 0
  
- Get(key K) (V, bool)
  └─ Retrieve value; ok=false if not found
  
- Count(key K) int
  └─ Ref count for key
  
- List() map[K]V
  └─ Snapshot of all live entries
  
- SetHook(hook Hook[K, V])
  └─ Observe registration/unregistration events
```

#### Hook Interface

```
type Hook[K comparable, V any] interface {
  OnRegister(ownerID K, key K, value V)
  OnUnregister(ownerID K)
}
```

#### Design Rationale

- **Owner-Scoped Lifecycle**: Multiple owners can register same key; entry persists as long as any owner claims it
- **Ref Counting**: Automatic cleanup when last owner unregisters
- **Hook Support**: Observe all changes; useful for logging, metrics, side effects
- **Snapshot Safety**: List() returns copy to avoid mutation issues
- **Thread-Safe**: sync.RWMutex for concurrent access

#### User Stories

**Story PKG-4.1: Multi-Owner Service Registry**
```
reg := registry.New[string, ServiceInfo]()

// Two agents register the same service (with different ports)
reg.Register("agent-a", "api-svc", ServiceInfo{Port: 8080})
reg.Register("agent-b", "api-svc", ServiceInfo{Port: 8081})

// Service is active as long as either owner is active
svc, ok := reg.Get("api-svc")  // (ServiceInfo{...}, true)
count := reg.Count("api-svc")  // 2

// Agent-A shuts down
reg.Unregister("agent-a")
count = reg.Count("api-svc")  // 1

// Service still active (agent-b still owns it)
svc, ok = reg.Get("api-svc")  // (ServiceInfo{...}, true)

// Agent-B shuts down
reg.Unregister("agent-b")
count = reg.Count("api-svc")  // 0
svc, ok = reg.Get("api-svc")  // (zero, false) — entry removed
```

**Story PKG-4.2: Observing Registry Changes**
```
type MyHook struct{}

func (h *MyHook) OnRegister(ownerID string, key string, value ServiceInfo) {
  log.Printf("service registered: %s by %s", key, ownerID)
}

func (h *MyHook) OnUnregister(ownerID string) {
  log.Printf("owner unregistered: %s", ownerID)
}

reg := registry.New[string, ServiceInfo]()
reg.SetHook(&MyHook{})
```

**Story PKG-4.3: Agent Lifecycle Management**
```
// Agents register their services on startup
agentID := "agent-42"
services := map[string]ServiceInfo{
  "task-queue": ServiceInfo{Host: "localhost", Port: 5672},
  "cache":      ServiceInfo{Host: "localhost", Port: 6379},
}

for name, info := range services {
  reg.Register(agentID, name, info)
}

// On graceful shutdown, unregister all
defer reg.Unregister(agentID)  // Removes all services owned by agent
```

**Story PKG-4.4: Snapshot and Inspection**
```
// Inspect all registered services
allServices := reg.List()  // map[string]ServiceInfo
for key, svc := range allServices {
  fmt.Printf("Service %s: %s:%d\n", key, svc.Host, svc.Port)
}

// Count owners of a service
refs := reg.Count("api-svc")
fmt.Printf("Service api-svc has %d owner(s)\n", refs)
```

---

## Technical Requirements

### FR-PKG1-001: logctx Context Injection
- **Shall**: Provide `WithLogger(ctx, logger)` and `From(ctx)` functions
- **Shall**: Panic on `From()` if no logger in context (intentional failure)
- **Shall**: Work with `log/slog` from Go 1.21+ stdlib
- **Shall**: Support nested context chains (parent→child contexts)

### FR-PKG1-002: logctx Type Safety
- **Shall**: Return `*slog.Logger` directly (no wrapper type)
- **Shall**: Use context.ContextKey private type to prevent collisions
- **Shall**: Thread-safe (context.WithValue is thread-safe)

### FR-PKG2-001: ringbuffer Generics
- **Shall**: Support any type T via `[T any]` generics (Go 1.18+)
- **Shall**: Fixed capacity set at creation; immutable
- **Shall**: Circular buffer with modulo arithmetic
- **Shall**: GetAll() returns copy to avoid mutation

### FR-PKG2-002: ringbuffer Thread Safety
- **Shall**: Protect with sync.RWMutex for concurrent access
- **Shall**: Push() and GetAll() safe for concurrent reads/writes
- **Shall**: No allocation beyond what's necessary

### FR-PKG3-001: waitfor Polling
- **Shall**: Support exponential backoff: MinInterval * 2^attempt (capped at MaxInterval)
- **Shall**: Accept timeout, min/max intervals, InitialWait flag
- **Shall**: Return ErrTimedOut or condition error
- **Shall**: Handle context cancellation (context deadline)

### FR-PKG3-002: waitfor Testability
- **Shall**: Integrate with github.com/coder/quartz for testable clocks
- **Shall**: After() function accepts quartz.Clock (nil → real clock)
- **Shall**: Tests can advance fake clock and verify polling behavior

### FR-PKG4-001: registry Generics
- **Shall**: Support generic keys K comparable, values V any
- **Shall**: Register(ownerID, key, value), Unregister(ownerID), Get(key), Count(key), List()
- **Shall**: Ref counting: entry removed only when last owner unregisters
- **Shall**: Thread-safe with sync.RWMutex

### FR-PKG4-002: registry Hooks
- **Shall**: Support Hook[K, V] interface with OnRegister() and OnUnregister() methods
- **Shall**: Optional; if no hook, no observability overhead
- **Shall**: Hooks called synchronously during Register/Unregister

### FR-PKG4-003: registry Snapshots
- **Shall**: List() returns map copy (no pointer sharing)
- **Shall**: Safe for iteration without holding locks

### FR-CROSS-001: Zero Cross-Package Dependencies
- **Shall**: logctx independent of ringbuffer, waitfor, registry
- **Shall**: ringbuffer independent of logctx, waitfor, registry
- **Shall**: waitfor independent of other packages
- **Shall**: registry independent of other packages
- **Shall**: No circular imports

### FR-CROSS-002: Testing
- **Shall**: All packages have unit tests with >90% coverage
- **Shall**: logctx tests verify panic on missing logger
- **Shall**: ringbuffer tests verify overflow and FIFO ordering
- **Shall**: waitfor tests use fake clocks (quartz) for deterministic timing
- **Shall**: registry tests verify ref counting and hook invocation

### FR-CROSS-003: Documentation
- **Shall**: README.md with quick-start examples for each package
- **Shall**: Godoc comments on all public types and functions
- **Shall**: Example code files (example_test.go) demonstrating usage

### FR-CROSS-004: Build & Quality
- **Shall**: Go 1.22+ required (generics, slog)
- **Shall**: `go test -race ./...` passes (race detector enabled)
- **Shall**: `go vet ./...` clean (static analysis)
- **Shall**: `gofumpt -l .` formatting check
- **Shall**: `golangci-lint run` full lint suite

---

## Success Metrics

| Metric | Target |
|--------|--------|
| Test coverage | >= 95% |
| Package independence | Zero cross-package imports |
| Documentation | 100% Godoc coverage |
| Performance | < 1µs overhead per logctx lookup |
| Registry ref count accuracy | 100% (unit test) |

---

## Out of Scope (v1)

- Web UI for registry inspection
- Metrics/observability integration (future: prometheus)
- Distributed registry (future: consensus-based)
- Custom serialization plugins
- Plugin architecture for hooks

---

## Dependencies

### Required (stdlib only)
- Go 1.22+
- `context`
- `log/slog`
- `sync`
- `time`

### Optional (tests only)
- `github.com/coder/quartz` (v0.1.2+) — fake clock for waitfor tests

### Zero Runtime Dependencies
- All packages use only Go stdlib

---

## Version 1 Release Checklist

- [x] logctx: WithLogger, From, panic semantics
- [x] ringbuffer: generic circular buffer with thread safety
- [x] waitfor: exponential backoff polling with quartz integration
- [x] registry: ref-counted generic registry with hooks
- [x] Unit tests with >95% coverage for all packages
- [x] README.md with examples
- [x] Godoc coverage for all public API
- [x] go.mod / go.sum setup
- [x] CI: tests, lint, race detector
- [ ] Comprehensive examples (e.g., multi-package integration)
- [ ] Performance benchmarks
- [ ] API stability guarantee (semantic versioning)
