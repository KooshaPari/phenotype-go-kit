# ADR-002: GRASP Patterns & xDD Methodologies Reference

**Status**: Reference
**Date**: 2026-03-25
**Type**: Best Practices Guide

---

## GRASP Patterns (General Responsibility Assignment Software Patterns)

| Pattern | Description | Applied In |
|---------|-------------|------------|
| **Information Expert** | Assign responsibility to class with most info needed | Domain entities |
| **Creator** | Class A creates B if A contains/aggregates/uses B | Factory pattern |
| **Controller** | First object beyond UI layer handling events | UseCase handlers |
| **Low Coupling** | Minimize dependencies between classes | Ports interfaces |
| **High Cohesion** | Keep related responsibilities together | Domain services |
| **Polymorphism** | Handle variations based on type | Adapter dispatch |
| **Pure Fabrication** | Create artificial class for high cohesion | Application services |
| **Indirection** | Introduce mediator to decouple | Ports/adapters |
| **Protected Variations** | Isolate unstable elements | Plugin system |

---

## xDD Methodologies (50+)

### Development Methodologies

```
┌─────────┬────────────────────────────────────────┐
│ TDD     │ Test-Driven Development                │
│ BDD     │ Behavior-Driven Development            │
│ DDD     │ Domain-Driven Design                    │
│ ATDD    │ Acceptance Test-Driven Development     │
│ SDD     │ Specification-Driven Development       │
│ FDD     │ Feature-Driven Development             │
│ CDD     │ Contract-Driven Development            │
│ IDD     │ Interaction-Driven Development         │
│ MDD     │ Model-Driven Development               │
│ RDD     │ Repository-Driven Development          │
│ PDD     │ Pattern-Driven Development             │
│ QFDD    │ Quality-Function-Driven Development   │
│ StepDD  │ Stepwise Development                   │
└─────────┴────────────────────────────────────────┘
```

### Design Principles

```
┌─────────┬────────────────────────────────────────┐
│ DRY     │ Don't Repeat Yourself                 │
│ KISS    │ Keep It Simple, Stupid                │
│ YAGNI   │ You Aren't Gonna Need It              │
│ SOLID   │ Single Responsibility, OCP, LSP, ISP, DIP │
│ GRASP   │ General Responsibility Assignment      │
│ LoD     │ Law of Demeter (Principle of Least Knowledge) │
│ SoC     │ Separation of Concerns                │
│ CoC     │ Convention over Configuration         │
│ PoLA    │ Principle of Least Astonishment       │
│ AIP     │ Interface Segregation Principle       │
│ SAP     │ Stable Abstractions Principle          │
│ SDP     │ Stable Dependencies Principle         │
│ REP     │ Reuse-Release Equivalence Principle   │
│ CCP     │ Common Closure Principle              │
│ CRP     │ Common Reuse Principle                │
└─────────┴────────────────────────────────────────┘
```

### Architectural Patterns

```
┌──────────────┬────────────────────────────────────────┐
│ Clean        │ Uncle Bob's Clean Architecture         │
│ Hexagonal    │ Ports & Adapters                       │
│ Onion        │ Inside-out dependency flow             │
│ CQRS         │ Command Query Responsibility Segregation│
│ EDA          │ Event-Driven Architecture              │
│ Event Sourcing│ Store events, derive state            │
│ Microservices │ Distributed small services             │
│ Modular      │ Package by feature/module              │
│ Layered      │ Presentation, Domain, Infrastructure  │
│ Pipe-Filter  │ Data processing pipeline               │
│ Broker       │ Mediated communication                  │
│ MVC          │ Model-View-Controller                   │
│ MVP          │ Model-View-Presenter                    │
│ MVVM         │ Model-View-ViewModel                    │
│ ADR          │ Architecture Decision Records          │
└──────────────┴────────────────────────────────────────┘
```

### Quality Assurance

```
┌─────────────────┬────────────────────────────────────┐
│ Property-Based  │ Fast-check/proptest-style testing  │
│ Mutation        │ Mutate code to verify test quality │
│ Contract        │ Pre/post condition testing         │
│ Shift-Left      │ Move testing earlier in lifecycle  │
│ Chaos           │ Deliberate failure testing         │
│ Performance     │ Load, stress, spike testing        │
│ Security        │ Pen testing, SAST, DAST            │
│ Smoke           │ Quick sanity tests                 │
│ Regression      │ Prevent feature breakage           │
│ Canary          │ Gradual rollout testing            │
│ A/B             │ Comparative testing                │
│ Exploratory     │ Manual creative testing            │
│ Snapshot        │ Visual regression testing          │
└─────────────────┴────────────────────────────────────┘
```

### Process & DevOps

```
┌─────────┬────────────────────────────────────────┐
│ CI/CD   │ Continuous Integration/Delivery       │
│ Agile   │ Iterative incremental development     │
│ Scrum   │ Sprint-based project management       │
│ Kanban  │ Flow-based task management            │
│ XP      │ Extreme Programming practices         │
│ SAFe    │ Scaled Agile Framework                │
│ LeSS    │ Large-Scale Scrum                      │
│ ShapeUp │ Basecamp's methodology                 │
│ Ops     │ Operations integration                │
│ SRE     │ Site Reliability Engineering         │
│ DevSecOps│ Security in DevOps pipeline          │
│ GitOps  │ Git-based infrastructure management   │
│ MLOps   │ ML deployment and operations           │
│ FinOps  │ Cloud cost optimization               │
└─────────┴────────────────────────────────────────┘
```

### Documentation & Communication

```
┌──────────┬───────────────────────────────────────┐
│ ADRs     │ Architecture Decision Records         │
│ RFC      │ Request for Comments                  │
│ Design   │ Technical design documents            │
│ Runbooks │ Operational procedures                │
│ SpecDD   │ Specification-Driven Development     │
│ StoryDD  │ Story-driven documentation           │
│ TraceDD  │ Trace-based requirements             │
│ Diátaxis │ User guide types (tutorials/how-to)  │
│ ADR      │ Architecture Decision Records        │
│ Changelog│ Version history documentation        │
│ API Spec │ OpenAPI/Swagger specifications       │
└──────────┴───────────────────────────────────────┘
```

### Emerging & AI-Assisted

```
┌─────────┬────────────────────────────────────────┐
│ AI-DD   │ AI-assisted Development               │
│ Prompt-DD│ Prompt-driven specification          │
│ LLM-Assist│ Large Language Model assistance    │
│ Copilot │ AI pair programming                   │
│ Claude-DD│ Anthropic Claude-assisted development│
│ RAG     │ Retrieval-Augmented Generation       │
│ MLOps   │ ML operational practices              │
│ AIOps   │ AI for IT operations                  │
└─────────┴────────────────────────────────────────┘
```

---

## Applied Examples in phenotype-go-kit

### TDD + BDD + ATDD Flow
```
1. ATDD: Write acceptance criteria (Given-When-Then)
2. TDD: Write failing unit tests
3. BDD: Automate acceptance tests with ginkgo
4. Refactor: Apply SOLID
```

### Hexagonal + Clean Architecture
```
┌─────────────────────────────────────────┐
│           Adapters (Outer)              │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐  │
│  │  REST   │ │  gRPC   │ │  CLI    │  │
│  └────┬────┘ └────┬────┘ └────┬────┘  │
│       │          │          │         │
│       └──────────┼──────────┘         │
│            Ports │                      │
│       ┌───────────┼───────────┐         │
│       │           │           │         │
│  ┌────┴────┐ ┌────┴────┐ ┌────┴────┐  │
│  │  Use    │ │Command  │ │ Query   │  │
│  │  Case   │ │Handler  │ │ Handler │  │
│  └────┬────┘ └────┬────┘ └────┬────┘  │
│       │           │           │         │
│       └───────────┬┴───────────┘        │
│               Domain │                   │
│       ┌─────────────┼─────────────┐    │
│       │             │             │     │
│  ┌────┴────┐  ┌─────┴────┐ ┌──────┴───┐ │
│  │Entities │  │ Services │ │  Value   │ │
│  │         │  │          │ │  Objects │ │
│  └─────────┘  └──────────┘ └──────────┘ │
│              Domain (Inner)              │
└─────────────────────────────────────────┘
```

### GRASP + SOLID Application

| GRASP | SOLID | Example |
|-------|-------|---------|
| Creator | (n/a) | `RepositoryFactory` creates repositories |
| Controller | (n/a) | `UseCaseHandler` receives input |
| Information Expert | SRP | `Entity.GetID()` knows its own ID |
| Low Coupling | ISP, DIP | Ports define minimal interfaces |
| High Cohesion | SRP | `UserService` only user operations |
| Polymorphism | OCP | `Handler` interface for all transports |
| Pure Fabrication | ISP | `ApplicationService` for orchestration |
| Indirection | DIP | `Repository` interface decouples |
| Protected Variations | OCP, DIP | Plugin system isolates changes |

---

## Decision Criteria Matrix

```
When to apply which pattern:

┌──────────────────────────────────────────────────────────────┐
│ Need testability + multiple transports?    → Hexagonal      │
│ Complex domain with business rules?         → DDD            │
│ Need contract verification?                → CDD + Contract │
│ Many variations of algorithms?             → Strategy + TDD  │
│ Need clear boundaries?                     → Modular + ADR   │
│ Rapid iteration?                           → Agile + CI/CD   │
│ AI-assisted development?                   → AI-DD + PromptDD│
└──────────────────────────────────────────────────────────────┘
```

---

## References

- [GRASP Patterns - Craig Larman](https://www.craiglarman.com/wiki/index.php?title=GRASP_Patterns)
- [xDD Methods - Wikipedia](https://en.wikipedia.org/wiki/Category:Software_developmentPhilosophies)
- [Clean Architecture](blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design - Eric Evans](https://domainlanguage.com/ddd/)
