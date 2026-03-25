# ADR-009: Hexagonal Architecture Implementation

## Status
**Proposed** | 2024-01-15

## Context

We need to implement hexagonal (ports & adapters) architecture in phenotype-go-kit to achieve:
- Swappable implementations (Redis ↔ Memcached, Postgres ↔ MongoDB)
- Easy unit testing with mocks
- Clear separation between domain logic and infrastructure
- Support for multiple entry points (CLI, API, gRPC)

## Decision

### Directory Structure

```
phenotype-go-kit/
├── domain/                    # Pure business logic (KISS)
│   ├── entities/             # Domain entities (DDD)
│   │   └── common.go
│   ├── services/             # Domain services (DDD)
│   │   └── user_service.go
│   ├── value_objects/        # Immutable value types
│   │   └── email.go
│   └── errors/               # Domain errors
│       └── errors.go
│
├── ports/                    # INTERFACES (Hexagonal)
│   ├── cache.go              # Cache port interface
│   ├── repository.go          # Repository port interface
│   └── observability.go       # Logger, metrics, tracing ports
│
├── adapters/                  # IMPLEMENTATIONS (Hexagonal)
│   ├── cache/
│   │   ├── memory.go         # In-memory adapter
│   │   └── redis.go          # Redis adapter
│   ├── persistence/
│   │   ├── postgres.go       # Postgres adapter
│   │   └── mongo.go          # MongoDB adapter
│   └── observability/
│       ├── slog.go            # slog adapter
│       └── prometheus.go      # Prometheus adapter
│
├── application/              # USE CASES (Clean Architecture)
│   ├── commands/             # Write operations (CQRS)
│   │   └── create_user.go
│   ├── queries/              # Read operations (CQRS)
│   │   └── get_user.go
│   └── handlers/             # Command/query handlers
│       └── user_handler.go
│
├── presentation/             # ENTRY POINTS
│   ├── cli/                  # CLI commands
│   ├── api/                  # HTTP handlers
│   └── grpc/                 # gRPC handlers
│
└── infrastructure/          # CROSS-CUTTING
    ├── config/               # Configuration
    ├── wire/                 # Dependency injection wiring
    └── middleware/            # HTTP middleware
```

### Port Interface Pattern

Following **Interface Segregation** (SOLID) and **Dependency Inversion** (SOLID):

```go
// ports/repository.go - Ports defined in domain layer

// QueryRepositoryPort is a read-only repository.
type QueryRepositoryPort[T any, ID any] interface {
    FindByID(ctx context.Context, id ID) (T, error)
    FindByFilter(ctx context.Context, filter map[string]any) ([]T, error)
    Exists(ctx context.Context, id ID) (bool, error)
}

// CommandRepositoryPort is a write-only repository.
type CommandRepositoryPort[T any, ID any] interface {
    Save(ctx context.Context, entity T) (T, error)
    Update(ctx context.Context, entity T) (T, error)
    Delete(ctx context.Context, id ID) error
}

// RepositoryPort combines both (for simple cases).
type RepositoryPort[T any, ID any] interface {
    QueryRepositoryPort[T, ID]
    CommandRepositoryPort[T, ID]
}
```

### Adapter Pattern

Following **Adapter** pattern and **Open/Closed** (SOLID):

```go
// adapters/cache/memory.go - In-memory implementation

type MemoryCache struct {
    mu    sync.RWMutex
    items map[string][]byte
    ttl   time.Duration
}

func (c *MemoryCache) Get(ctx context.Context, key string) ([]byte, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    val, ok := c.items[key]
    if !ok {
        return nil, ports.ErrCacheMiss{Key: key}
    }
    return val, nil
}

func (c *MemoryCache) Set(ctx context.Context, key string, val []byte, ttl time.Duration) error {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.items[key] = val
    return nil
}

// Satisfies ports.CachePort
```

### Dependency Injection Wiring

Following **Composition Root** pattern and **Dependency Inversion**:

```go
// infrastructure/wire/wire.go

type Container struct {
    UserService *domain.UserService

    // Adapters
    userRepo    ports.RepositoryPort[*domain.User, string]
    cache       ports.CachePort
    logger      ports.LoggerPort
}

func Wire(cfg *Config) (*Container, error) {
    // Create adapters
    var logger ports.LoggerPort
    if cfg.UseSlog {
        logger = &observability.SlogAdapter{}
    } else {
        logger = &observability.LogAdapter{}
    }

    var cache ports.CachePort
    if cfg.UseRedis {
        cache = &cache.RedisAdapter{Redis: redisClient}
    } else {
        cache = &cache.MemoryAdapter{}
    }

    var userRepo ports.RepositoryPort[*domain.User, string]
    if cfg.UsePostgres {
        userRepo = &persistence.PostgresUserRepo{DB: db}
    }

    // Create domain services with injected dependencies
    userSvc := domain.NewUserService(userRepo, cache, logger)

    return &Container{
        UserService: userSvc,
        userRepo:    userRepo,
        cache:       cache,
        logger:      logger,
    }, nil
}
```

## Consequences

### Positive
- **Testability**: Mock any adapter for unit tests
- **Flexibility**: Swap implementations without changing domain
- **Clarity**: Clear boundaries between layers
- **Parallel Development**: Teams can work on different layers
- **Microservice Ready**: Easy to extract as standalone service

### Negative
- **Boilerplate**: More interfaces to maintain
- **Indirection**: Harder to trace through layers
- **Learning Curve**: Team needs to understand patterns

## Implementation Plan

### Phase 1: Define Ports
- [ ] Define `CachePort` interface
- [ ] Define `RepositoryPort` interfaces
- [ ] Define `LoggerPort`, `MetricsPort`, `TracerPort`
- [ ] Create ADR for each port

### Phase 2: Extract Adapters
- [ ] Create `adapters/cache/` directory
- [ ] Implement `MemoryCache` adapter
- [ ] Implement `RedisCache` adapter
- [ ] Implement `SlogAdapter` for logging

### Phase 3: Create Application Layer
- [ ] Define command/query handlers
- [ ] Implement use cases
- [ ] Add DTOs for input/output

### Phase 4: Migrate Existing Code
- [ ] Update packages to use ports
- [ ] Add integration tests
- [ ] Deprecate old entry points

### Phase 5: Polish
- [ ] Add property-based tests
- [ ] Add mutation tests
- [ ] Document migration guide

## Alternatives Considered

1. **Pure Clean Architecture** - rejected, too many layers for library
2. **Onion Architecture** - similar, but hexagonal is more explicit for ports
3. **Full DDD with Aggregates** - overkill for infrastructure library

## References

- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Ports & Adapters](https://jmg.im/coding/honeycomb-architecture/)
- [Go dependency injection](https://github.com/google/wire)
