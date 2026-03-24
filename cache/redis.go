package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrKeyNotFound = errors.New("key not found")
	ErrInvalidTTL  = errors.New("invalid TTL")
)

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

// Cache provides Redis caching functionality.
type Cache struct {
	client *redis.Client
	prefix string
}

// New creates a new Redis cache instance.
func New(cfg Config) *Cache {
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

	return &Cache{
		client: client,
		prefix: "phenotype:",
	}
}

// NewFromClient creates a cache from an existing Redis client.
func NewFromClient(client *redis.Client) *Cache {
	return &Cache{
		client: client,
		prefix: "phenotype:",
	}
}

// Ping checks Redis connectivity.
func (c *Cache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

// Close closes the Redis connection.
func (c *Cache) Close() error {
	return c.client.Close()
}

// Get retrieves a value by key.
func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, c.prefix+key).Result()
	if err == redis.Nil {
		return "", ErrKeyNotFound
	}
	return val, err
}

// Set stores a value with TTL.
func (c *Cache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	if ttl <= 0 {
		return ErrInvalidTTL
	}
	return c.client.Set(ctx, c.prefix+key, value, ttl).Err()
}

// SetNX stores a value only if key doesn't exist.
func (c *Cache) SetNX(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
	return c.client.SetNX(ctx, c.prefix+key, value, ttl).Result()
}

// Delete removes a key.
func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, c.prefix+key).Err()
}

// Exists checks if a key exists.
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Exists(ctx, c.prefix+key).Result()
	return n > 0, err
}

// Expire sets expiration on a key.
func (c *Cache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return c.client.Expire(ctx, c.prefix+key, ttl).Err()
}

// TTL returns remaining TTL for a key.
func (c *Cache) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.client.TTL(ctx, c.prefix+key).Result()
}

// Incr increments a counter.
func (c *Cache) Incr(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, c.prefix+key).Result()
}

// Decr decrements a counter.
func (c *Cache) Decr(ctx context.Context, key string) (int64, error) {
	return c.client.Decr(ctx, c.prefix+key).Result()
}

// GetJSON retrieves and unmarshals JSON.
func (c *Cache) GetJSON(ctx context.Context, key string, dest interface{}) error {
	val, err := c.Get(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

// SetJSON marshals and stores JSON.
func (c *Cache) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Set(ctx, key, string(data), ttl)
}

// DeletePattern deletes all keys matching a pattern.
func (c *Cache) DeletePattern(ctx context.Context, pattern string) (int64, error) {
	keys, err := c.client.Keys(ctx, c.prefix+pattern).Result()
	if err != nil {
		return 0, err
	}
	if len(keys) == 0 {
		return 0, nil
	}
	return c.client.Del(ctx, keys...).Result()
}

// SortedSetAdd adds to a sorted set.
func (c *Cache) SortedSetAdd(ctx context.Context, key string, member string, score float64) error {
	return c.client.ZAdd(ctx, c.prefix+key, redis.Z{Score: score, Member: member}).Err()
}

// SortedSetRange retrieves sorted set members by rank.
func (c *Cache) SortedSetRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.client.ZRange(ctx, c.prefix+key, start, stop).Result()
}

// SortedSetRemove removes from a sorted set.
func (c *Cache) SortedSetRemove(ctx context.Context, key string, members ...interface{}) error {
	formatted := make([]interface{}, len(members))
	for i, m := range members {
		formatted[i] = m
	}
	return c.client.ZRem(ctx, c.prefix+key, formatted...).Err()
}

// SetPrefix sets the key prefix.
func (c *Cache) SetPrefix(prefix string) {
	c.prefix = prefix
}

// Client returns the underlying Redis client.
func (c *Cache) Client() *redis.Client {
	return c.client
}

// HealthCheck returns cache health status.
func (c *Cache) HealthCheck(ctx context.Context) error {
	if err := c.Ping(ctx); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}
	return nil
}
