package cache

import (
	"context"
	"log/slog"
	"strings"
	"time"
)

// InvalidationStrategy defines cache invalidation behavior.
type InvalidationStrategy interface {
	Invalidate(ctx context.Context, cache *Cache, keys ...string) error
}

// InvalidateAll removes all matching keys.
type InvalidateAll struct {
	prefix string
}

// NewInvalidateAll creates a strategy that deletes keys matching a prefix.
func NewInvalidateAll(prefix string) *InvalidateAll {
	return &InvalidateAll{prefix: prefix}
}

func (i *InvalidateAll) Invalidate(ctx context.Context, cache *Cache, keys ...string) error {
	for _, key := range keys {
		pattern := key + "*"
		n, err := cache.DeletePattern(ctx, pattern)
		if err != nil {
			return err
		}
		slog.Debug("cache invalidated", "pattern", pattern, "count", n)
	}
	return nil
}

// InvalidateByTags invalidates keys by tags.
type InvalidateByTags struct {
	tagKey string
}

// NewInvalidateByTags creates a tag-based invalidation strategy.
func NewInvalidateByTags() *InvalidateByTags {
	return &InvalidateByTags{tagKey: "tags"}
}

func (i *InvalidateByTags) Invalidate(ctx context.Context, cache *Cache, tags ...string) error {
	for _, tag := range tags {
		pattern := i.tagKey + ":" + tag + ":*"
		keys, err := cache.client.Keys(ctx, "phenotype:"+pattern).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			_, err = cache.client.Del(ctx, keys...).Result()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// CacheWarmer preloads frequently accessed data into cache.
type CacheWarmer struct {
	cache   *Cache
	fetcher func(ctx context.Context, key string) (string, time.Duration, error)
	logger  *slog.Logger
}

// NewCacheWarmer creates a new cache warmer.
func NewCacheWarmer(cache *Cache, fetcher func(ctx context.Context, key string) (string, time.Duration, error)) *CacheWarmer {
	return &CacheWarmer{
		cache:   cache,
		fetcher: fetcher,
		logger:  slog.Default(),
	}
}

// WarmKey loads a single key into cache.
func (w *CacheWarmer) WarmKey(ctx context.Context, key string) error {
	value, ttl, err := w.fetcher(ctx, key)
	if err != nil {
		w.logger.Error("cache warm failed", "key", key, "error", err)
		return err
	}

	if err := w.cache.Set(ctx, key, value, ttl); err != nil {
		w.logger.Error("cache set failed", "key", key, "error", err)
		return err
	}

	w.logger.Debug("cache warmed", "key", key, "ttl", ttl)
	return nil
}

// WarmKeys loads multiple keys into cache.
func (w *CacheWarmer) WarmKeys(ctx context.Context, keys []string) error {
	for _, key := range keys {
		if err := w.WarmKey(ctx, key); err != nil {
			w.logger.Warn("key warm skipped", "key", key, "error", err)
		}
	}
	return nil
}

// WarmByPattern loads keys matching a pattern.
func (w *CacheWarmer) WarmByPattern(ctx context.Context, pattern string, fetchKeys func(ctx context.Context) ([]string, error)) error {
	keys, err := fetchKeys(ctx)
	if err != nil {
		return err
	}

	// Filter to only missing keys
	missing := make([]string, 0)
	for _, key := range keys {
		exists, _ := w.cache.Exists(ctx, key)
		if !exists {
			missing = append(missing, key)
		}
	}

	if len(missing) == 0 {
		w.logger.Debug("all keys already cached", "pattern", pattern)
		return nil
	}

	w.logger.Info("warming cache", "pattern", pattern, "count", len(missing))
	return w.WarmKeys(ctx, missing)
}

// TagBasedCache adds tags to cached items for targeted invalidation.
type TagBasedCache struct {
	cache  *Cache
	tagKey string
}

// NewTagBasedCache creates a cache with tag support.
func NewTagBasedCache(cache *Cache) *TagBasedCache {
	return &TagBasedCache{
		cache:  cache,
		tagKey: "tags",
	}
}

// SetWithTags stores a value with associated tags.
func (t *TagBasedCache) SetWithTags(ctx context.Context, key string, value string, ttl time.Duration, tags ...string) error {
	pipe := t.cache.client.Pipeline()

	// Store the value
	pipe.Set(ctx, t.cache.prefix+key, value, ttl)

	// Store tag associations
	for _, tag := range tags {
		tagKey := t.tagKey + ":" + tag + ":" + key
		pipe.Set(ctx, t.cache.prefix+tagKey, "1", ttl)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// GetTags retrieves all tags for a key.
func (t *TagBasedCache) GetTags(ctx context.Context, key string) ([]string, error) {
	pattern := t.tagKey + ":*:" + key
	keys, err := t.cache.client.Keys(ctx, t.cache.prefix+pattern).Result()
	if err != nil {
		return nil, err
	}

	tags := make([]string, 0, len(keys))
	for _, k := range keys {
		// Extract tag from key format: tags:tagName:originalKey
		parts := strings.Split(k, ":")
		if len(parts) >= 2 {
			tags = append(tags, parts[1])
		}
	}
	return tags, nil
}

// InvalidateByTag removes all cached items with a specific tag.
func (t *TagBasedCache) InvalidateByTag(ctx context.Context, tag string) error {
	pattern := t.tagKey + ":" + tag + ":*"
	keys, err := t.cache.client.Keys(ctx, t.cache.prefix+pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		// Extract original keys from tag keys
		originalKeys := make([]string, 0)
		for _, k := range keys {
			parts := strings.Split(k, ":")
			if len(parts) >= 3 {
				originalKeys = append(originalKeys, parts[2])
			}
		}

		// Delete original keys and tag keys
		allKeys := append(originalKeys, keys...)
		_, err = t.cache.client.Del(ctx, allKeys...).Result()
		return err
	}
	return nil
}
