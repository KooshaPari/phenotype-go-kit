# xDD Methodologies & Best Practices Reference

## 50+ Development Methodologies Reference

### Development Methodologies (12)

| Acronym | Name | Core Principle |
|---------|------|----------------|
| TDD | Test-Driven Development | Write tests before code |
| BDD | Behavior-Driven Development | Define behavior via scenarios |
| DDD | Domain-Driven Design | Model based on business domain |
| ATDD | Acceptance Test-Driven Development | Tests from acceptance criteria |
| SDD | Specification-Driven Development | Specifications as first-class artifacts |
| FDD | Feature-Driven Development | Feature-centric iterative building |
| CDD | Contract-Driven Development | API contracts define interactions |
| IDD | Intent-Driven Development | Explicit intent in code |
| MDD | Model-Driven Development | Models as primary artifacts |
| RDD | README-Driven Development | Docs as project foundation |
| AI-DD | AI-Assisted Development | AI augments development |
| Prompt-DD | Prompt-Driven Development | Prompts drive AI coding |

### Design Principles (15)

| Acronym | Name | Core Principle |
|---------|------|----------------|
| DRY | Don't Repeat Yourself | Single source of truth |
| KISS | Keep It Simple, Stupid | Prefer simplicity |
| YAGNI | You Aren't Gonna Need It | No premature features |
| SOLID | Single Responsibility, Open/Closed, Liskov, Interface Segregation, Dependency Inversion | OOP design principles |
| GRASP | General Responsibility Assignment Software Patterns | Responsibility assignment |
| LoD | Law of Demeter | Minimize coupling |
| SoC | Separation of Concerns | Divide by concern |
| CoC | Convention over Configuration | Sensible defaults |
| PoLA | Principle of Least Astonishment | Expected behavior |
| BDUF | Big Design Up Front | Design before coding |
| WET | Write Everything Twice | Anti-DRY anti-pattern |
| MVP | Minimum Viable Product | Ship core functionality |
| MVM | Model-View-ViewModel | UI separation |
| RCP | Release Candidate Principle | Stable release candidates |
| CRP | Customer Ratio Principle | Customer-centric metrics |

### Architecture Patterns (15)

| Pattern | Description | Use Case |
|---------|-------------|----------|
| Clean Architecture | Layers: domain, application, infrastructure, UI | Enterprise applications |
| Hexagonal Architecture | Ports & Adapters isolation | Plugable systems |
| Onion Architecture | Layered from core outward | Domain-centric design |
| CQRS | Command Query Responsibility Segregation | Read/write optimization |
| EDA | Event-Driven Architecture | Async, reactive systems |
| Microservices | Small, independent services | Scalable systems |
| Event Sourcing | Store events, not state | Audit, replay |
| Serverless | Function-as-a-Service | Event-driven workloads |
| SAGA | Distributed transactions | Long-running processes |
| Strangler Fig | Incrementally replace systems | Legacy modernization |
| Anti-Corruption | Protect domain from external | Integration |
| Facade | Simplified interface | Complex subsystems |
| Decorator | Add behavior dynamically | Extensible design |
| Proxy | Control access | Lazy loading, security |
| Repository | Abstract data access | Persistence independence |

### Quality Assurance (12)

| Method | Description | Goal |
|--------|-------------|------|
| Property-Based Testing | Test invariants across inputs | Fuzzing, proptest |
| Mutation Testing | Mutate code, verify tests catch | Test quality |
| Contract Testing | Verify API contracts | Integration |
| Shift-Left Testing | Test earlier in lifecycle | Faster feedback |
| Chaos Engineering | Inject failures intentionally | Resilience |
| Canary Testing | Gradual rollouts | Risk reduction |
| A/B Testing | Compare variants | Data-driven decisions |
| SRE Practices | SLIs/SLOs/SLAs | Reliability |
| Code Review | Peer inspection | Knowledge sharing |
| Static Analysis | Automated code inspection | Quality gates |
| Dynamic Analysis | Runtime behavior analysis | Memory leaks |
| Performance Testing | Load, stress, spike tests | Scalability |

### Documentation (8)

| Method | Description | Output |
|--------|-------------|--------|
| ADRs | Architecture Decision Records | Decision history |
| RFCs | Request for Comments | Design proposals |
| Design Docs | Technical specifications | Implementation guides |
| Runbooks | Operational procedures | Incident response |
| SpecDD | Specification-Driven Development | Living specs |
| Storybook | UI component catalog | Visual testing |
| API Docs | OpenAPI/Swagger specs | API reference |
| Changelog | Version history | Release notes |

### Process & Agile (10)

| Method | Description | Focus |
|--------|-------------|-------|
| DevOps | Dev/Ops collaboration | Delivery |
| CI/CD | Continuous Integration/Delivery | Automation |
| Agile | Iterative development | Adaptability |
| Scrum | Sprint-based framework | Predictability |
| Kanban | Flow-based delivery | Throughput |
| XP | Extreme Programming | Quality |
| Lean | Eliminate waste | Efficiency |
| SAFe | Scaled Agile Framework | Enterprise |
| DevEx | Developer Experience | Productivity |
| Inner Loop | Personal dev cycle | Speed |

### Emerging & AI-Augmented (10)

| Method | Description | Tools |
|--------|-------------|-------|
| AI-DD | AI-Assisted Development | Claude, Copilot |
| Prompt-DD | Prompt-Driven Development | AI agents |
| Story-DD | User story decomposition | AI generation |
| Trace-DD | Distributed trace-driven | Telemetry-first |
| Test-DD | AI-generated tests | Property-based |
| Refactor-DD | AI-guided refactoring | Smell detection |
| Review-DD | AI code review | Style, bugs |
| Doc-DD | AI documentation | Auto-generate |
| Spec-DD | AI specification | From code |
| Commit-DD | AI commit messages | Conventional |

### Git & Versioning (8)

| Method | Description | Practice |
|--------|-------------|----------|
| Trunk-Based | Short-lived branches | Simplicity |
| GitFlow | Feature/release/hotfix | Release mgmt |
| Squash & Merge | Single commit per PR | Clean history |
| Conventional Commits | Semantic messages | Changelog gen |
| Semantic Versioning | MAJOR.MINOR.PATCH | Dependency mgmt |
| Release Branches | Stable release lines | Support |
| Stacked PRs | Dependent PRs | Modular review |
| Fork Strategy | Fork per contributor | Isolation |

---

## Application Checklist

### Pre-Implementation
- [ ] Read existing CLAUDE.md
- [ ] Review current structure
- [ ] Identify domain boundaries
- [ ] Map existing code to layers

### Architecture Implementation
- [ ] Define domain models (DDD)
- [ ] Identify ports/adapters (Hexagonal)
- [ ] Map layers (Clean Architecture)
- [ ] Apply SOLID principles
- [ ] Extract interfaces (Dependency Inversion)

### Quality Implementation
- [ ] Add property-based tests (proptest)
- [ ] Document ADRs for decisions
- [ ] Add contract tests (Pact)
- [ ] Implement chaos testing
- [ ] Add mutation coverage

### Documentation Implementation
- [ ] Create ARCHITECTURE.md
- [ ] Add ADRs/ folder
- [ ] Document APIs with OpenAPI
- [ ] Create runbooks
- [ ] Add codeowners

---

## Hexagonal + Clean Architecture Template

```
src/
├── domain/                 # Core business logic (DDD)
│   ├── entities/          # Domain entities
│   ├── value_objects/     # Immutable value types
│   ├── services/          # Domain services
│   ├── events/            # Domain events
│   ├── exceptions/        # Domain exceptions
│   └── interfaces/        # Repository & port interfaces
│
├── application/           # Use cases (CDD)
│   ├── commands/          # Write operations
│   ├── queries/           # Read operations
│   ├── handlers/          # Command/query handlers
│   ├── services/          # Application services
│   └── dto/               # Data transfer objects
│
├── infrastructure/        # External adapters
│   ├── persistence/       # Database implementations
│   ├── external/         # External service clients
│   ├── messaging/         # Event publishers
│   └── adapters/          # Port implementations
│
└── presentation/         # UI/API layer
    ├── api/               # REST/gRPC endpoints
    ├── cli/               # CLI commands
    └── ui/                # Web UI
```

### Port Interface Pattern

```go
// domain/interfaces/repository.go
type UserRepository interface {
    Save(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id ID) (*User, error)
    FindAll(ctx context.Context) ([]*User, error)
}

// infrastructure/persistence/postgres/user_repo.go
type PostgresUserRepository struct {
    db *sql.DB
}

func (r *PostgresUserRepository) Save(ctx context.Context, user *User) error {
    // Implementation
}
```

---

## SOLID Principles Checklist

| Principle | Checklist |
|-----------|-----------|
| **S**ingle Responsibility | Does the type have one reason to change? |
| **O**pen/Closed | Can you extend without modifying existing code? |
| **L**iskov Substitution | Can subclasses be used interchangeably? |
| **I**nterface Segregation | Are interfaces small and focused? |
| **D**ependency Inversion | Do high-level modules depend on abstractions? |

---

## Testing Pyramid

```
         /\
        /  \       E2E Tests (few, slow)
       /    \
      /------\     Integration Tests (some, medium)
     /        \
    /----------\  Unit Tests (many, fast)
```

### Test Naming Convention (BDD)

```go
func TestUserRepository_Save_WithValidUser_PersistsToDatabase(t *testing.T)
func TestUserRepository_Save_WithDuplicateEmail_ReturnsError(t *testing.T)
func TestUserService_Create_WithValidInput_ReturnsCreatedUser(t *testing.T)
```

---

## ADR Template

```markdown
# ADR-001: Use Hexagonal Architecture

## Status
Accepted

## Context
[Problem statement]

## Decision
[Chosen approach]

## Consequences
### Positive
- [Benefit 1]
- [Benefit 2]

### Negative
- [Drawback 1]

## Alternatives Considered
- [Alternative 1]
- [Alternative 2]
```
