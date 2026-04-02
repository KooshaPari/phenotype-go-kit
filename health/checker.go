package health

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// CheckResult represents the result of a health check.
type CheckResult struct {
	Name     string                 `json:"name"`
	Status   string                 `json:"status"` // "healthy", "degraded", "unhealthy"
	Message  string                 `json:"message,omitempty"`
	Duration string                 `json:"duration_ms,omitempty"`
	Details  map[string]interface{} `json:"details,omitempty"`
}

// Checker defines the interface for health checks.
type Checker interface {
	Name() string
	Check(ctx context.Context) CheckResult
}

// HealthChecker aggregates multiple health checkers.
type HealthChecker struct {
	lock      sync.RWMutex
	checks    []Checker
	timeout   time.Duration
	lastCheck map[string]CheckResult
}

// NewHealthChecker creates a new health checker.
func NewHealthChecker(timeout time.Duration) *HealthChecker {
	return &HealthChecker{
		checks:    make([]Checker, 0),
		timeout:   timeout,
		lastCheck: make(map[string]CheckResult),
	}
}

// Register adds a health check.
func (h *HealthChecker) Register(checker Checker) {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.checks = append(h.checks, checker)
}

// RunAll runs all registered health checks.
func (h *HealthChecker) RunAll(ctx context.Context) []CheckResult {
	h.lock.RLock()
	defer h.lock.RUnlock()

	results := make([]CheckResult, 0, len(h.checks))
	for _, checker := range h.checks {
		result := h.runCheck(ctx, checker)
		results = append(results, result)
		h.lastCheck[checker.Name()] = result
	}

	return results
}

func (h *HealthChecker) runCheck(ctx context.Context, checker Checker) CheckResult {
	start := time.Now()
	resultCtx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	result := checker.Check(resultCtx)
	result.Duration = time.Since(start).String()

	return result
}

// LastResults returns the last check results.
func (h *HealthChecker) LastResults() map[string]CheckResult {
	h.lock.RLock()
	defer h.lock.RUnlock()

	result := make(map[string]CheckResult, len(h.lastCheck))
	for k, v := range h.lastCheck {
		result[k] = v
	}
	return result
}

// LivenessHandler returns a handler for liveness probes.
func LivenessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

// ReadinessHandler returns a handler for readiness probes.
func ReadinessHandler(h *HealthChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		results := h.RunAll(r.Context())

		// Check if any check is unhealthy
		hasUnhealthy := false
		for _, result := range results {
			if result.Status == "unhealthy" {
				hasUnhealthy = true
				break
			}
		}

		w.Header().Set("Content-Type", "application/json")

		if hasUnhealthy {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status": "unhealthy",
				"checks": results,
			})
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status": "healthy",
				"checks": results,
			})
		}
	}
}

// JSONHandler returns health status as JSON.
func JSONHandler(h *HealthChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		results := h.RunAll(r.Context())

		// Determine overall status
		overallStatus := "healthy"
		for _, result := range results {
			if result.Status == "unhealthy" {
				overallStatus = "unhealthy"
				break
			}
			if result.Status == "degraded" {
				overallStatus = "degraded"
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    overallStatus,
			"timestamp": time.Now().Format(time.RFC3339),
			"checks":    results,
		})
	}
}

// DatabaseChecker implements a database health check.
type DatabaseChecker struct {
	name    string
	checkFn func(ctx context.Context) error
}

// NewDatabaseChecker creates a new database health checker.
func NewDatabaseChecker(name string, checkFn func(ctx context.Context) error) *DatabaseChecker {
	return &DatabaseChecker{name: name, checkFn: checkFn}
}

func (c *DatabaseChecker) Name() string {
	return c.name
}

func (c *DatabaseChecker) Check(ctx context.Context) CheckResult {
	err := c.checkFn(ctx)
	if err != nil {
		return CheckResult{
			Name:    c.name,
			Status:  "unhealthy",
			Message: err.Error(),
		}
	}
	return CheckResult{
		Name:   c.name,
		Status: "healthy",
	}
}

// RedisChecker implements a Redis health check.
type RedisChecker struct {
	name    string
	checkFn func(ctx context.Context) error
}

// NewRedisChecker creates a new Redis health checker.
func NewRedisChecker(name string, checkFn func(ctx context.Context) error) *RedisChecker {
	return &RedisChecker{name: name, checkFn: checkFn}
}

func (c *RedisChecker) Name() string {
	return c.name
}

func (c *RedisChecker) Check(ctx context.Context) CheckResult {
	err := c.checkFn(ctx)
	if err != nil {
		return CheckResult{
			Name:    c.name,
			Status:  "unhealthy",
			Message: err.Error(),
		}
	}
	return CheckResult{
		Name:   c.name,
		Status: "healthy",
	}
}

// ComponentChecker implements a generic component health check.
type ComponentChecker struct {
	name    string
	checkFn func(ctx context.Context) error
}

// NewComponentChecker creates a new component health checker.
func NewComponentChecker(name string, checkFn func(ctx context.Context) error) *ComponentChecker {
	return &ComponentChecker{name: name, checkFn: checkFn}
}

func (c *ComponentChecker) Name() string {
	return c.name
}

func (c *ComponentChecker) Check(ctx context.Context) CheckResult {
	err := c.checkFn(ctx)
	if err != nil {
		return CheckResult{
			Name:    c.name,
			Status:  "unhealthy",
			Message: err.Error(),
		}
	}
	return CheckResult{
		Name:   c.name,
		Status: "healthy",
	}
}
