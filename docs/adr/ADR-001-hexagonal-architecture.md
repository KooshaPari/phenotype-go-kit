# ADR-001: Hexagonal Architecture with Ports & Adapters

**Status**: Accepted
**Date**: 2026-03-25
**Deciders**: Phenotype Architecture Team

---

## Context

The phenotype-go-kit repository needs a scalable architecture that supports:
- Multiple transport adapters (REST, gRPC, CLI)
- Multiple storage backends (PostgreSQL, MongoDB, Redis)
- Plugin-based extensibility
- Testability and maintainability

Current state has mixed concerns with infrastructure code co-located with domain logic.

---

## Decision

We will adopt **Hexagonal Architecture** (also known as Ports & Adapters) with the following principles:

### 1. Dependency Rule
Outer layers depend on inner layers. Never the reverse.

```
Adapters → Ports → Domain
```

### 2. Ports (Interfaces)

**Inbound Ports (Driving)**:
- `UseCase` - Standard use case interface
- `CommandHandler` - CQRS command processing
- `QueryHandler` - CQRS query processing
- `EventHandler` - Domain event handling

**Outbound Ports (Driven)**:
- `Repository` - Data persistence
- `Cache` - Caching operations
- `EventPublisher` - Event distribution
- `ExternalService` - HTTP calls
- `SecretStore` - Secret management

### 3. Directory Structure

```
phenotype-go-kit/
├── contracts/              # Ports (interfaces)
│   ├── ports/
│   │   ├── inbound/      # Driving ports
│   │   └── outbound/     # Driven ports
│   ├── models/           # Domain models
│   └── plugins/          # Plugin contracts
├── internal/              # Private code
│   ├── domain/           # Domain entities, services
│   ├── application/       # Use cases, commands, queries
│   └── adapters/         # Concrete implementations
│       ├── primary/      # REST, gRPC handlers
│       └── secondary/    # DB, Cache, External
└── pkg/                  # Public libraries
```

### 4. Plugin System

Dynamic loading of plugins following the `Plugin` interface:
- Registry pattern for plugin management
- Manifest-based configuration
- Health monitoring and lifecycle management

---

## Consequences

### Positive
- Clear separation of concerns
- Easy to test with mock adapters
- Multiple transport options
- Plugin extensibility
- Independent development of layers

### Negative
- Initial complexity overhead
- More interfaces to maintain
- Learning curve for team

### Risks
- Over-engineering for simple features → Mitigate with YAGNI
- Interface explosion → Mitigate with GRASP/ISP

---

## Alternatives Considered

### 1. Clean Architecture (Microsoft)
Similar to hexagonal but with explicit use cases layer. Chosen hexagonal for simpler port concept.

### 2. Onion Architecture
Similar layering but no explicit adapters concept. Chosen hexagonal for clearer adapter naming.

### 3. Layered Architecture
Simpler but leads to dependency violations. Not chosen for scalability.

### 4. Microservices
Overkill for library. Not chosen.

---

## Implementation

See [contracts/README.md](../../contracts/README.md) for detailed implementation.

---

## References

- [Ports and Adapters - Alistair Cockburn](https://alistair.cockburn.us/strategic-use-of-package-structure/)
- [Hexagonal Architecture - Siemens](https://hexagonalarchitecture.org/)
- [Clean Architecture - Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
