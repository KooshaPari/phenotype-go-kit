package adapter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/KooshaPari/phenotype-go-kit/contracts/ports/outbound"
	"github.com/redis/go-redis/v9"
)

var (
	_ outbound.CachePort        = (*RedisCacheAdapter)(nil)
	_ outbound.CacheJSONPort    = (*RedisCacheAdapter)(nil)
	_ outbound.CacheCounterPort = (*RedisCacheAdapter)(nil)
)

// ErrKeyNotFound is returned when a key does not exist.
var ErrKeyNotFound = errors.New("key not found")

// ErrInvalidTTL is returned when TTL is invalid.
var ErrInvalidTTL = errors.New("invalid TTL")

// Config holds Redis cache configuration.
type Config struct {
	Addr         string        `default:"localhost:6379"`
	Password     string        `default:""`
	DB           int           `default:"0"`
	PoolSize     int           `default:"10"`
	MinIdleConns int           `default:"5"`
	DialTimeout  time.Duration `default:"5s"`
	ReadTimeout  time.Duration `default:"3s"`
	WriteTimeout time.Duration `default:"3s"`
}

// RedisCacheAdapter implements outbound.CachePort using Redis.
type RedisCacheAdapter struct {
	client *redis.Client
	prefix string
}

// NewRedisCacheAdapter creates a new Redis cache adapter.
func NewRedisCacheAdapter(cfg Config) *RedisCacheAdapter {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	return &RedisCacheAdapter{
		client: client,
		prefix: "phenotype:",
	}
}

// NewRedisCacheAdapterFromClient creates adapter from existing Redis client.
func NewRedisCacheAdapterFromClient(client *redis.Client) *RedisCacheAdapter {
	return &RedisCacheAdapter{
		client: client,
		prefix: "phenotype:",
	}
}

// Get retrieves a value by key.
func (a *RedisCacheAdapter) Get(ctx context.Context, key string) (string, error) {
	val, err := a.client.Get(ctx, a.prefix+key).Result()
	if err == redis.Nil {
		return "", ErrKeyNotFound
	}
	return val, err
}

// Set stores a value with TTL.
func (a *RedisCacheAdapter) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	if ttl <= 0 {
		return ErrInvalidTTL
	}
	return a.client.Set(ctx, a.prefix+key, value, ttl).Err()
}

// SetNX stores a value only if key doesn't exist.
func (a *RedisCacheAdapter) SetNX(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
	return a.client.SetNX(ctx, a.prefix+key, value, ttl).Result()
}

// Delete removes a key.
func (a *RedisCacheAdapter) Delete(ctx context.Context, key string) error {
	return a.client.Del(ctx, a.prefix+key).Err()
}

// Exists checks if a key exists.
func (a *RedisCacheAdapter) Exists(ctx context.Context, key string) (bool, error) {
	n, err := a.client.Exists(ctx, a.prefix+key).Result()
	return n > 0, err
}

// Expire sets expiration on a key.
func (a *RedisCacheAdapter) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return a.client.Expire(ctx, a.prefix+key, ttl).Err()
}

// TTL returns remaining TTL for a key.
func (a *RedisCacheAdapter) TTL(ctx context.Context, key string) (time.Duration, error) {
	return a.client.TTL(ctx, a.prefix+key).Result()
}

// Ping checks connectivity.
func (a *RedisCacheAdapter) Ping(ctx context.Context) error {
	return a.client.Ping(ctx).Err()
}

// Close closes the Redis connection.
func (a *RedisCacheAdapter) Close() error {
	return a.client.Close()
}

// Incr increments a counter.
func (a *RedisCacheAdapter) Incr(ctx context.Context, key string) (int64, error) {
	return a.client.Incr(ctx, a.prefix+key).Result()
}

// Decr decrements a counter.
func (a *RedisCacheAdapter) Decr(ctx context.Context, key string) (int64, error) {
	return a.client.Decr(ctx, a.prefix+key).Result()
}

// GetJSON retrieves and unmarshals JSON.
func (a *RedisCacheAdapter) GetJSON(ctx context.Context, key string, dest interface{}) error {
	val, err := a.Get(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

// SetJSON marshals and stores JSON.
func (a *RedisCacheAdapter) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return a.Set(ctx, key, string(data), ttl)
}

// DeletePattern deletes all keys matching a pattern.
func (a *RedisCacheAdapter) DeletePattern(ctx context.Context, pattern string) (int64, error) {
	keys, err := a.client.Keys(ctx, a.prefix+pattern).Result()
	if err != nil {
		return 0, err
	}
	if len(keys) == 0 {
		return 0, nil
	}
	return a.client.Del(ctx, keys...).Result()
}

// SortedSetAdd adds to a sorted set.
func (a *RedisCacheAdapter) SortedSetAdd(ctx context.Context, key string, member string, score float64) error {
	return a.client.ZAdd(ctx, a.prefix+key, redis.Z{Score: score, Member: member}).Err()
}

// SortedSetRange retrieves sorted set members by rank.
func (a *RedisCacheAdapter) SortedSetRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return a.client.ZRange(ctx, a.prefix+key, start, stop).Result()
}

// SortedSetRemove removes from a sorted set.
func (a *RedisCacheAdapter) SortedSetRemove(ctx context.Context, key string, members ...string) error {
	formatted := make([]interface{}, len(members))
	for i, m := range members {
		formatted[i] = m
	}
	return a.client.ZRem(ctx, a.prefix+key, formatted...).Err()
}

// Client returns the underlying Redis client.
func (a *RedisCacheAdapter) Client() *redis.Client {
	return a.client
}

// HealthCheck returns cache health status.
func (a *RedisCacheAdapter) HealthCheck(ctx context.Context) error {
	if err := a.Ping(ctx); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}
	return nil
}
