package outbound

import (
	"context"
	"errors"
	"time"
)

// Standard errors for outbound ports.
var (
	ErrKeyNotFound = errors.New("key not found")
	ErrInvalidTTL  = errors.New("invalid TTL")
)

// CachePort defines the interface for cache implementations.
// This is an outbound (driven) port used by the application/domain layers.
//
// # Design Patterns Applied
//
//   - Interface Segregation (ISP): Focused, minimal interface
//   - Dependency Inversion (DIP): Domain depends on abstraction
//   - Low Coupling: Minimal methods, focused responsibility
//
// # Implementations
//
//   - Redis: github.com/KooshaPari/phenotype-go-kit/cache/adapter (RedisCacheAdapter)
type CachePort interface {
	// Get retrieves a value by key.
	// Returns ErrKeyNotFound if key does not exist.
	Get(ctx context.Context, key string) (string, error)

	// Set stores a value with TTL.
	// Returns ErrInvalidTTL if ttl <= 0.
	Set(ctx context.Context, key string, value string, ttl time.Duration) error

	// SetNX stores a value only if key doesn't exist.
	SetNX(ctx context.Context, key string, value string, ttl time.Duration) (bool, error)

	// Delete removes a key.
	Delete(ctx context.Context, key string) error

	// Exists checks if a key exists.
	Exists(ctx context.Context, key string) (bool, error)

	// Expire sets expiration on a key.
	Expire(ctx context.Context, key string, ttl time.Duration) error

	// TTL returns remaining TTL for a key.
	TTL(ctx context.Context, key string) (time.Duration, error)

	// Ping checks connectivity.
	Ping(ctx context.Context) error

	// Close closes the connection.
	Close() error
}

// CacheJSONPort extends CachePort with JSON serialization support.
type CacheJSONPort interface {
	CachePort

	// GetJSON retrieves and unmarshals JSON.
	GetJSON(ctx context.Context, key string, dest interface{}) error

	// SetJSON marshals and stores JSON.
	SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error
}

// CacheCounterPort extends CachePort with atomic counter operations.
type CacheCounterPort interface {
	CachePort

	// Incr increments a counter.
	Incr(ctx context.Context, key string) (int64, error)

	// Decr decrements a counter.
	Decr(ctx context.Context, key string) (int64, error)
}

// CacheSortedSetPort extends CachePort with sorted set operations.
type CacheSortedSetPort interface {
	CachePort

	// SortedSetAdd adds to a sorted set.
	SortedSetAdd(ctx context.Context, key string, member string, score float64) error

	// SortedSetRange retrieves sorted set members by rank.
	SortedSetRange(ctx context.Context, key string, start, stop int64) ([]string, error)

	// SortedSetRemove removes from a sorted set.
	SortedSetRemove(ctx context.Context, key string, members ...string) error
}

// CacheInvalidationPort extends CachePort with cache invalidation strategies.
type CacheInvalidationPort interface {
	CachePort

	// DeletePattern deletes all keys matching a pattern.
	DeletePattern(ctx context.Context, pattern string) (int64, error)
}

// InvalidationStrategy defines cache invalidation behavior.
type InvalidationStrategy interface {
	Invalidate(ctx context.Context, cache CachePort, keys ...string) error
}
