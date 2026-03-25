package inbound

import (
	"context"
	"time"
)

// CacheCommands defines CQRS commands for cache operations.
type CacheCommands struct{}

// NewCacheCommands creates new cache command handlers.
func NewCacheCommands() *CacheCommands {
	return &CacheCommands{}
}

// SetCacheCommand represents a command to set a cache value.
type SetCacheCommand struct {
	Key   string
	Value string
	TTL   time.Duration
}

// SetCacheHandler handles SetCacheCommand.
type SetCacheHandler func(ctx context.Context, cmd SetCacheCommand) error

// DeleteCacheCommand represents a command to delete a cache key.
type DeleteCacheCommand struct {
	Key string
}

// DeleteCacheHandler handles DeleteCacheCommand.
type DeleteCacheHandler func(ctx context.Context, cmd DeleteCacheCommand) error

// ExpireCacheCommand represents a command to set expiration on a key.
type ExpireCacheCommand struct {
	Key string
	TTL time.Duration
}

// ExpireCacheHandler handles ExpireCacheCommand.
type ExpireCacheHandler func(ctx context.Context, cmd ExpireCacheCommand) error

// InvalidateByPatternCommand represents a command to invalidate keys by pattern.
type InvalidateByPatternCommand struct {
	Pattern string
}

// InvalidateByPatternHandler handles InvalidateByPatternCommand.
type InvalidateByPatternHandler func(ctx context.Context, cmd InvalidateByPatternCommand) (int64, error)

// InvalidateByTagsCommand represents a command to invalidate keys by tags.
type InvalidateByTagsCommand struct {
	Tags []string
}

// InvalidateByTagsHandler handles InvalidateByTagsCommand.
type InvalidateByTagsHandler func(ctx context.Context, cmd InvalidateByTagsCommand) error
