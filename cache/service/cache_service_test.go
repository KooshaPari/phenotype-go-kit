package service

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/KooshaPari/phenotype-go-kit/contracts/ports/inbound"
	"github.com/KooshaPari/phenotype-go-kit/contracts/ports/outbound"
)

// mockCache implements outbound.CacheJSONPort for testing.
// Following Law of Demeter - only exposes needed operations.
type mockCache struct {
	mu   sync.RWMutex
	data map[string]string
	ttls map[string]time.Duration
}

func newMockCache() *mockCache {
	return &mockCache{
		data: make(map[string]string),
		ttls: make(map[string]time.Duration),
	}
}

func (m *mockCache) Get(ctx context.Context, key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if _, ok := m.data[key]; !ok {
		return "", outbound.ErrKeyNotFound
	}
	return m.data[key], nil
}

func (m *mockCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if ttl <= 0 {
		return outbound.ErrInvalidTTL
	}
	m.data[key] = value
	m.ttls[key] = ttl
	return nil
}

func (m *mockCache) SetNX(ctx context.Context, key, value string, ttl time.Duration) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.data[key]; ok {
		return false, nil
	}
	m.data[key] = value
	m.ttls[key] = ttl
	return true, nil
}

func (m *mockCache) Delete(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.data, key)
	delete(m.ttls, key)
	return nil
}

func (m *mockCache) Exists(ctx context.Context, key string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, ok := m.data[key]
	return ok, nil
}

func (m *mockCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.data[key]; !ok {
		return outbound.ErrKeyNotFound
	}
	m.ttls[key] = ttl
	return nil
}

func (m *mockCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ttl, ok := m.ttls[key]
	if !ok {
		return 0, outbound.ErrKeyNotFound
	}
	return ttl, nil
}

func (m *mockCache) Ping(ctx context.Context) error {
	return nil
}

func (m *mockCache) Close() error {
	return nil
}

func (m *mockCache) GetJSON(ctx context.Context, key string, dest interface{}) error {
	val, err := m.Get(ctx, key)
	if err != nil {
		return err
	}
	// Simple string storage for testing
	if s, ok := dest.(*string); ok {
		*s = val
	}
	return nil
}

func (m *mockCache) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	// Simple string storage for testing
	if s, ok := value.(string); ok {
		return m.Set(ctx, key, s, ttl)
	}
	return m.Set(ctx, key, "json-value", ttl)
}

// Tests

func TestCacheService_SetCacheHandler(t *testing.T) {
	cache := newMockCache()
	service := NewCacheService(cache)
	handler := service.SetCacheHandler()

	ctx := context.Background()
	cmd := inbound.SetCacheCommand{
		Key:   "test-key",
		Value: "test-value",
		TTL:   time.Hour,
	}

	err := handler(ctx, cmd)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify
	val, err := cache.Get(ctx, "test-key")
	if err != nil {
		t.Fatalf("expected value, got error: %v", err)
	}
	if val != "test-value" {
		t.Errorf("expected 'test-value', got %q", val)
	}
}

func TestCacheService_SetCacheHandler_InvalidTTL(t *testing.T) {
	cache := newMockCache()
	service := NewCacheService(cache)
	handler := service.SetCacheHandler()

	ctx := context.Background()
	cmd := inbound.SetCacheCommand{
		Key:   "test-key",
		Value: "test-value",
		TTL:   0, // Invalid
	}

	err := handler(ctx, cmd)
	if !errors.Is(err, ErrInvalidTTL) {
		t.Errorf("expected ErrInvalidTTL, got %v", err)
	}
}

func TestCacheService_GetCacheHandler(t *testing.T) {
	cache := newMockCache()
	cache.data["existing-key"] = "existing-value"
	cache.ttls["existing-key"] = time.Hour

	service := NewCacheService(cache)
	handler := service.GetCacheHandler()

	ctx := context.Background()

	// Test existing key
	query := inbound.GetCacheQuery{Key: "existing-key"}
	val, err := handler(ctx, query)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if val != "existing-value" {
		t.Errorf("expected 'existing-value', got %q", val)
	}

	// Test missing key
	query = inbound.GetCacheQuery{Key: "missing-key"}
	_, err = handler(ctx, query)
	if !errors.Is(err, ErrKeyNotFound) {
		t.Errorf("expected ErrKeyNotFound, got %v", err)
	}
}

func TestCacheService_DeleteCacheHandler(t *testing.T) {
	cache := newMockCache()
	cache.data["delete-key"] = "delete-value"

	service := NewCacheService(cache)
	handler := service.DeleteCacheHandler()

	ctx := context.Background()
	cmd := inbound.DeleteCacheCommand{Key: "delete-key"}

	err := handler(ctx, cmd)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify deleted
	_, err = cache.Get(ctx, "delete-key")
	if !errors.Is(err, outbound.ErrKeyNotFound) {
		t.Errorf("expected ErrKeyNotFound after delete, got %v", err)
	}
}

func TestCacheService_ExpireCacheHandler(t *testing.T) {
	cache := newMockCache()
	cache.data["expire-key"] = "expire-value"

	service := NewCacheService(cache)
	handler := service.ExpireCacheHandler()

	ctx := context.Background()
	cmd := inbound.ExpireCacheCommand{
		Key: "expire-key",
		TTL: 2 * time.Hour,
	}

	err := handler(ctx, cmd)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify TTL set
	ttl, err := cache.TTL(ctx, "expire-key")
	if err != nil {
		t.Fatalf("expected TTL, got error: %v", err)
	}
	if ttl != 2*time.Hour {
		t.Errorf("expected TTL 2h, got %v", ttl)
	}
}

func TestCacheService_ExistsCacheHandler(t *testing.T) {
	cache := newMockCache()
	cache.data["exists-key"] = "exists-value"

	service := NewCacheService(cache)
	handler := service.ExistsCacheHandler()

	ctx := context.Background()

	// Test existing key
	query := inbound.ExistsCacheQuery{Key: "exists-key"}
	exists, err := handler(ctx, query)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !exists {
		t.Error("expected key to exist")
	}

	// Test missing key
	query = inbound.ExistsCacheQuery{Key: "missing-key"}
	exists, err = handler(ctx, query)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if exists {
		t.Error("expected key to not exist")
	}
}

// Benchmark

func BenchmarkCacheService_SetCacheHandler(b *testing.B) {
	cache := newMockCache()
	service := NewCacheService(cache)
	handler := service.SetCacheHandler()

	ctx := context.Background()
	cmd := inbound.SetCacheCommand{
		Key:   "bench-key",
		Value: "bench-value",
		TTL:   time.Hour,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler(ctx, cmd)
	}
}
