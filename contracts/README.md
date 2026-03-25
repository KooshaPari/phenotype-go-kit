# Architecture: Contracts (Ports & Adapters)

This directory contains the hexagonal architecture contracts defining ports and domain models.

## Directory Structure

```
contracts/
├── doc.go                    # Architecture overview
├── ports/                    # Port interfaces
│   ├── doc.go               # Ports documentation
│   ├── inbound/             # Driving (inbound) ports
│   │   ├── doc.go
│   │   └── ports.go         # UseCase, CommandHandler, QueryHandler, etc.
│   └── outbound/            # Driven (outbound) ports
│       ├── doc.go
│       └── ports.go         # Repository, Cache, EventBus, etc.
├── models/                  # Domain models and DTOs
│   ├── doc.go
│   └── events.go           # DomainEvent, CommandResult, etc.
└── plugins/                 # Plugin system contracts
    ├── doc.go
    ├── plugin.go            # Plugin interface, Metadata, Registry
    └── manifest.go          # Plugin manifest schema
```

## Hexagonal Architecture

```
                        ┌───────────────────────────────────────┐
                        │           Driving Adapters              │
                        │  (REST, gRPC, CLI, Message Handlers)   │
                        └─────────────────────┬─────────────────┘
                                              │
                        ┌─────────────────────▼─────────────────┐
                        │         Inbound Ports (Driving)        │
                        │  UseCase, CommandHandler, QueryHandler  │
                        └─────────────────────┬─────────────────┘
                                              │
┌─────────────────────────────────────────────┼─────────────────────────────┐
│                                             │                              │
│                   ┌────────────────────────▼────────────────────────┐     │
│                   │               Domain Core                          │     │
│                   │  Entities, Value Objects, Domain Services, Events   │     │
│                   └────────────────────────▲────────────────────────┘     │
│                                             │                              │
└─────────────────────────────────────────────┼─────────────────────────────┘
                                              │
                        ┌─────────────────────▼─────────────────┐
                        │         Outbound Ports (Driven)         │
                        │  Repository, Cache, EventBus, External  │
                        └─────────────────────┬─────────────────┘
                                              │
                        ┌─────────────────────▼─────────────────┐
                        │           Driven Adapters              │
                        │  (Postgres, Redis, Kafka, HTTP Client)  │
                        └───────────────────────────────────────┘
```

## Ports

### Inbound Ports (Driving)

- **UseCase**: Standard use case interface
- **CommandHandler**: Handles commands (CQRS)
- **QueryHandler**: Handles queries (CQRS)
- **EventHandler**: Handles domain events
- **InputPort**: Generic input port interface
- **Validator**: Input validation
- **Interceptor**: Request/response middleware

### Outbound Ports (Driven)

- **Repository**: Data persistence abstraction
- **QueryRepository**: Read-only data access (CQRS)
- **EventStore**: Event persistence (Event Sourcing)
- **EventPublisher**: Event distribution
- **Cache**: Caching abstraction
- **ExternalService**: HTTP client abstraction
- **SecretStore**: Secret management
- **MetricsCollector**: Metrics collection
- **Logger**: Logging abstraction
- **ConfigProvider**: Configuration management

## Plugin System

Plugins are external adapters that extend functionality through dynamic loading.

### Plugin Lifecycle

1. **Register**: Register plugin factory with registry
2. **Load**: Load plugin from source
3. **Init**: Initialize with configuration
4. **Start**: Start plugin operation
5. **Stop**: Graceful shutdown

## Usage Example

```go
import (
    "github.com/Phenotype/phenotype-go-kit/contracts/ports/inbound"
    "github.com/Phenotype/phenotype-go-kit/contracts/ports/outbound"
)

// Define your use case
type CreateOrderUseCase struct {
    orderRepo outbound.Repository[Order, string]
}

func (uc *CreateOrderUseCase) Execute(ctx context.Context, input interface{}) (interface{}, error) {
    // Implementation
}

// Use in adapter
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
    useCase := h.registry.Get(inbound.UseCaseType, "CreateOrder")
    result, err := useCase.Execute(r.Context(), input)
    // Handle response
}
```

## Related

- [ADR-001: Hexagonal Architecture Decision](../docs/adr/ADR-001-hexagonal-architecture.md)
- [ADR-002: Plugin System Design](../docs/adr/ADR-002-plugin-system.md)
