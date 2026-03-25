// Package valueobjects contains immutable value types.
// Value objects have no identity and are compared by their attributes.
package valueobjects

import "time"

// AlertRule defines an alerting rule configuration.
type AlertRule struct {
	ID          string
	Name        string
	Condition   Condition
	Severity    string
	Labels      map[string]string
	Annotations map[string]string
	For         time.Duration
}

// Condition defines the triggering condition for an alert.
type Condition struct {
	Metric    string
	Operator  string
	Threshold float64
	Duration  time.Duration
}

// RetryPolicy defines how to handle retries.
type RetryPolicy struct {
	MaxAttempts int
	Backoff     BackoffStrategy
	Timeout     time.Duration
}

// BackoffStrategy defines the backoff algorithm.
type BackoffStrategy string

const (
	BackoffLinear    BackoffStrategy = "linear"
	BackoffExponential BackoffStrategy = "exponential"
	BackoffConstant  BackoffStrategy = "constant"
)

// CircuitState represents the circuit breaker state.
type CircuitState string

const (
	CircuitStateClosed   CircuitState = "closed"
	CircuitStateOpen     CircuitState = "open"
	CircuitStateHalfOpen CircuitState = "half-open"
)

// RateLimitConfig defines rate limiting parameters.
type RateLimitConfig struct {
	RequestsPerSecond int
	Burst             int
	Strategy          string
}
