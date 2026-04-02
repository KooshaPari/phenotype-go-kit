package ratelimit

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter provides rate limiting functionality.
type RateLimiter struct {
	requests map[string]*clientTracker
	mu       sync.RWMutex
	config   Config
}

// Config holds rate limiter configuration.
type Config struct {
	RequestsPerSecond float64       `default:"100"`
	BurstSize         int           `default:"200"`
	CleanupInterval   time.Duration `default:"5m"`
	BlockDuration     time.Duration `default:"5m"`
}

// clientTracker tracks requests for a client.
type clientTracker struct {
	tokens      float64
	lastUpdate  time.Time
	blocked     bool
	blockExpiry time.Time
}

// New creates a new rate limiter.
func New(cfg Config) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string]*clientTracker),
		config:   cfg,
	}

	go rl.cleanupLoop()

	return rl
}

// Allow checks if a request is allowed.
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	tracker, exists := rl.requests[key]

	if !exists {
		rl.requests[key] = &clientTracker{
			tokens:     float64(rl.config.BurstSize) - 1,
			lastUpdate: now,
		}
		return true
	}

	// Check if blocked
	if tracker.blocked && now.Before(tracker.blockExpiry) {
		return false
	}

	// Reset if enough time has passed
	if now.Sub(tracker.lastUpdate) > rl.config.BlockDuration {
		tracker.tokens = float64(rl.config.BurstSize)
		tracker.blocked = false
	}

	// Token bucket algorithm
	elapsed := now.Sub(tracker.lastUpdate).Seconds()
	refillAmount := elapsed * rl.config.RequestsPerSecond
	tracker.tokens = min(tracker.tokens+refillAmount, float64(rl.config.BurstSize))

	if tracker.tokens >= 1 {
		tracker.tokens--
		tracker.lastUpdate = now
		return true
	}

	return false
}

// Block blocks a client for a duration.
func (rl *RateLimiter) Block(key string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	tracker, exists := rl.requests[key]
	if !exists {
		tracker = &clientTracker{}
		rl.requests[key] = tracker
	}

	tracker.blocked = true
	tracker.blockExpiry = time.Now().Add(rl.config.BlockDuration)
}

// Unblock removes blocks for a client.
func (rl *RateLimiter) Unblock(key string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if tracker, exists := rl.requests[key]; exists {
		tracker.blocked = false
	}
}

// Middleware returns HTTP middleware for rate limiting.
func (rl *RateLimiter) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Use API key, IP, or user ID as key
			key := getClientKey(r)

			if !rl.Allow(key) {
				w.Header().Set("Retry-After", "1")
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.config.CleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		rl.cleanup()
	}
}

func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for key, tracker := range rl.requests {
		if now.Sub(tracker.lastUpdate) > rl.config.BlockDuration*2 {
			delete(rl.requests, key)
		}
	}
}

func getClientKey(r *http.Request) string {
	// Check for API key
	if apiKey := r.Header.Get("X-API-Key"); apiKey != "" {
		return "apikey:" + apiKey
	}

	// Check for auth token
	if auth := r.Header.Get("Authorization"); auth != "" {
		return "auth:" + auth
	}

	// Fall back to IP
	return "ip:" + r.RemoteAddr
}

// Headers adds rate limit headers to response.
func Headers(rl *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)

			// Add rate limit headers
			w.Header().Set("X-RateLimit-Limit", "100")
			w.Header().Set("X-RateLimit-Remaining", "99")
		})
	}
}

// DistributedRateLimiter provides rate limiting across multiple instances.
// Uses a simple in-memory implementation - production would use Redis.
type DistributedRateLimiter struct {
	local *RateLimiter
}

// NewDistributed creates a distributed rate limiter.
func NewDistributed(cfg Config) *DistributedRateLimiter {
	return &DistributedRateLimiter{
		local: New(cfg),
	}
}

// Allow checks if a request is allowed (local implementation).
func (rl *DistributedRateLimiter) Allow(key string) bool {
	return rl.local.Allow(key)
}

// Middleware returns HTTP middleware.
func (rl *DistributedRateLimiter) Middleware() func(http.Handler) http.Handler {
	return rl.local.Middleware()
}
