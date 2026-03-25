# Comparison Matrix

## Feature Comparison

This document compares **phenotype-go-kit** with similar tools in the Go infrastructure toolkit space.

| Repository | Purpose | Key Features | Language/Framework | Maturity | Comparison |
|------------|---------|--------------|-------------------|----------|------------|
| **phenotype-go-kit (this repo)** | Go infrastructure toolkit | logctx, ringbuffer, waitfor, registry | Go | Stable | Phenotype ecosystem |
| [zerolog](https://github.com/rs/zerolog) | Logging | Zero allocation, JSON, Structured | Go | Stable | Production logging |
| [zap](https://github.com/uber-go/zap) | Logging | High performance, Structured | Go | Stable | Uber's logging |
| [slog](https://github.com/golang/example) | Standard logging | Go 1.21+, Structured | Go | Stable | stdlib solution |
| [goval](https://github.com/mitchellh/go-watcher) | Polling/waiting | Retry logic | Go | Stable | Retry polling |
| [go-redis](https://github.com/redis/go-redis) | Redis client | Redis operations | Go | Stable | Redis integration |

## Detailed Feature Comparison

### Packages

| Package | phenotype-go-kit | zerolog | zap | slog |
|---------|-----------------|---------|-----|------|
| logctx | ✅ | ❌ | ❌ | ✅ (via context) |
| ringbuffer | ✅ | ❌ | ❌ | ❌ |
| waitfor | ✅ | ❌ | ❌ | ❌ |
| registry | ✅ | ❌ | ❌ | ❌ |

### logctx Features

| Feature | phenotype-go-kit | slog (stdlib) | zerolog | zap |
|---------|-----------------|---------------|---------|-----|
| Context-scoped | ✅ | ✅ | ❌ | ❌ |
| slog integration | ✅ | ✅ | ❌ | ❌ |
| Panic on missing | ✅ | ❌ | N/A | N/A |
| Type-safe | ✅ | ✅ | N/A | N/A |

### waitfor Features

| Feature | phenotype-go-kit | goval | golang.org/x/sync |
|---------|-----------------|-------|--------------------|
| Exponential backoff | ✅ | ✅ | ✅ (ErrWait) |
| Configurable timeout | ✅ | ✅ | ❌ |
| Testable clocks | ✅ | ❌ | ❌ |
| Initial wait option | ✅ | ❌ | ❌ |

### registry Features

| Feature | phenotype-go-kit | sync.Map | redis |
|---------|-----------------|----------|-------|
| Generic types | ✅ | ❌ | N/A |
| Owner tracking | ✅ | ❌ | ❌ |
| Ref counting | ✅ | ❌ | ❌ |
| Change hooks | ✅ | ❌ | ❌ |

## Unique Value Proposition

phenotype-go-kit provides:

1. **Context-Scoped Logging**: `logctx` for request-bound logger injection
2. **Circular Buffer**: Generic ring buffer for bounded queues
3. **Polling with Backoff**: `waitfor` with exponential backoff and testable clocks
4. **Owner-Based Registry**: Ref-counted registry with lifecycle management

## Packages

| Package | Description |
|---------|-------------|
| `logctx` | Context-scoped slog.Logger injection and retrieval |
| `ringbuffer` | Generic fixed-capacity circular buffer |
| `waitfor` | Polling with exponential backoff and configurable timeout |
| `registry` | Generic thread-safe key-value registry with owner tracking |

## References

- zerolog: [rs/zerolog](https://github.com/rs/zerolog)
- zap: [uber-go/zap](https://github.com/uber-go/zap)
- slog: [golang/example](https://github.com/golang/example)
