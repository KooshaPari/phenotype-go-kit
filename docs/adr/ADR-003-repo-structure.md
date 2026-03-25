# ADR-003: Top-Level Directory Structure & File Organization

**Status**: Accepted
**Date**: 2026-03-25
**Deciders**: Phenotype Architecture Team

---

## Context

Phenotype repositories need consistent, maintainable top-level organization that supports:
- Clear project boundaries
- Easy navigation
- CI/CD integration
- Documentation discovery
- Polyrepo management

Current repos have inconsistent structures causing confusion.

---

## Decision

### Standard Top-Level Structure

```
phenotype-{service}/
├── .github/                    # GitHub configs ( workflows, actions)
│   ├── workflows/            # CI/CD pipelines
│   ├── actions/              # Custom actions
│   └── ISSUE_TEMPLATE/       # Issue templates
│
├── .vscode/                   # Editor configs
│   ├── settings.json
│   └── extensions.json
│
├── cmd/                       # Command-line entry points
│   └── {binary}/             # One dir per binary
│       └── main.go
│
├── internal/                  # Private application code (not importable)
│   ├── domain/               # Domain entities, value objects, services
│   ├── application/          # Use cases, commands, queries
│   └── adapters/            # Implementations
│       ├── primary/          # Driving (REST, gRPC, CLI)
│       └── secondary/        # Driven (DB, Cache, External)
│
├── pkg/                       # Public libraries (importable)
│   ├── {library}/            # One pkg per concern
│   └── ...
│
├── contracts/                 # Interface definitions
│   ├── ports/
│   ├── models/
│   └── plugins/
│
├── api/                       # API definitions
│   ├── openapi/
│   ├── proto/
│   └── graphql/
│
├── configs/                   # Configuration files
│   ├── config.yaml
│   └── config.schema.json
│
├── scripts/                   # Build, deployment scripts
│   ├── build.sh
│   └── deploy.sh
│
├── docs/                      # Documentation
│   ├── adr/                  # Architecture Decision Records
│   ├── guides/
│   └── runbooks/
│
├── tests/                     # External test data and helpers
│   ├── fixtures/
│   └── e2e/
│
├── migrations/                # Database migrations
│
├── examples/                  # Usage examples
│
├── CONTRIBUTING.md
├── CHANGELOG.md
├── LICENSE
├── README.md
├── Makefile
├── Taskfile.yml
├── go.mod / package.json     # Language-specific
└── .gitignore
```

### Naming Conventions

| Type | Convention | Example |
|------|------------|---------|
| Directories | kebab-case | `phenotype-config`, `cli-proxy` |
| Packages | lowercase, no separators | `domain`, `adapters` |
| Binaries | kebab-case | `phenotype-agent`, `helios-cli` |
| Files | kebab-case | `api-client.go`, `user-service.go` |
| Go modules | kebab-case, prefixed | `github.com/phenotype/agent-api` |
| Docker images | lowercase with dashes | `phenotype/agent:v1` |

### File Organization Principles

#### 1. Package-by-Feature (Preferred)
```
internal/
├── user/
│   ├── user.go           # Entity
│   ├── user_service.go   # Domain service
│   ├── user_repo.go      # Repository interface
│   ├── commands.go       # CQRS commands
│   ├── queries.go        # CQRS queries
│   └── user_test.go
├── order/
│   └── ...
```

#### 2. Layered (Alternative for Simple Projects)
```
internal/
├── models/           # All entities
├── services/         # All domain services
├── handlers/         # All adapters
└── repository/       # All persistence
```

#### 3. Group by Kind (for utilities)
```
pkg/
├── validator/
│   ├── validator.go
│   ├── email.go
│   ├── url.go
│   └── validator_test.go
├── logger/
│   └── ...
```

---

## Special Directory Rules

### `.github/`
- Contains all GitHub-related configs
- `.github/workflows/` for CI/CD
- `.github/actions/` for reusable actions
- `.github/ISSUE_TEMPLATE/` for issue templates

### `cmd/`
- One directory per application binary
- Minimal main.go (delegate to app initialization)
- No shared code here

### `internal/`
- Cannot be imported by external packages
- Sub-package structure follows domain boundaries
- Use `internal/` for everything private

### `pkg/`
- Public, versioned APIs
- Documentation required
- Semantic versioning
- Examples in `/examples/`

### `api/`
- OpenAPI/Swagger specs
- Protocol Buffer definitions
- GraphQL schemas
- Breaking changes require version bump

### `docs/`
```
docs/
├── adr/                    # Architecture Decision Records
├── guides/                 # How-to guides
├── tutorials/              # Step-by-step tutorials
├── reference/              # API reference
├── concepts/               # Architecture explanations
└── runbooks/               # Operational procedures
```

### `tests/`
- External test data only
- Fixtures, mocks, test utilities
- e2e tests with full app initialization
- Unit tests colocated with code

---

## Consequences

### Positive
- Predictable navigation
- Clear separation of concerns
- Easy onboarding
- Standardized automation

### Negative
- Learning curve for newcomers
- May need adjustment for legacy projects

### Risks
- Over-structuring simple projects → Apply YAGNI, start simple
- Inconsistent migration → Use automated tools

---

## Alternative Structures Considered

### Google Style (OSS)
```
go/
├── doc.go
├── doc_test.go
├── impl.go
└── impl_test.go
```
- Good for single-package libraries
- Not scalable for large apps

### Ruby on Rails Style
```
app/
controllers/
models/
views/
helpers/
```
- Convention over configuration
- Good for web apps
- Tends toward anemic domain

### Node.js Style
```
src/
  index.js
  app.js
  config/
  routes/
```
- Simple
- Tends toward spaghetti

---

## Migration Guide

### From Legacy Structure
```bash
# 1. Create new structure
mkdir -p cmd internal pkg api docs

# 2. Move code by type
mv src/domain internal/
mv src/handlers internal/adapters/primary/
mv src/repositories internal/adapters/secondary/

# 3. Update import paths
# Use IDE refactoring or sed

# 4. Update CI/CD
# Update paths in workflows

# 5. Verify
# Run tests, linters, formatters
```

---

## References

- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Package by Feature](https://phauer.com/2020/package-by-feature/)
- [Standard Go Project Layout](medium.com/@benjkjames)
