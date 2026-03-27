package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/KooshaPari/phenotype-go-kit/contracts/ports/inbound"
	"github.com/KooshaPari/phenotype-go-kit/contracts/ports/outbound"
)

// Errors for cache service.
var (
	ErrKeyNotFound = errors.New("key not found")
	ErrInvalidTTL  = errors.New("invalid TTL")
)

// CacheService handles cache operations following CQRS pattern.
type CacheService struct {
	cache  outbound.CacheJSONPort
	logger *slog.Logger
}

// NewCacheService creates a new cache service.
func NewCacheService(cache outbound.CacheJSONPort) *CacheService {
	return &CacheService{
		cache:  cache,
		logger: slog.Default(),
	}
}

// SetCacheHandler handles SetCacheCommand.
func (s *CacheService) SetCacheHandler() inbound.SetCacheHandler {
	return func(ctx context.Context, cmd inbound.SetCacheCommand) error {
		if cmd.TTL <= 0 {
			return ErrInvalidTTL
		}
		return s.cache.Set(ctx, cmd.Key, cmd.Value, cmd.TTL)
	}
}

// SetCacheJSONHandler handles setting a JSON value.
func (s *CacheService) SetCacheJSONHandler() func(ctx context.Context, cmd inbound.SetCacheCommand) error {
	return func(ctx context.Context, cmd inbound.SetCacheCommand) error {
		if cmd.TTL <= 0 {
			return ErrInvalidTTL
		}
		return s.cache.SetJSON(ctx, cmd.Key, cmd.Value, cmd.TTL)
	}
}

// DeleteCacheHandler handles DeleteCacheCommand.
func (s *CacheService) DeleteCacheHandler() inbound.DeleteCacheHandler {
	return func(ctx context.Context, cmd inbound.DeleteCacheCommand) error {
		return s.cache.Delete(ctx, cmd.Key)
	}
}

// ExpireCacheHandler handles ExpireCacheCommand.
func (s *CacheService) ExpireCacheHandler() inbound.ExpireCacheHandler {
	return func(ctx context.Context, cmd inbound.ExpireCacheCommand) error {
		return s.cache.Expire(ctx, cmd.Key, cmd.TTL)
	}
}

// InvalidateByPatternHandler handles InvalidateByPatternCommand.
func (s *CacheService) InvalidateByPatternHandler() inbound.InvalidateByPatternHandler {
	return func(ctx context.Context, cmd inbound.InvalidateByPatternCommand) (int64, error) {
		invalidator := &PatternInvalidator{}
		return invalidator.Invalidate(ctx, s.cache, cmd.Pattern)
	}
}

// PatternInvalidator implements outbound.InvalidationStrategy.
type PatternInvalidator struct{}

func (p *PatternInvalidator) Invalidate(ctx context.Context, cache outbound.CachePort, keys ...string) (int64, error) {
	var total int64
	for _, key := range keys {
		pattern := key + "*"
		// Use DeletePattern if available
		if inv, ok := cache.(interface {
			DeletePattern(ctx context.Context, pattern string) (int64, error)
		}); ok {
			n, err := inv.DeletePattern(ctx, pattern)
			if err != nil {
				return total, err
			}
			total += n
		}
	}
	return total, nil
}

// GetCacheHandler handles GetCacheQuery.
func (s *CacheService) GetCacheHandler() inbound.GetCacheHandler {
	return func(ctx context.Context, query inbound.GetCacheQuery) (string, error) {
		val, err := s.cache.Get(ctx, query.Key)
		if errors.Is(err, outbound.ErrKeyNotFound) {
			return "", ErrKeyNotFound
		}
		return val, err
	}
}

// GetCacheJSONHandler handles getting a JSON value.
func (s *CacheService) GetCacheJSONHandler(dest interface{}) inbound.GetCacheHandler {
	return func(ctx context.Context, query inbound.GetCacheQuery) (string, error) {
		err := s.cache.GetJSON(ctx, query.Key, dest)
		if errors.Is(err, outbound.ErrKeyNotFound) {
			return "", ErrKeyNotFound
		}
		return "", err
	}
}

// ExistsCacheHandler handles ExistsCacheQuery.
func (s *CacheService) ExistsCacheHandler() inbound.ExistsCacheHandler {
	return func(ctx context.Context, query inbound.ExistsCacheQuery) (bool, error) {
		return s.cache.Exists(ctx, query.Key)
	}
}

// TTLCacheHandler handles TTLCacheQuery.
func (s *CacheService) TTLCacheHandler() inbound.TTLHandler {
	return func(ctx context.Context, query inbound.TTLCacheQuery) (inbound.TTLResult, error) {
		ttl, err := s.cache.TTL(ctx, query.Key)
		if err != nil {
			return inbound.TTLResult{}, err
		}

		result := inbound.TTLResult{
			Exists:   ttl > 0,
			HasTTL:   ttl > 0,
			Duration: int64(ttl.Seconds()),
		}

		// Check if key exists
		exists, _ := s.cache.Exists(ctx, query.Key)
		result.Exists = exists

		return result, nil
	}
}

// CacheWarmer preloads frequently accessed data into cache.
type CacheWarmer struct {
	cache   outbound.CachePort
	fetcher func(ctx context.Context, key string) (string, time.Duration, error)
	logger  *slog.Logger
}

// NewCacheWarmer creates a new cache warmer.
func NewCacheWarmer(cache outbound.CachePort, fetcher func(ctx context.Context, key string) (string, time.Duration, error)) *CacheWarmer {
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
