# phenotype-go-kit

[![CI](https://github.com/KooshaPari/phenotype-go-kit/actions/workflows/ci.yml/badge.svg)](https://github.com/KooshaPari/phenotype-go-kit/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/KooshaPari/phenotype-go-kit.svg)](https://pkg.go.dev/github.com/KooshaPari/phenotype-go-kit)

Go infrastructure toolkit extracted from the Phenotype ecosystem. Small, focused packages for logging, data structures, polling, and registries.

## Packages

| Package | Description |
|---------|-------------|
| [`logctx`](logctx/) | Context-scoped `slog.Logger` injection and retrieval |
| [`ringbuffer`](ringbuffer/) | Generic fixed-capacity circular buffer |
| [`waitfor`](waitfor/) | Polling with exponential backoff, configurable timeout, and testable clocks |
| [`registry`](registry/) | Generic thread-safe key-value registry with owner tracking, ref counting, and change hooks |

## Install

```bash
go get github.com/KooshaPari/phenotype-go-kit
```

Requires Go 1.22+.

## Usage

### logctx

Attach a `*slog.Logger` to a `context.Context` and retrieve it anywhere downstream.
`From` panics if no logger has been injected — this is intentional: a missing logger is a
programmer error and should fail loudly.

```go
import (
    "context"
    "log/slog"

    "github.com/KooshaPari/phenotype-go-kit/logctx"
)

ctx := logctx.WithLogger(context.Background(), slog.Default())
logger := logctx.From(ctx)
logger.Info("hello from context logger")
```

### ringbuffer

A generic, fixed-capacity circular buffer. When full, `Push` overwrites the oldest entry.

```go
import "github.com/KooshaPari/phenotype-go-kit/ringbuffer"

rb := ringbuffer.New[int](3)
rb.Push(1)
rb.Push(2)
rb.Push(3)
rb.Push(4)           // overwrites 1
items := rb.GetAll() // [2, 3, 4] — oldest first
fmt.Println(rb.Len()) // 3
fmt.Println(rb.Cap()) // 3
```

### waitfor

Poll a condition with exponential backoff until it returns `true`, an error occurs, or the
timeout expires. Integrates with [`github.com/coder/quartz`](https://github.com/coder/quartz)
for deterministic testing.

```go
import (
    "context"
    "time"

    "github.com/KooshaPari/phenotype-go-kit/waitfor"
)

err := waitfor.WaitFor(ctx, waitfor.WaitTimeout{
    Timeout:     10 * time.Second,
    MinInterval: 50 * time.Millisecond,
    MaxInterval: 500 * time.Millisecond,
    InitialWait: false, // check condition immediately before first sleep
}, func() (bool, error) {
    return isReady(), nil
})
if err != nil {
    // err is waitfor.ErrTimedOut or a condition error
}
```

`After` is a helper that returns a channel that fires after a duration using any
`quartz.Clock` (pass `nil` for the real clock):

```go
<-waitfor.After(nil, 5*time.Second)
```

### registry

A generic, thread-safe key-value store with owner-scoped lifecycle management. Multiple owners
may hold the same key; the entry is removed only when the last owner unregisters. An optional
`Hook` interface observes all changes.

```go
import "github.com/KooshaPari/phenotype-go-kit/registry"

type ServiceInfo struct{ Port int }

reg := registry.New[string, ServiceInfo]()

// Two owners register the same key.
reg.Register("owner-a", "api-svc", ServiceInfo{Port: 8080})
reg.Register("owner-b", "api-svc", ServiceInfo{Port: 8080})

svc, ok := reg.Get("api-svc") // (ServiceInfo{8080}, true)
count := reg.Count("api-svc") // 2

// Removing one owner decrements the ref count.
reg.Unregister("owner-a")
count = reg.Count("api-svc") // 1

// Removing the last owner deletes the entry.
reg.Unregister("owner-b")
_, ok = reg.Get("api-svc") // (zero, false)

// Snapshot of all live entries.
all := reg.List() // map[string]ServiceInfo
```

Implement the `Hook` interface to observe changes:

```go
type myHook struct{}

func (h *myHook) OnRegister(ownerID string, key string, value ServiceInfo) {
    fmt.Printf("registered %s by %s\n", key, ownerID)
}
func (h *myHook) OnUnregister(ownerID string) {
    fmt.Printf("unregistered owner %s\n", ownerID)
}

reg.SetHook(&myHook{})
```

## Development

```bash
go test -race ./...   # Run all tests with race detector
go vet ./...          # Static analysis
gofumpt -l .          # Format check
golangci-lint run     # Full lint suite
```

## License

MIT
