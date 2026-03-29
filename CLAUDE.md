# Project Instructions

**This project is managed through AgilePlus.**

## AgilePlus Mandate

All work MUST be tracked in AgilePlus:
- Reference: `/Users/kooshapari/CodeProjects/Phenotype/repos/AgilePlus`
- CLI: `cd /Users/kooshapari/CodeProjects/Phenotype/repos/AgilePlus && agileplus <command>`

## Work Requirements

1. **Check for AgilePlus spec before implementing** (`repos/AgilePlus/scripts/list-features.sh` or `repos/AgilePlus/kitty-specs/<slug>/`).
2. **Create or revise spec**: `agileplus specify [--feature <kebab-slug>] [--from-file path.md] [--force]`
3. **Delivery**: `agileplus validate --feature <slug>`, `agileplus plan --feature <slug>`, `agileplus implement --feature <slug>`, and `agileplus queue list` / `agileplus cycle list` as needed.
4. **No code without corresponding AgilePlus spec**

---

# phenotype-go-kit

## Project

Go module containing generic infrastructure packages extracted from the Phenotype ecosystem.
Each package is independent with minimal dependencies.

## Stack

- **Language**: Go 1.24+
- **Test**: `go test ./...`
- **Lint**: `go vet ./...`
- **Format**: `gofumpt`

## Structure

### Top-Level Directories

```
phenotype-go-kit/
├── contracts/              # Port interfaces (not implementation)
│   ├── ports/
│   │   ├── inbound/       # Driving ports (UseCase, CommandHandler, QueryHandler)
│   │   └── outbound/      # Driven ports (Repository, Cache, EventPublisher)
│   ├── models/            # Domain models and events
│   └── plugins/           # Plugin system interfaces
│
├── plugins/                # Plugin implementations
│   └── embeddings/        # AI embeddings providers (OpenAI, Ollama)
│
├── cache/                  # Cache implementation (Redis)
│   ├── adapter/           # Secondary adapters (RedisCacheAdapter)
│   └── service/           # Application services (CQRS handlers)
│
├── auth/                  # Authentication
├── db/                    # Database utilities
├── logging/               # Logging utilities
└── ...
```

## Design Principles

This repository follows **Hexagonal Architecture** with **SOLID**, **GRASP**, **Law of Demeter**, and **SoC**.

### SOLID Principles

| Principle | Description |
|-----------|-------------|
| **SRP** | Single Responsibility - one reason to change |
| **OCP** | Open/Closed - extend, don't modify |
| **LSP** | Liskov Substitution - subtypes substitutable |
| **ISP** | Interface Segregation - small, focused interfaces |
| **DIP** | Dependency Inversion - depend on abstractions |

### GRASP Patterns

| Pattern | Applied As |
|---------|------------|
| Information Expert | Entity knows its own ID |
| Creator | Factory creates repositories |
| Controller | UseCaseHandler receives input |
| Low Coupling | Ports define minimal interfaces |
| High Cohesion | Package-by-feature structure |

### Law of Demeter (LoD)

> Only talk to immediate collaborators.

```go
// GOOD: Tell, don't ask
service.Call(ctx, request)

// BAD: Train-wreck
user.GetSession().GetToken().Use()
```

## ADRs (Architecture Decision Records)

| ID | Title | Status |
|----|-------|--------|
| ADR-001 | Hexagonal Architecture | Accepted |
| ADR-002 | xDD Methodologies Reference | Reference |
| ADR-003 | Top-Level Directory Structure | Accepted |
| ADR-004 | Plugin System & Extensibility | Accepted |
| ADR-005 | AI Embeddings Plugin System | Accepted |
| ADR-006 | Design Principles (SOLID, GRASP, LoD) | Accepted |

## Shared Governance Protocols

These governance blocks are maintained centrally:
- Worktree discipline, reuse protocol, git delivery, stability, CI, child-agent delegation
- Source: `KooshaPari/thegent` -> `templates/claude/governance-blocks/`
- Do not duplicate these blocks here — reference the source instead.

<!-- governance: see thegent/templates/claude/governance-blocks/ for shared protocols -->

## Worktree Discipline

- Feature work goes in `.worktrees/<topic>/`
- Legacy `PROJECT-wtrees/` and `repo-wtrees/` roots are for migration only and must not receive new work.
- Canonical repository remains on `main` for final integration and verification.
