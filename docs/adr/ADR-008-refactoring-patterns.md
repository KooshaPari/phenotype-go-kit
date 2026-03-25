# ADR-008: Refactoring Patterns

## Status
**Accepted** | 2024-01-15

## Context

Code needs to evolve. We need patterns for safe, incremental refactoring.

## Decision

### 1. Branch by Abstraction

When changing an implementation in-place:

```go
// BEFORE: Direct implementation
func (s *Service) Process() error {
    return s.legacy.Process()
}

// STEP 1: Add abstraction layer
type Processor interface {
    Process() error
}

// STEP 2: Make injectable
func (s *Service) SetProcessor(p Processor) {
    s.processor = p
}

// STEP 3: Old implementation as default
func NewService() *Service {
    return &Service{
        processor: &LegacyProcessor{},  // default
    }
}

// STEP 4: New implementation can now be injected
// No callers need to change yet

// LATER: Change default
func NewService() *Service {
    return &Service{
        processor: &NewProcessor{},  // switched
    }
}
```

### 2. Expand-Contract Pattern

For breaking changes:

```go
// PHASE 1: Expand - add new API
func (s *Service) ProcessV2(ctx context.Context, input Input) (Output, error) {
    // New implementation
}

// PHASE 2: Migrate callers
// Update all callers to V2

// PHASE 3: Contract - remove old API
// Once all migrated, remove Process()
```

### 3. Parallel Run

For risky migrations:

```go
func (s *Service) Process(ctx context.Context, input Input) (Output, error) {
    // Run both implementations
    result, err := s.newImpl.Process(ctx, input)
    
    // Compare results in dev/staging
    if s.enableComparison {
        oldResult, _ := s.oldImpl.Process(ctx, input)
        if !cmp.Equal(result, oldResult) {
            log.Printf("DIFF: new=%v old=%v", result, oldResult)
        }
    }
    
    return result, err
}
```

### 4. Strangler Fig

Incrementally replace a system:

```
┌─────────────────────────────────────────────┐
│              API Gateway                     │
└─────────────────┬───────────────────────────┘
                  │
        ┌─────────┴─────────┐
        ▼                   ▼
┌───────────────┐   ┌───────────────┐
│  Old System   │   │  New System   │
│  (Legacy)     │   │  (Greenfield) │
└───────────────┘   └───────┬───────┘
                            │
                    Traffic increases as
                    features are migrated
```

### 5. Extract Method Object

For long methods:

```go
// BEFORE: Long method
func (s *Service) Process(order Order) error {
    // 100 lines of logic
}

// AFTER: Extract to class
type OrderProcessor struct {
    service *Service
    order   Order
}

func (op *OrderProcessor) Execute() error {
    // Extract logic here
    op.validate()
    op.calculate()
    op.save()
    op.notify()
}

func (s *Service) Process(order Order) error {
    return (&OrderProcessor{service: s, order: order}).Execute()
}
```

### 6. Introduce Parameter Object

For long parameter lists:

```go
// BEFORE
func CreateUser(name, email, phone, address, city, state, zip, country string) error

// AFTER
type CreateUserParams struct {
    Name    string
    Email   string
    Phone   string
    Address Address
}

type Address struct {
    Street  string
    City    string
    State   string
    Zip     string
    Country string
}

func CreateUser(ctx context.Context, params CreateUserParams) error
```

### 7. Replace Error with nil

For methods returning errors:

```go
// BEFORE: Error as special value
func (r *Repo) FindByID(id string) *User {
    user := r.db.Find(id)
    if user == nil {
        return nil  // Confusing: nil vs error?
    }
    return user
}

// AFTER: Error as second return
func (r *Repo) FindByID(id string) (*User, error) {
    user, err := r.db.Find(id)
    if err != nil {
        return nil, fmt.Errorf("find: %w", err)
    }
    if user == nil {
        return nil, ErrNotFound
    }
    return user, nil
}
```

## Consequences

### Positive
- Safe incremental changes
- Easy rollback
- Reduced risk
- Continuous delivery

### Negative
- More upfront planning
- Temporary complexity
- Longer timelines

## References
- [Martin Fowler Refactoring Catalog](https://refactoring.com/catalog/)
- [Working Effectively with Legacy Code](https://www.legacybook.com/)
