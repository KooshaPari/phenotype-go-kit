# ADR-003: Context-Propagated Logging

## Status
**Accepted** | 2024-01-15

## Context

Distributed services need correlated logging across multiple components (request ID, user ID, trace ID).

**Problem:** Without context propagation, logs from different components of one request are scattered.

## Decision

Use `context.Context` to propagate request-scoped values including:

- `trace_id` - distributed trace identifier
- `request_id` - unique request identifier
- `user_id` - authenticated user (if applicable)
- `org_id` - organization context

### Interface

```go
package logctx

type Logger interface {
    Debug(msg string, fields ...Field)
    Info(msg string, fields ...Field)
    Warn(msg string, fields ...Field)
    Error(msg string, fields ...Field)
    With(fields ...Field) Logger
}

type Field struct {
    Key   string
    Value any
}

func FromContext(ctx context.Context) Logger
func WithContext(ctx context.Context, logger Logger) context.Context
```

### Usage

```go
// Entry point - create logger with context
logger := logctx.FromContext(ctx).With(
    logctx.String("trace_id", traceID),
    logctx.String("user_id", userID),
)

// Propagate to downstream
ctx = logctx.WithContext(ctx, logger)
downstream.Call(ctx)

// All logs automatically include trace_id, user_id
logger.Info("request completed", logctx.Int("status", 200))
```

## Consequences

### Positive
- Correlated logs across services
- Easy to filter by request/trace
- Standardized fields
- Zero-cost when not used

### Negative
- Must pass context everywhere
- Need middleware to extract context values

## Implementation

See `logctx/` package for concrete implementation.

## References
- [OpenTelemetry: Context Propagation](https://opentelemetry.io/docs/reference/specification/context/)
- [Structured Logging](https://www.honeycomb.io/blog/observability-a-primer/)
