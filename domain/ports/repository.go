// Package ports defines interfaces for repository pattern.
// Following DDD and Clean Architecture principles.
package ports

import (
	"context"
)

// RepositoryPort defines the interface for generic CRUD operations.
// This follows the Repository pattern from DDD.
//
// # Usage
//
// Define repository interfaces in the domain layer (ports),
// then implement in infrastructure layer (adapters).
//
// Example:
//
//	type UserRepository interface {
//	    RepositoryPort[User, string]
//	    FindByEmail(ctx context.Context, email string) (*User, error)
//	}
type RepositoryPort[Entity any, ID comparable] interface {
	// Save creates or updates an entity.
	Save(ctx context.Context, entity Entity) (Entity, error)

	// FindByID retrieves an entity by ID.
	// Returns (nil, ErrNotFound) if not found.
	FindByID(ctx context.Context, id ID) (Entity, error)

	// Delete removes an entity by ID.
	Delete(ctx context.Context, id ID) error

	// List returns all entities (use with caution for large datasets).
	List(ctx context.Context) ([]Entity, error)
}

// QueryRepositoryPort extends RepositoryPort with query capabilities.
type QueryRepositoryPort[Entity any, ID comparable] interface {
	RepositoryPort[Entity, ID]

	// FindByFilter returns entities matching the given filter.
	FindByFilter(ctx context.Context, filter map[string]any) ([]Entity, error)

	// Count returns the total number of entities.
	Count(ctx context.Context) (int64, error)
}

// NotFoundError signals an entity was not found.
type ErrNotFound struct {
	Entity string
	ID     any
}

func (e *ErrNotFound) Error() string {
	return e.Entity + " not found: " + toString(e.ID)
}

func toString(v any) string {
	switch val := v.(type) {
	case string:
		return val
	default:
		return "<unknown>"
	}
}
