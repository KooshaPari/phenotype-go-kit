package inbound

import "context"

// CacheQueries defines CQRS queries for cache operations.
type CacheQueries struct{}

// NewCacheQueries creates new cache query handlers.
func NewCacheQueries() *CacheQueries {
	return &CacheQueries{}
}

// GetCacheQuery represents a query to get a cache value.
type GetCacheQuery struct {
	Key string
}

// GetCacheHandler handles GetCacheQuery.
type GetCacheHandler func(ctx context.Context, query GetCacheQuery) (string, error)

// ExistsCacheQuery represents a query to check if a key exists.
type ExistsCacheQuery struct {
	Key string
}

// ExistsCacheHandler handles ExistsCacheQuery.
type ExistsCacheHandler func(ctx context.Context, query ExistsCacheQuery) (bool, error)

// TTLCacheQuery represents a query to get remaining TTL.
type TTLCacheQuery struct {
	Key string
}

// TTLCacheHandler handles TTLCacheQuery.
type TTLHandler func(ctx context.Context, query TTLCacheQuery) (TTLResult, error)

// TTLResult contains the TTL query result.
type TTLResult struct {
	Duration int64 // seconds
	Exists   bool
	HasTTL   bool // false if key doesn't exist or has no expiry
}

// GetTagsQuery represents a query to get tags for a cached key.
type GetTagsQuery struct {
	Key string
}

// GetTagsHandler handles GetTagsQuery.
type GetTagsHandler func(ctx context.Context, query GetTagsQuery) ([]string, error)

// SortedSetRangeQuery represents a query to get sorted set members.
type SortedSetRangeQuery struct {
	Key   string
	Start int64
	Stop  int64
}

// SortedSetRangeHandler handles SortedSetRangeQuery.
type SortedSetRangeHandler func(ctx context.Context, query SortedSetRangeQuery) ([]string, error)
