# ADR-007: API Design Guidelines

## Status
**Accepted** | 2024-01-15

## Context

We need consistent, intuitive APIs that follow Go conventions and are easy to use correctly.

## Decision

### 1. Constructor Functions

```go
// New prefix for constructors
func NewService(cfg Config) *Service
func NewClient(opts ...Option) *Client
func NewRingBuffer[T any](capacity int) *RingBuffer[T]
```

### 2. Option Pattern for Configuration

```go
// Functional options for flexible configuration
type Option func(*Server)

func WithTimeout(d time.Duration) Option {
    return func(s *Server) {
        s.timeout = d
    }
}

func WithMiddleware(mw ...Middleware) Option {
    return func(s *Server) {
        s.middleware = mw
    }
}

// Usage
srv := NewServer(
    WithTimeout(30*time.Second),
    WithMiddleware(mw1, mw2),
)
```

### 3. Error Handling

```go
// Return errors, don't panic
func (r *Repository) FindByID(id string) (*User, error) {
    user, err := r.db.Query(...)
    if err != nil {
        return nil, fmt.Errorf("find user: %w", err)  // Wrap errors
    }
    return user, nil
}

// Sentinel errors for known conditions
var (
    ErrNotFound     = errors.New("not found")
    ErrExists       = errors.New("already exists")
    ErrInvalidInput = errors.New("invalid input")
)

// Check with errors.Is
if errors.Is(err, ErrNotFound) {
    // Handle not found
}
```

### 4. Context Propagation

```go
// Always accept Context as first parameter
func FindByID(ctx context.Context, id string) (*User, error)

// Don't use context.Background() internally
// Pass context from caller

// Timeout patterns
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()
```

### 5. Method Chaining (Builder Pattern)

```go
// For complex object construction
query := NewQuery().
    Where("status = ?", "active").
    OrderBy("created_at DESC").
    Limit(10).
    Offset(20)

// For queries
results, err := client.Query().
    Select("id", "name", "email").
    From("users").
    Where("active = true").
    Execute(ctx)
```

### 6. Nil Safety

```go
// Return nil slices, not empty slices
func (r *Repo) FindAll() ([]*User, error) {
    if len(r.users) == 0 {
        return nil, nil  // NOT []User{}
    }
    return r.users, nil
}

// Document nil behavior
// Return value is nil if not found, not an empty slice
```

### 7. Documentation

```go
// Package-level doc
// Package cache provides thread-safe caching with automatic invalidation.
//
// # Usage
//
//   cache := cache.NewRedisCache(client)
//   err := cache.Set(ctx, "key", value, ttl)
//
//   val, err := cache.Get(ctx, "key")
//
// # Eviction
//
// Items are evicted when TTL expires or when manually deleted.
package cache

// Function-level doc
// Get retrieves a value from the cache.
//
// Returns nil if key doesn't exist. The returned slice is a copy;
// modifying it won't affect the cache.
func (c *Cache) Get(ctx context.Context, key string) ([]byte, error)
```

## Consequences

### Positive
- Consistent API across packages
- Easy to use correctly
- Harder to misuse
- Self-documenting code

### Negative
- More upfront design
- Option pattern adds indirection

## References
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://go.dev/doc/effective_go)
