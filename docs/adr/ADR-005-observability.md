# ADR-005: Observability Standards

## Status
**Accepted** | 2024-01-15

## Context

Distributed systems require comprehensive observability to debug issues and understand behavior.

## Decision

### Three Pillars

```
┌─────────────┐   ┌─────────────┐   ┌─────────────┐
│   LOGS     │   │   METRICS   │   │   TRACES    │
│             │   │             │   │             │
│ "What       │   │ "How much   │   │ "How does  │
│  happened?" │   │  and when?" │   │  it flow?"  │
└─────────────┘   └─────────────┘   └─────────────┘
```

### Structured Logging (logctx)
- JSON format for machine parsing
- Correlation IDs (trace_id, request_id)
- Levels: DEBUG, INFO, WARN, ERROR
- Log context propagation via `context.Context`

### Metrics (metrics)
- RED metrics: Rate, Errors, Duration
- USE metrics: Utilization, Saturation, Errors
- Cardinality-aware labeling
- Prometheus-compatible format

### Distributed Tracing (tracing)
- OpenTelemetry standard
- Automatic span propagation
- Sampling strategy documented
- Service mesh integration

## Key Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `http_requests_total` | Counter | Total HTTP requests by method, path, status |
| `http_request_duration_seconds` | Histogram | Request latency distribution |
| `db_queries_total` | Counter | Database queries by query type |
| `db_query_duration_seconds` | Histogram | Query latency |
| `cache_hits_total` | Counter | Cache hit/miss ratio |
| `errors_total` | Counter | Errors by type and location |

## Implementation

```go
// Middleware for HTTP observability
func ObservedHandler(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        
        // Extract or generate trace ID
        span := trace.SpanFromContext(ctx)
        spanCtx := span.SpanContext()
        
        // Record start time
        start := time.Now()
        
        // Wrap response writer to capture status
        wrapped := &statusWriter{ResponseWriter: w, status: 200}
        
        handler.ServeHTTP(wrapped, r)
        
        // Record metrics
        metrics.RecordHTTP(ctx, 
            r.Method, r.URL.Path,
            wrapped.status,
            time.Since(start),
        )
        
        // Add trace context to logs
        logctx.FromContext(ctx).Info("request completed",
            logctx.Int("status", wrapped.status),
            logctx.Duration("duration", time.Since(start)),
        )
    })
}
```

## Consequences

### Positive
- Full visibility into system behavior
- Correlated data across pillars
- Production debugging capability
- Performance optimization insights

### Negative
- Added instrumentation code
- Storage costs for traces/logs
- Sampling strategy needed for high-volume

## References
- [OpenTelemetry](https://opentelemetry.io/)
- [RED Method](https://www.weave.works/blog/red-prometheus)
- [USE Method](https://www.brendangregg.com/usemethod.html)
