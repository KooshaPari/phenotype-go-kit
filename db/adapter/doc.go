// Package adapter provides hexagonal architecture adapters for database operations.
// Adapters implement the outbound ports defined in contracts/ports/outbound.
//
// Architecture: Hexagonal (Ports & Adapters)
//
// This package contains concrete implementations that connect the domain
// to specific database technologies (PostgreSQL, MySQL, SQLite, etc.).
//
// # Adapter Pattern
//
// Adapters wrap external dependencies and implement domain-defined interfaces:
//
//	type PostgresAdapter struct {
//	    db *sql.DB
//	}
//
//	func (a *PostgresAdapter) Query(ctx context.Context, query string, args ...any) (outbound.Rows, error) {
//	    return a.db.QueryContext(ctx, query, args...)
//	}
//
// # Dependency Rule
//
// Adapters depend on domain ports, never the other way around:
//
//	contracts/ports/outbound/  <-  adapter/  <-  infrastructure
//	       (interface)              (impl)
//
// This ensures:
//   - Domain has no dependencies on infrastructure
//   - Easy testing via mock adapters
//   - Technology substitution without domain changes
package adapter
