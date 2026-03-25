# ADR-007: Database & Secrets Hexagonal Architecture

**Status**: Accepted
**Date**: 2026-03-25
**Author**: Forge

## Context

The `db/` and `secrets/` packages in `phenotype-go-kit` were implemented with direct dependencies on external libraries, making them difficult to test, extend, and swap implementations.

## Decision

We will apply Hexagonal Architecture (Ports & Adapters) to both packages:

### Database Package Structure

```
db/
├── adapter/           # Infrastructure adapters
│   ├── doc.go
│   ├── postgres.go   # PostgreSQL implementation
│   ├── mysql.go      # MySQL implementation
│   └── sqlite.go     # SQLite implementation
├── indexes.go         # Domain index definitions
├── pool.go           # Connection pool utilities
└── query.go          # Query builder utilities
```

**Outbound Ports** (`contracts/ports/outbound/db.go`):
- `QueryExecutor` - Core query operations
- `Transaction` - Transaction management
- `ConnectionPool` - Pool configuration
- `IndexManager` - Index operations
- `MigrationExecutor` - Database migrations

### Secrets Package Structure

```
secrets/
├── adapter/           # Infrastructure adapters
│   ├── doc.go
│   ├── vault.go      # HashiCorp Vault implementation
│   ├── aws.go        # AWS Secrets Manager implementation
│   └── env.go        # Environment variables (dev)
└── manager.go        # Manager utilities
```

**Outbound Ports** (`contracts/ports/outbound/secrets.go`):
- `SecretPort` - CRUD operations
- `SecretReader` - Read-only access
- `SecretWriter` - Write-only access
- `SecretEncryptor` - Encryption/decryption
- `SecretRotator` - Secret rotation

## Consequences

### Positive
- **Testability**: Easy to mock adapters for unit tests
- **Flexibility**: Swap database/secrets providers without changing domain code
- **SOLID Compliance**: ISP (segregated interfaces), DIP (domain depends on abstractions)
- **GRASP Patterns**: Low Coupling, High Cohesion

### Negative
- Additional interface layer adds indirection
- More files to maintain

## References
- ADR-001: Hexagonal Architecture
- ADR-006: Design Principles
