// Package ports defines interfaces (ports) for hexagonal architecture.
// These interfaces are defined by the domain/application layer and
// implemented by the infrastructure layer (adapters).
//
// # Hexagonal Architecture
//
// The core domain is isolated from infrastructure concerns:
//   - Ports (interfaces) define how the domain interacts with the outside world
//   - Adapters implement these ports and handle external concerns (DB, cache, etc.)
//
// This follows the Dependency Inversion Principle (D from SOLID):
// high-level modules depend on abstractions, not concretions.
package ports

import (
	"context"
)

// CachePort defines the interface for caching implementations.
// This allows easy swapping between Redis, Memcached, in-memory, etc.
type CachePort interface {
	// Get retrieves a value by key.
	// Returns ErrKeyNotFound if the key doesn't exist.
	Get(ctx context.Context, key string) (string, error)

	// Set stores a value with TTL.
	// TTL must be positive; returns ErrInvalidTTL otherwise.
	Set(ctx context.Context, key string, value string, ttl any) error

	// Delete removes a key.
	Delete(ctx context.Context, key string) error

	// Exists checks if a key exists.
	Exists(ctx context.Context, key string) (bool, error)

	// Ping checks connectivity.
	Ping(ctx context.Context) error

	// Close releases resources.
	Close() error
}

// KeyNotFoundError signals a cache miss.
type ErrKeyNotFound struct{}

func (e ErrKeyNotFound) Error() string { return "key not found" }

// InvalidTTLError signals an invalid TTL value.
type ErrInvalidTTL struct{}

func (e ErrInvalidTTL) Error() string { return "invalid TTL" }
