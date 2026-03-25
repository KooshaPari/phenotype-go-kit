# xDD Methodologies Compendium (100+ Methodologies)

## Comprehensive Software Engineering Methodology Reference

> **Purpose:** Provide a complete reference of development, design, architecture, quality, process, documentation, and emerging methodologies for comprehensive software engineering.

---

## Table of Contents

1. [Development Methodologies (10)](#development-methodologies)
2. [Design Principles (15)](#design-principles)
3. [Architectural Patterns (15)](#architectural-patterns)
4. [Quality Assurance (10)](#quality-assurance)
5. [Process & Delivery (10)](#process--delivery)
6. [Documentation (10)](#documentation)
7. [Emerging AI-Driven (10)](#emerging-ai-driven)
8. [Data & Persistence (10)](#data--persistence)
9. [Security & Reliability (10)](#security--reliability)
10. [Operations & Observability (10)](#operations--observability)

---

## 1. Development Methodologies (10) {#development-methodologies}

| Acronym | Full Name | Description |
|---------|-----------|-------------|
| **TDD** | Test-Driven Development | Write tests before code; red-green-refactor cycle |
| **BDD** | Behavior-Driven Development | Tests written in natural language (Gherkin) |
| **ATDD** | Acceptance Test-Driven Development | Define acceptance criteria before implementation |
| **FDD** | Feature-Driven Development | Model first, then develop features iteratively |
| **SDD** | Specification-Driven Development | Focus on detailed specifications before coding |
| **CDD** | Consumer-Driven Contracts | API consumers define contract tests |
| **DDD** | Domain-Driven Design | Focus on core domain logic and ubiquitous language |
| **MDD** | Model-Driven Development | Generate code from abstract models |
| **RDD** | Responsibility-Driven Design | Focus on object responsibilities and collaborations |
| **IDD** | Interaction-Driven Development | Design based on user interactions/workflows |

---

## 2. Design Principles (15) {#design-principles}

| Acronym | Full Name | Description |
|---------|-----------|-------------|
| **DRY** | Don't Repeat Yourself | Eliminate duplication across codebase |
| **KISS** | Keep It Simple, Stupid | Prefer simple solutions over complex ones |
| **YAGNI** | You Aren't Gonna Need It | Don't implement features until needed |
| **SOLID** | Single Responsibility, Open/Closed, Liskov Substitution, Interface Segregation, Dependency Inversion | Five principles for OOP design |
| **GRASP** | General Responsibility Assignment Software Patterns | Patterns for assigning responsibilities |
| **LoD** | Law of Demeter | Principle of least knowledge; minimize dependencies |
| **SoC** | Separation of Concerns | Divide code into distinct modules |
| **CoC** | Convention over Configuration | Sensible defaults reduce explicit configuration |
| **PoLA** | Principle of Least Astonishment | Behavior should not surprise users |
| **ISP** | Interface Segregation Principle | Prefer small, specific interfaces |
| **CCP** | Common Closure Principle | Classes that change together should be packaged together |
| **CRP** | Common Reuse Principle | Classes used together should be packaged together |
| **SAP** | Stable Abstractions Principle | Stable packages should be abstract |
| **ADP** | Acyclic Dependencies Principle | No cycles in dependency graph |
| **SVP** | Stable Dependencies Principle | Dependencies should point toward stability |

---

## 3. Architectural Patterns (15) {#architectural-patterns}

| Pattern | Description |
|---------|-------------|
| **Hexagonal Architecture** (Ports & Adapters) | Core business logic isolated with ports; adapters for external systems |
| **Clean Architecture** | Layered architecture with domain at center (entities, use cases, interfaces, frameworks) |
| **Onion Architecture** | Layers radiate outward: domain core, application, infrastructure |
| **CQRS** | Command Query Responsibility Segregation: separate read and write models |
| **Event Sourcing** | Store events, not state; rebuild state from event log |
| **Event-Driven Architecture** | Systems communicate via events; async, loosely coupled |
| **Microservices** | Small, independent services that own their data |
| **Monolithic** | Single deployable unit; simpler for small teams |
| **Modular Monolith** | Monolith with clear module boundaries |
| **Serverless** | Functions as a service; auto-scale, pay-per-use |
| **Service Mesh** | Infrastructure layer for service-to-service communication |
| **Saga Pattern** | Manage distributed transactions across services |
| **Strangler Fig** | Gradually migrate from legacy by replacing pieces |
| **Bulkhead** | Isolate failures; prevent cascading |
| **Sidecar** | Deploy auxiliary components alongside main service |

---

## 4. Quality Assurance (10) {#quality-assurance}

| Acronym | Full Name | Description |
|---------|-----------|-------------|
| **PBT** | Property-Based Testing | Test invariants rather than specific examples |
| **Mutation Testing** | Introduce bugs to verify test effectiveness |
| **Contract Testing** | Verify API compatibility between services |
| **Shift-Left** | Move testing earlier in development lifecycle |
| **Shift-Right** | Test in production environments |
| **Chaos Engineering** | Intentionally break systems to find weaknesses |
| **Fuzz Testing** | Random inputs to find edge cases and vulnerabilities |
| **Snapshot Testing** | Compare output against stored snapshots |
| **Regression Testing** | Ensure new changes don't break existing functionality |
| **Smoke Testing** | Quick sanity checks for critical functionality |

---

## 5. Process & Delivery (10) {#process--delivery}

| Acronym | Full Name | Description |
|---------|-----------|-------------|
| **DevOps** | Development + Operations integration | Break down silos between dev and ops |
| **CI/CD** | Continuous Integration/Delivery/Deployment | Automate build, test, and deployment |
| **Agile** | Iterative development with adaptive planning | Embrace change, deliver frequently |
| **Scrum** | Agile framework with sprints, roles, ceremonies | Fixed-length iterations |
| **Kanban** | Visual workflow management; limit WIP | Continuous flow, no fixed sprints |
| **XP** | Extreme Programming | Pair programming, TDD, short iterations |
| **SAFe** | Scaled Agile Framework | Agile at enterprise scale |
| **LeSS** | Large-Scale Scrum | Scrum scaled across multiple teams |
| **Continuous Release** | Release every change that passes CI | Fast feedback, reduced risk |
| **GitOps** | Git as single source of truth for infrastructure | Declarative, versioned deployments |

---

## 6. Documentation (10) {#documentation}

| Acronym | Full Name | Description |
|---------|-----------|-------------|
| **ADR** | Architecture Decision Record | Document significant architectural decisions |
| **RFC** | Request for Comments | Proposal for discussion before implementation |
| **SpecDD** | Specification-Driven Development | Specifications as executable tests |
| **Design Docs** | Technical design documents | Explain "why" and "how" |
| **Runbooks** | Operational procedures | Step-by-step operational guides |
| **ADRs** | Architecture Decision Records | Capture context, decisions, consequences |
| **API Docs** | API specifications | OpenAPI/Swagger, endpoints, schemas |
| **Code Comments** | Inline documentation | Explain non-obvious code |
| **README** | Project overview | Getting started, quick reference |
| **CHANGELOG** | Version history | Track changes across releases |

---

## 7. Emerging AI-Driven (10) {#emerging-ai-driven}

| Acronym | Full Name | Description |
|---------|-----------|-------------|
| **AI-DD** | AI-Driven Development | AI assists in code generation and review |
| **Prompt-Driven Development** | Natural language prompts generate code | LLM-based implementation |
| **StoryDD** | Story-Driven Development | User stories drive AI-generated code |
| **TraceDD** | Trace-Driven Development | Execution traces guide implementation |
| **Prompt Engineering** | Crafting effective prompts for AI tools | Structured input for better output |
| **AI-Augmented Review** | AI assists in code review | Automated suggestions and analysis |
| **AI-Assisted Refactoring** | AI suggests improvements | Pattern recognition for code quality |
| **Synthetic Data Generation** | AI generates test data | Cover edge cases efficiently |
| **Documentation Generation** | AI generates docs from code | Keep docs in sync |
| **AI-Driven Test Generation** | AI generates tests from requirements | Increase coverage |

---

## 8. Data & Persistence (10) {#data--persistence}

| Acronym | Full Name | Description |
|---------|-----------|-------------|
| **ACID** | Atomicity, Consistency, Isolation, Durability | Transaction guarantees |
| **BASE** | Basically Available, Soft state, Eventually consistent | Distributed systems tradeoff |
| **CAP** | Consistency, Availability, Partition tolerance | Choose 2 of 3 in distributed systems |
| **CRDT** | Conflict-free Replicated Data Types | Merge concurrent edits without conflicts |
| **ORM** | Object-Relational Mapping | Map objects to database tables |
| **Repository Pattern** | Abstract data layer | Decouple domain from persistence |
| **Unit of Work** | Atomic operations | Group related changes |
| **Outbox Pattern** | Reliable event publishing | Store events with business data |
| **Saga Pattern** | Distributed transactions | Choreography or orchestration |
| **Materialized View** | Pre-computed query results | Optimize read performance |

---

## 9. Security & Reliability (10) {#security--reliability}

| Acronym | Full Name | Description |
|---------|-----------|-------------|
| **Zero Trust** | Never trust, always verify | Every request is authenticated |
| **Defense in Depth** | Multiple layers of security | No single point of failure |
| **Fail Open/Closed** | Behavior on security failure | Default to more secure option |
| **Rate Limiting** | Throttle excessive requests | Prevent abuse |
| **Circuit Breaker** | Prevent cascading failures | Trip on repeated failures |
| **Retry Pattern** | Automatic retry with backoff | Handle transient failures |
| **Timeout Pattern** | Set maximum wait times | Prevent indefinite hangs |
| **Bulkhead** | Isolate resources | Prevent resource exhaustion |
| **Health Check** | Monitor service health | Liveness and readiness probes |
| **Graceful Degradation** | Reduce functionality on failure | Stay partially available |

---

## 10. Operations & Observability (10) {#operations--observability}

| Acronym | Full Name | Description |
|---------|-----------|-------------|
| **SLO** | Service Level Objective | Target availability/performance |
| **SLA** | Service Level Agreement | Contractual SLO commitments |
| **SLI** | Service Level Indicator | Actual measured performance |
| **DORA** | DevOps Research and Assessment | Four key metrics: velocity and stability |
| **Telemetry** | Automated data collection | Metrics, logs, traces, events |
| **On-Call** | 24/7 availability response | Incident management |
| **Postmortem** | Incident analysis | Blameless review of failures |
| **Runbook** | Operational procedures | Step-by-step guidance |
| **Chaos Engineering** | Test resilience | Deliberately inject failures |
| **Feature Flags** | Toggle features remotely | Gradual rollouts, quick rollback |

---

## Application to Phenotype Go Kit

### Architecture Applied

```
phenotype-go-kit/
├── domain/                    # DDD: Core business logic
│   ├── entities/             # Domain models
│   ├── valueobjects/         # Immutable value types
│   ├── services/             # Domain services
│   ├── events/               # Domain events
│   └── ports/                # Interfaces (Hexagonal)
│
├── application/              # CDD + SpecDD: Use cases
│   ├── commands/             # Write operations
│   ├── queries/              # Read operations (CQRS)
│   ├── handlers/              # Request handlers
│   └── services/              # Application services
│
├── infrastructure/            # Hexagonal: Adapters
│   ├── persistence/          # Database adapters
│   ├── cache/                # Redis adapters
│   ├── messaging/             # Event bus adapters
│   ├── http/                 # HTTP adapters
│   └── external/             # Third-party adapters
│
├── pkg/                      # Shared utilities
│   ├── config/               # Configuration
│   ├── middleware/           # HTTP middleware
│   └── logging/              # Logging utilities
│
└── cmd/                      # Entry points
    ├── server/               # API server
    └── worker/               # Background worker
```

### Principles Applied

| Principle | Implementation |
|-----------|---------------|
| **Hexagonal** | Ports define interfaces; adapters implement |
| **Clean** | Domain core, no external dependencies |
| **SOLID** | Single responsibility per package |
| **DDD** | Entities, Value Objects, Aggregates, Services |
| **CQRS** | Separate command and query handlers |
| **Event Sourcing** | Domain events for state changes |
| **Repository** | Abstract persistence behind interfaces |
| **TDD** | Tests drive implementation |
| **SpecDD** | Specifications as tests |
| **BDD** | Behavior specs in natural language |

---

## References

- Clean Architecture: Robert C. Martin
- Domain-Driven Design: Eric Evans
- Hexagonal Architecture: Alistair Cockburn
- Implementing Domain-Driven Design: Vaughn Vernon
- Refactoring: Martin Fowler
- Building Microservices: Sam Newman
