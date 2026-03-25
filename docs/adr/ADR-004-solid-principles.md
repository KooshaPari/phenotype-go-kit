# ADR-004: SOLID Principles Implementation

## Status
**Accepted** | 2024-01-15

## Context

We need consistent code quality across a multi-contributor project. SOLID provides a checklist for maintainability.

## Decision

### S - Single Responsibility Principle
Each type has **one reason to change**.

```go
// BAD - handles data, caching, AND logging
type UserService struct {
    db *sql.DB
    cache *redis.Client
    logger *slog.Logger
}

// GOOD - each has one job
type UserRepository struct { db *sql.DB }           // Data access
type UserCache struct { client *redis.Client }      // Caching
type UserLogger struct { logger *slog.Logger }    // Logging
```

### O - Open/Closed Principle
Open for extension, closed for modification.

```go
// BAD - adding new formats requires modifying this
func (e *Exporter) Export(format string) {
    if format == "json" { /* ... */ }
    if format == "csv" { /* ... */ }
}

// GOOD - extend via interface
type Exporter interface { Export(io.Writer) error }
type JSONExporter struct {}
type CSVExporter struct {}
```

### L - Liskov Substitution Principle
Subtypes must be substitutable for base types.

```go
// BAD - violates Liskov
type IntQueue struct { items []int }
func (q *IntQueue) Dequeue() int {
    if len(q.items) == 0 { return 0 } // Returns zero, can't distinguish from actual zero
}

// GOOD
type IntQueue struct { items []int }
func (q *IntQueue) Dequeue() (int, error) {
    if len(q.items) == 0 { return 0, ErrEmpty }
    item := q.items[0]
    q.items = q.items[1:]
    return item, nil
}
```

### I - Interface Segregation Principle
Many small interfaces > one large interface.

```go
// BAD - callers forced to implement unused methods
type Persister interface {
    Save() error
    Delete() error
    Query() ([]Item, error)
    Connect() error
    Disconnect() error
}

// GOOD - small, focused interfaces
type Saver interface { Save() error }
type Deleter interface { Delete() error }
type Querier interface { Query() ([]Item, error) }
```

### D - Dependency Inversion Principle
High-level modules depend on abstractions, not concretions.

```go
// BAD - direct dependency on concrete type
type UserService struct {
    repo *PostgresUserRepository  // concrete, not interface
}

// GOOD - depend on abstraction
type UserService struct {
    repo UserRepository  // interface
}
```

## Consequences

### Positive
- Easier to test (mock interfaces)
- Easier to extend (add implementations)
- Clearer responsibility boundaries
- Reduced risk when changing code

### Negative
- More interfaces to manage
- Slight overhead for small projects
- Requires discipline to maintain

## Checklist

| Principle | Question |
|-----------|----------|
| SRP | Does this type have only one reason to change? |
| OCP | Can I add new behavior without changing existing code? |
| LSP | Can I substitute any implementation without breaking? |
| ISP | Are all interface methods actually used by consumers? |
| DIP | Do high-level modules depend on abstractions? |

## References
- [SOLID Go Design](https://dave.cheney.net/2016/08/20/solid-go-design)
