# Project Instructions

**This project is managed through AgilePlus.**

## AgilePlus Mandate

All work MUST be tracked in AgilePlus:
- Reference: `/Users/kooshapari/CodeProjects/Phenotype/repos/AgilePlus`
- CLI: `cd /Users/kooshapari/CodeProjects/Phenotype/repos/AgilePlus && agileplus <command>`

## Work Requirements

1. **Check for AgilePlus spec before implementing**
2. **Create spec for new work**: `agileplus specify --title "<feature>" --description "<desc>"`
3. **Update work package status**: `agileplus status <feature-id> --wp <wp-id> --state <state>`
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

## CI Completeness Policy

- Always evaluate and fix ALL CI check failures on a PR, including pre-existing failures inherited from main.
- Never dismiss a CI failure as "pre-existing" or "unrelated to our changes" — if it fails on the PR, fix it in the PR.
- This includes: build, lint, test, docs build, security scanning (CodeQL), code review gates (CodeRabbit), workflow guard checks, and any other CI jobs.
- When a failure is caused by infrastructure outside the branch (e.g., rate limits, external service outages), implement or improve automated retry/bypass mechanisms in CI workflows.
- After fixing CI failures, verify locally where possible (build, vet, tests) before pushing.

## Phenotype Git and Delivery Workflow Protocol <!-- PHENOTYPE_GIT_DELIVERY_PROTOCOL -->

- Use branch-based delivery with pull requests; do not rely on direct default-branch writes where rulesets apply.
- Prefer stacked PRs for multi-part changes so each PR is small, reviewable, and independently mergeable.
- Keep PRs linear and scoped: one concern per PR, explicit dependency order for stacks, and clear migration steps.
- Enforce CI and required checks strictly: do not merge until all required checks and policy gates are green.
- Resolve all review threads and substantive PR comments before merge; do not leave unresolved reviewer feedback.
- Follow repository coding standards and best practices (typing, tests, lint, docs, security) before requesting merge.
- Rebase or restack to keep branches current with target branch and to avoid stale/conflicting stacks.
- When a ruleset or merge policy blocks progress, surface the blocker explicitly and adapt the plan (for example: open PR path, restack, or split changes).

## Phenotype Org Cross-Project Reuse Protocol <!-- PHENOTYPE_SHARED_REUSE_PROTOCOL -->

- Treat this repository as part of the broader Phenotype organization project collection, not an isolated codebase.
- During research and implementation, actively identify code that is sharable, modularizable, splittable, or decomposable for reuse across repositories.
- When reusable logic is found, prefer extraction into existing shared modules/projects first; if none fit, propose creating a new shared module/project.
- Include a `Cross-Project Reuse Opportunities` section in plans with candidate code, target shared location, impacted repos, and migration order.
- For cross-repo moves or ownership-impacting extractions, ask the user for confirmation on destination and rollout, then bake that into the execution plan.
- Execute forward-only migrations: extract shared code, update all callers, and remove duplicated local implementations.

## Phenotype Long-Term Stability and Non-Destructive Change Protocol <!-- PHENOTYPE_LONGTERM_STABILITY_PROTOCOL -->

- Optimize for long-term platform value over short-term convenience; choose durable solutions even when implementation complexity is higher.
- Classify proposed changes as `quick_fix` or `stable_solution`; prefer `stable_solution` unless an incident response explicitly requires a temporary fix.
- Do not use deletions/reversions as the default strategy; prefer targeted edits, forward fixes, and incremental hardening.
- Prefer moving obsolete or superseded material into `.archive/` over destructive removal when retention is operationally useful.
- Prefer clean manual merges, explicit conflict resolution, and auditable history over forceful rewrites, force merges, or history-destructive workflows.
- Prefer completing unused stubs into production-quality implementations when they represent intended product direction; avoid leaving stubs ignored indefinitely.
- Do not merge any PR while any check is failing, including non-required checks, unless the user gives explicit exception approval.
- When proposing a quick fix, include a scheduled follow-up path to a stable solution in the same plan.

## Worktree Discipline

- Feature work goes in `.worktrees/<topic>/`
- Legacy `PROJECT-wtrees/` and `repo-wtrees/` roots are for migration only and must not receive new work.
- Canonical repository remains on `main` for final integration and verification.
