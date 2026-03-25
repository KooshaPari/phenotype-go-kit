package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/KooshaPari/phenotype-go-kit/contracts/ports/inbound"
github.com/KooshaPari/phenotype-go-kit/contracts/ports/inbound
	"github.com/KooshaPari/phenotype-go-kit/contracts/ports/outbound"

// MockCachePort is a mock implementation of CacheJSONPort for testing.
type MockCachePort struct {
	data      map[string]string
	ttls      map[string]time.Duration
	getErr    error
	setErr    error
	deleteErr error
	expireErr error
}

func NewMockCache() *MockCachePort {
	return &MockCachePort{
		data:  make(map[string]string),
		ttls:  make(map[string]time.Duration),
	}
}

func (m *MockCachePort) Get(ctx context.Context, key string) (string, error) {
	if m.getErr != nil {
		return "", m.getErr
	}
	v, ok := m.data[key]
	if !ok {
		return "", outbound.ErrKeyNotFound
	}
	return v, nil
}

func (m *MockCachePort) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	if m.setErr != nil {
		return m.setErr
	}
	m.data[key] = value
	m.ttls[key] = ttl
	return nil
}

func (m *MockCachePort) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return m.Set(ctx, key, serializeJSON(value), ttl)
}

func (m *MockCachePort) GetJSON(ctx context.Context, key string, dest interface{}) error {
	v, err := m.Get(ctx, key)
	if err != nil {
		return err
	}
	return deserializeJSON(v, dest)
}

func (m *MockCachePort) Delete(ctx context.Context, key string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	delete(m.data, key)
	delete(m.ttls, key)
	return nil
}

func (m *MockCachePort) Exists(ctx context.Context, key string) (bool, error) {
	_, ok := m.data[key]
	return ok, nil
}

func (m *MockCachePort) Expire(ctx context.Context, key string, ttl time.Duration) error {
	if m.expireErr != nil {
		return m.expireErr
	}
	m.ttls[key] = ttl
	return nil
}

func (m *MockCachePort) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, ok := m.ttls[key]
	if !ok {
		return 0, outbound.ErrKeyNotFound
	}
	return ttl, nil
}

func (m *MockCachePort) GetOrSet(ctx context.Context, key string, compute func() (string, error), ttl time.Duration) (string, error) {
	if v, ok := m.data[key]; ok {
		return v, nil
	}
	v, err := compute()
	if err != nil {
		return "", err
	}
	_ = m.Set(ctx, key, v, ttl)
	return v, nil
}

func serializeJSON(v interface{}) string {
	return `{"value":"` + toString(v) + `"}`
}

func toString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func deserializeJSON(data string, dest interface{}) error {
	return nil
}

// Test helper functions
func TestCacheService_SetCache_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	cache := NewMockCache()
	service := NewCacheService(cache)

	// Act - Using the handler function pattern
	handler := service.SetCacheHandler()
	cmd := inbound.SetCacheCommand{
		Key:   "test-key",
		Value: "test-value",
		TTL:   5 * time.Minute,
	}

	// Assert - Execute handler
	err := handler(ctx, cmd)

	// Verify
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	val, err := cache.Get(ctx, "test-key")
	if err != nil {
		t.Errorf("expected value in cache, got error: %v", err)
	}
	if val != "test-value" {
		t.Errorf("expected 'test-value', got '%s'", val)
	}
}

func TestCacheService_SetCache_InvalidTTL(t *testing.T) {
	// Arrange
	ctx := context.Background()
	cache := NewMockCache()
	service := NewCacheService(cache)

	// Act
	handler := service.SetCacheHandler()
	cmd := inbound.SetCacheCommand{
		Key:   "test-key",
		Value: "test-value",
		TTL:   0, // Invalid TTL
	}

	// Assert
	err := handler(ctx, cmd)

	if err != ErrInvalidTTL {
		t.Errorf("expected ErrInvalidTTL, got %v", err)
	}
}

func TestCacheService_SetCache_NegativeTTL(t *testing.T) {
	// Arrange
	ctx := context.Background()
	cache := NewMockCache()
	service := NewCacheService(cache)

	// Act
	handler := service.SetCacheHandler()
	cmd := inbound.SetCacheCommand{
		Key:   "test-key",
		Value: "test-value",
		TTL:   -1 * time.Second, // Negative TTL
	}

	// Assert
	err := handler(ctx, cmd)

	if err != ErrInvalidTTL {
		t.Errorf("expected ErrInvalidTTL, got %v", err)
	}
}

func TestCacheService_DeleteCache_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	cache := NewMockCache()
	_ = cache.Set(ctx, "test-key", "test-value", time.Minute)
	service := NewCacheService(cache)

	// Act
	handler := service.DeleteCacheHandler()
	cmd := inbound.DeleteCacheCommand{Key: "test-key"}
	err := handler(ctx, cmd)

	// Assert
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	exists, _ := cache.Exists(ctx, "test-key")
	if exists {
		t.Error("expected key to be deleted")
	}
}

func TestCacheService_GetCache_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	cache := NewMockCache()
	_ = cache.Set(ctx, "test-key", "test-value", time.Minute)
	service := NewCacheService(cache)

	// Act
	handler := service.GetCacheHandler()
	query := inbound.GetCacheQuery{Key: "test-key"}
	val, err := handler(ctx, query)

	// Assert
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if val != "test-value" {
		t.Errorf("expected 'test-value', got '%s'", val)
	}
}

func TestCacheService_GetCache_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	cache := NewMockCache()
	service := NewCacheService(cache)

	// Act
	handler := service.GetCacheHandler()
	query := inbound.GetCacheQuery{Key: "nonexistent"}
	_, err := handler(ctx, query)

	// Assert
	if !errors.Is(err, ErrKeyNotFound) {
		t.Errorf("expected ErrKeyNotFound, got %v", err)
	}
}

func TestCacheService_Exists_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	cache := NewMockCache()
	_ = cache.Set(ctx, "test-key", "test-value", time.Minute)
	service := NewCacheService(cache)

	// Act
	handler := service.ExistsCacheHandler()
	query := inbound.ExistsCacheQuery{Key: "test-key"}
	exists, err := handler(ctx, query)

	// Assert
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !exists {
		t.Error("expected key to exist")
	}
}

func TestCacheService_Exists_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	cache := NewMockCache()
	service := NewCacheService(cache)

	// Act
	handler := service.ExistsCacheHandler()
	query := inbound.ExistsCacheQuery{Key: "nonexistent"}
	exists, err := handler(ctx, query)

	// Assert
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if exists {
		t.Error("expected key to not exist")
	}
}

func TestCacheService_Expire_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	cache := NewMockCache()
	_ = cache.Set(ctx, "test-key", "test-value", time.Minute)
	service := NewCacheService(cache)

	// Act
	handler := service.ExpireCacheHandler()
	cmd := inbound.ExpireCacheCommand{Key: "test-key", TTL: 10 * time.Minute}
	err := handler(ctx, cmd)

	// Assert
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// CacheWarmer Tests
func TestCacheWarmer_WarmKey_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	cache := NewMockCache()
	fetcher := func(ctx context.Context, key string) (string, time.Duration, error) {
		return "fetched-value", 5 * time.Minute, nil
	}
	warmer := NewCacheWarmer(cache, fetcher)

	// Act
	err := warmer.WarmKey(ctx, "warm-key")

	// Assert
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	val, _ := cache.Get(ctx, "warm-key")
	if val != "fetched-value" {
		t.Errorf("expected 'fetched-value', got '%s'", val)
	}
}

func TestCacheWarmer_WarmKey_FetchError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	cache := NewMockCache()
	expectedErr := errors.New("fetch error")
	fetcher := func(ctx context.Context, key string) (string, time.Duration, error) {
		return "", 0, expectedErr
	}
	warmer := NewCacheWarmer(cache, fetcher)

	// Act
	err := warmer.WarmKey(ctx, "warm-key")

	// Assert
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected fetch error, got %v", err)
	}
}

func TestCacheWarmer_WarmKeys_SkipsOnError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	cache := NewMockCache()
	fetcher := func(ctx context.Context, key string) (string, time.Duration, error) {
		if key == "fail-key" {
			return "", 0, errors.New("fetch failed")
		}
		return "value-" + key, time.Minute, nil
	}
	warmer := NewCacheWarmer(cache, fetcher)

	// Act
	err := warmer.WarmKeys(ctx, []string{"key1", "fail-key", "key2"})

	// Assert - Should not return error (skips failed keys)
	if err != nil {
		t.Errorf("expected no error (skip), got %v", err)
	}

	// Verify key1 and key2 were cached
	_, err = cache.Get(ctx, "key1")
	if err != nil {
		t.Errorf("expected key1 to be cached, got error")
	}
	_, err = cache.Get(ctx, "key2")
	if err != nil {
		t.Errorf("expected key2 to be cached, got error")
	}
}

// PatternInvalidator Tests
func TestPatternInvalidator_Invalidate_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	cache := NewMockCache()
	_ = cache.Set(ctx, "user:1", "value1", time.Minute)
	_ = cache.Set(ctx, "user:2", "value2", time.Minute)
	_ = cache.Set(ctx, "order:1", "value3", time.Minute)

	invalidator := &PatternInvalidator{}

	// Act
	count, err := invalidator.Invalidate(ctx, cache, "user:*")

	// Assert
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 invalidated, got %d", count)
	}
}

// Benchmark Tests
func BenchmarkCacheService_SetCache(b *testing.B) {
	ctx := context.Background()
	cache := NewMockCache()
	service := NewCacheService(cache)
	handler := service.SetCacheHandler()
	cmd := inbound.SetCacheCommand{
		Key:   "bench-key",
		Value: "bench-value",
		TTL:   time.Minute,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = handler(ctx, cmd)
	}
}

func BenchmarkCacheService_GetCache(b *testing.B) {
	ctx := context.Background()
	cache := NewMockCache()
	_ = cache.Set(ctx, "bench-key", "bench-value", time.Minute)
	service := NewCacheService(cache)
	handler := service.GetCacheHandler()
	query := inbound.GetCacheQuery{Key: "bench-key"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = handler(ctx, query)
	}
}
