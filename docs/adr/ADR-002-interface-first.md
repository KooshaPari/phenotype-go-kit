# ADR-002: Interface-First Design

## Status
**Accepted** | 2024-01-15

## Context

We need clear contracts between layers to enable testing and swapping implementations (hexagonal architecture).

**Problem:** Without interfaces, tests require real implementations (DB, Redis, etc.), making unit testing difficult.

## Decision

### Core Principle
Define interfaces in the **consumer** package, not the **producer**.

### Pattern: Port Interface

```go
// application/service/user_service.go (consumer defines interface)
package service

// Port interface - defined by consumer
type UserRepository interface {
    FindByID(ctx context.Context, id string) (*User, error)
    Save(ctx context.Context, user *User) error
}

// Consumer depends on abstraction
type UserService struct {
    repo UserRepository  // interface, not concrete type
}
```

### Implementation in Infrastructure

```go
// infrastructure/persistence/postgres/user_repo.go (producer implements)
package postgres

type PostgresUserRepository struct {
    db *sql.DB
}

// Satisfies service.UserRepository interface
func (r *PostgresUserRepository) FindByID(ctx context.Context, id string) (*User, error) {
    // Implementation
}
```

## Consequences

### Positive
- Easy mocking in tests
- Swappable implementations
- Clear contracts
- Supports hexagonal architecture

### Negative
- Interface lives in consumer package (unusual for Go)
- Need to be deliberate about interface placement

## Implementation Guidelines

1. **Small interfaces** - aim for 1-3 methods
2. **Interface segregation** - multiple small interfaces > one large
3. **Return concrete types** - prefer returning interfaces only when needed
4. **Accept interfaces** - functions should accept interfaces, return concrete types

## References
- [Go Blog: Interfaces](https://go.dev/blog/interfaces)
- [Hexagonal Architecture: Ports](https://alistair.cockburn.us/hexagonal-architecture/)
