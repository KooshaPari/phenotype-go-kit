# phenotype-go-kit

Go infrastructure toolkit extracted from the Phenotype ecosystem. Small, focused packages for logging, data structures, polling, and registries.

## Packages

| Package | Description |
|---------|-------------|
| [`logctx`](logctx/) | Context-scoped `slog.Logger` injection and retrieval |
| [`ringbuffer`](ringbuffer/) | Generic fixed-capacity ring buffer with iteration |
| [`waitfor`](waitfor/) | Polling with configurable intervals and timeout using testable clocks |
| [`registry`](registry/) | Generic thread-safe registry with ref counting, quota tracking, and hooks |

## Install

```bash
go get github.com/KooshaPari/phenotype-go-kit
```

## Usage

### logctx

```go
import "github.com/KooshaPari/phenotype-go-kit/logctx"

ctx := logctx.WithLogger(ctx, slog.Default())
logger := logctx.From(ctx) // retrieves the logger
```

### ringbuffer

```go
import "github.com/KooshaPari/phenotype-go-kit/ringbuffer"

rb := ringbuffer.New[int](3)
rb.Push(1)
rb.Push(2)
rb.Push(3)
rb.Push(4) // overwrites 1
rb.Do(func(v int) { fmt.Println(v) }) // prints 2, 3, 4
```

### waitfor

```go
import "github.com/KooshaPari/phenotype-go-kit/waitfor"

err := waitfor.WaitFor(ctx, waitfor.WaitTimeout{
    Timeout:  10 * time.Second,
    Interval: 500 * time.Millisecond,
}, func() (bool, error) {
    return checkReady(), nil
})
```

### registry

```go
import "github.com/KooshaPari/phenotype-go-kit/registry"

reg := registry.New[string, ServiceInfo]()
reg.Register("owner-1", "api-svc", ServiceInfo{Port: 8080})
svc, ok := reg.Get("api-svc")
count := reg.Count("api-svc")
reg.Unregister("owner-1")
```

## Development

```bash
go test ./...       # Run all tests
go vet ./...        # Lint
gofumpt -l .        # Format check
```

## License

MIT
