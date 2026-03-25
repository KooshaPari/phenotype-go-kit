# ADR-006: Design Principles - SOLID, GRASP, LoD, SoC

**Status**: Accepted
**Date**: 2026-03-25
**Deciders**: Phenotype Architecture Team

---

## Context

This ADR documents the core design principles applied across Phenotype repositories to ensure maintainability, testability, and extensibility.

---

## SOLID Principles

| Principle | Description | Applied As |
|-----------|-------------|------------|
| **S**ingle Responsibility | A class has one reason to change | Each package does one thing |
| **O**pen/Closed | Open for extension, closed for modification | Use interfaces, not inheritance |
| **L**iskov Substitution | Subtypes must be substitutable | Implement contracts fully |
| **I**nterface Segregation | Prefer small, focused interfaces | Minimal port interfaces |
| **D**ependency Inversion | Depend on abstractions | Ports define contracts |

```
┌────────────────────────────────────────────────────────────────┐
│                      SOLID in Practice                          │
├────────────────────────────────────────────────────────────────┤
│                                                                 │
│  SRP: cache/service/ has CacheService, not UserAndCacheService  │
│  OCP: New adapter implements Port, existing code unchanged     │
│  LSP: MockProvider can replace RealProvider anywhere           │
│  ISP: CachePort has 10 methods, not 50                         │
│  DIP: Domain uses CachePort, not *RedisClient                  │
│                                                                 │
└────────────────────────────────────────────────────────────────┘
```

---

## GRASP Patterns

| Pattern | Description | Applied As |
|---------|-------------|------------|
| **Information Expert** | Assign to class with needed info | Entity knows its own ID |
| **Creator** | A creates B if it contains/uses B | Factory creates repositories |
| **Controller** | First object beyond UI handling events | UseCaseHandler |
| **Low Coupling** | Minimize dependencies | Ports interfaces |
| **High Cohesion** | Related responsibilities together | Package-by-feature |
| **Polymorphism** | Handle variations by type | Adapter dispatch |
| **Pure Fabrication** | Artificial class for cohesion | ApplicationService |
| **Indirection** | Introduce mediator | Ports/Adapters |
| **Protected Variations** | Isolate unstable elements | Plugin system |

---

## Law of Demeter (LoD)

Also known as "Principle of Least Knowledge".

### The Rule

> A method `M` of object `O` may only invoke the methods of these kinds of objects:
> 1. `O` itself
> 2. `M`'s parameters
> 3. Any objects created/instantiated within `M`
> 4. `O`'s direct component objects
> 5. A global variable accessible by `O`

### Bad Example (Violates LoD)

```go
// DON'T: Train-wreck pattern
user := auth.GetUser()
session := user.GetSession()
token := session.GetToken()
err := api.Call(token)

// ALSO DON'T: Global state
globalCache.Set("key", value) // No!
```

### Good Example (Follows LoD)

```go
// DO: Tell, don't ask
auth.CallAPI(ctx, "endpoint", request)

// DO: Dependency injection
type Service struct {
    cache CachePort  // Injected, not accessed via chain
}

// DO: Method chaining only on value objects
len := strings.Trim(" hello ").Trim()
```

### Why LoD Matters

```
┌────────────────────────────────────────────────────────────────┐
│                    Law of Demeter Benefits                      │
├────────────────────────────────────────────────────────────────┤
│                                                                 │
│  • Reduced coupling between classes                              │
│  • Better encapsulation                                        │
│  • Easier refactoring                                          │
│  • Improved testability                                        │
│  • Fewer unexpected side effects                                │
│                                                                 │
└────────────────────────────────────────────────────────────────┘
```

---

## Separation of Concerns (SoC)

### Definition

Separate code into distinct sections, each addressing a separate concern.

### Applied Structure

```
phenotype-go-kit/
├── contracts/        # What (interfaces) - business contracts
│   ├── ports/       # Inbound: use cases, outbound: data access
│   └── models/      # Domain entities, value objects
│
├── internal/         # How (private implementation)
│   ├── domain/      # Pure business logic
│   ├── application/ # Use cases, orchestration
│   └── adapters/    # Infrastructure implementations
│
├── plugins/          # Extensions, optional features
└── pkg/             # Public, reusable libraries
```

### Concern Categories

| Concern | Location | Example |
|---------|----------|---------|
| Business Rules | Domain | Validation, calculations |
| Data Access | Outbound Ports | Repository, Cache |
| API Transport | Inbound Ports | REST, gRPC handlers |
| Cross-Cutting | Infrastructure | Logging, tracing |
| Configuration | Init | Wiring dependencies |

---

## Convention Over Configuration (CoC)

### Principle

Establish sensible defaults, require explicit configuration only when needed.

### Applied Examples

```go
// Default config values
type Config struct {
    Addr     string        `default:"localhost:6379"`
    Port     int           `default:"8080"`
    Timeout  time.Duration `default:"30s"`
}

// Functional options for non-defaults
WithAPIKey("secret")
WithTimeout(60 * time.Second)
```

### Benefits

- Reduced boilerplate
- Consistent defaults
- Easy to override when needed

---

## DRY (Don't Repeat Yourself)

### Principle

Every piece of knowledge must have a single, unambiguous representation.

### Applied Examples

| Pattern | Implementation |
|---------|----------------|
| Error definitions | `var ErrKeyNotFound = errors.New("...")` |
| Port interfaces | `contracts/ports/` - defined once |
| Validation | Shared validators in `validation/` |
| Config structs | Single definition, YAML/JSON tags |

---

## KISS (Keep It Simple, Stupid)

### Principle

Favor simplicity over cleverness.

### Guidelines

```
┌────────────────────────────────────────────────────────────────┐
│                      KISS Guidelines                            │
├────────────────────────────────────────────────────────────────┤
│                                                                 │
│  1. Prefer clear over clever                                   │
│  2. Small functions that do one thing                          │
│  3. Obvious naming                                            │
│  4. Avoid premature optimization                               │
│  5. YAGNI - You Aren't Gonna Need It                          │
│                                                                 │
│  Bad: `x = (a|b|c) && d ? e(f(g())) : h(i(j()))`              │
│  Good: Extract to well-named functions                          │
│                                                                 │
└────────────────────────────────────────────────────────────────┘
```

---

## YAGNI (You Aren't Gonna Need It)

### Principle

Don't add functionality until it's actually needed.

### Applied As

- Start with simple implementation
- Add complexity only when required
- Refactor when new requirements emerge
- Don't over-engineer for "future" features

---

## Command-Query Separation (CQS)

### Principle

Methods should either:
- **Command**: Change state, return nothing
- **Query**: Return data, don't change state

### Applied in CQRS

```go
// Command (mutates)
SetCacheHandler() func(ctx context.Context, cmd SetCacheCommand) error

// Query (reads)
GetCacheHandler() func(ctx context.Context, query GetCacheQuery) (string, error)
```

---

## References

- [SOLID Principles - Robert C. Martin](https://en.wikipedia.org/wiki/SOLID)
- [GRASP Patterns - Craig Larman](https://www.craiglarman.com/wiki/index.php?title=GRASP_Patterns)
- [Law of Demeter - Wikipedia](https://en.wikipedia.org/wiki/Law_of_Demeter)
- [Don't Repeat Yourself - Wikipedia](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself)
