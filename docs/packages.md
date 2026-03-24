# Package Reference

## logctx

**Import:** `github.com/KooshaPari/phenotype-go-kit/logctx`

**Purpose:** Attaches and retrieves a `*slog.Logger` from a `context.Context`.

### API

| Symbol | Signature | Description |
|--------|-----------|-------------|
| `WithLogger` | `func WithLogger(ctx context.Context, logger *slog.Logger) context.Context` | Returns a child context carrying logger. |
| `From` | `func From(ctx context.Context) *slog.Logger` | Retrieves the logger or panics if none was set. |

**Design note:** `From` panics rather than returning a default logger so that missing logger
injection is caught at development time, not silently degraded in production.

---

## ringbuffer

**Import:** `github.com/KooshaPari/phenotype-go-kit/ringbuffer`

**Purpose:** Generic fixed-capacity circular buffer. `Push` overwrites the oldest entry when
the buffer is full.

### API

| Symbol | Signature | Description |
|--------|-----------|-------------|
| `New[T]` | `func New[T any](size int) *RingBuffer[T]` | Creates a ring buffer with the given capacity. |
| `Push` | `func (r *RingBuffer[T]) Push(item T)` | Appends item; evicts oldest if full. |
| `GetAll` | `func (r *RingBuffer[T]) GetAll() []T` | Returns all items, oldest first. |
| `Len` | `func (r *RingBuffer[T]) Len() int` | Current number of stored items. |
| `Cap` | `func (r *RingBuffer[T]) Cap() int` | Maximum capacity. |

**Not thread-safe.** Callers must synchronize if used from multiple goroutines.

---

## waitfor

**Import:** `github.com/KooshaPari/phenotype-go-kit/waitfor`

**Purpose:** Poll a condition with exponential backoff until it passes, returns an error, or
times out. Uses `github.com/coder/quartz` for a mockable clock — making timeout/interval
behaviour fully testable without `time.Sleep`.

### API

| Symbol | Description |
|--------|-------------|
| `WaitTimeout` | Configuration struct for `WaitFor`. |
| `WaitTimeout.Timeout` | Overall deadline (default 10 s). |
| `WaitTimeout.MinInterval` | Starting backoff interval (default 10 ms). |
| `WaitTimeout.MaxInterval` | Backoff ceiling (default 500 ms). |
| `WaitTimeout.InitialWait` | When true, sleep before the first condition check. |
| `WaitTimeout.Clock` | Optional `quartz.Clock`; nil uses a real clock. |
| `ErrTimedOut` | Sentinel error returned when the deadline expires. |
| `WaitFor` | `func WaitFor(ctx context.Context, timeout WaitTimeout, condition func() (bool, error)) error` |
| `After` | `func After(clk quartz.Clock, d time.Duration) <-chan time.Time` |

### Backoff behaviour

`WaitFor` doubles the sleep interval on each failed check, clamped between `MinInterval` and
`MaxInterval`. Context cancellation is respected at every sleep boundary.

---

## registry

**Import:** `github.com/KooshaPari/phenotype-go-kit/registry`

**Purpose:** Thread-safe key-value store with owner-scoped lifecycle. Each `Register` call
increments a reference count for the key; `Unregister` decrements all counts owned by the
caller and removes keys that reach zero.

### API

| Symbol | Description |
|--------|-------------|
| `New[K, V]` | `func New[K comparable, V any]() *Registry[K, V]` |
| `Register` | Adds or increments a key under ownerID. |
| `Unregister` | Decrements counts for all keys owned by ownerID; removes entries at zero. |
| `Get` | Returns the value for a key and whether it exists. |
| `List` | Returns a snapshot map of all live entries. |
| `Count` | Returns the current reference count for a key. |
| `SetHook` | Attaches a `Hook[K, V]` observer. |
| `Hook[K, V]` | Interface: `OnRegister(ownerID string, key K, value V)` and `OnUnregister(ownerID string)`. |

**Thread-safe.** Uses `sync.RWMutex` internally.
