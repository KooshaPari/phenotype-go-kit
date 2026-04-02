package circuit

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"
)

var (
	ErrCircuitOpen     = errors.New("circuit breaker is open")
	ErrCircuitHalfOpen = errors.New("circuit breaker is half-open")
)

// State represents the circuit breaker state.
type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

// Config holds circuit breaker configuration.
type Config struct {
	FailureThreshold int           `default:"5"`
	SuccessThreshold int           `default:"2"`
	Timeout          time.Duration `default:"30s"`
	RequestTimeout   time.Duration `default:"10s"`
}

// Breaker implements the circuit breaker pattern.
type Breaker struct {
	name    string
	config  Config
	state   State
	fails   int
	success int

	mu           sync.Mutex
	lastFailTime time.Time

	logger *slog.Logger
}

// New creates a new circuit breaker.
func New(name string, cfg Config) *Breaker {
	return &Breaker{
		name:   name,
		config: cfg,
		state:  StateClosed,
		logger: slog.Default(),
	}
}

// Execute runs a function with circuit breaker protection.
func (cb *Breaker) Execute(ctx context.Context, fn func() error) error {
	cb.mu.Lock()

	switch cb.state {
	case StateOpen:
		// Check if timeout has passed
		if time.Since(cb.lastFailTime) > cb.config.Timeout {
			cb.logger.Info("circuit breaker transitioning to half-open", "name", cb.name)
			cb.state = StateHalfOpen
			cb.success = 0
		} else {
			cb.mu.Unlock()
			return ErrCircuitOpen
		}
	case StateHalfOpen:
		// Allow one request through
	}

	cb.mu.Unlock()

	// Execute with timeout
	runCtx, cancel := context.WithTimeout(ctx, cb.config.RequestTimeout)
	defer cancel()

	result := make(chan error, 1)
	go func() {
		result <- fn()
	}()

	select {
	case <-runCtx.Done():
		cb.recordFailure()
		return runCtx.Err()
	case err := <-result:
		if err != nil {
			cb.recordFailure()
			return err
		}
		cb.recordSuccess()
		return nil
	}
}

func (cb *Breaker) recordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.fails = 0

	if cb.state == StateHalfOpen {
		cb.success++
		if cb.success >= cb.config.SuccessThreshold {
			cb.logger.Info("circuit breaker closing", "name", cb.name)
			cb.state = StateClosed
			cb.success = 0
		}
	}
}

func (cb *Breaker) recordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.fails++
	cb.lastFailTime = time.Now()

	if cb.state == StateHalfOpen {
		cb.logger.Info("circuit breaker opening from half-open", "name", cb.name)
		cb.state = StateOpen
	} else if cb.fails >= cb.config.FailureThreshold {
		cb.logger.Info("circuit breaker opening", "name", cb.name)
		cb.state = StateOpen
	}
}

// State returns the current circuit state.
func (cb *Breaker) State() State {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

// Reset clears the circuit breaker state.
func (cb *Breaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.state = StateClosed
	cb.fails = 0
	cb.success = 0
}

// Metrics holds circuit breaker metrics.
type Metrics struct {
	TotalRequests int
	Failures      int
	Successes     int
	CurrentState  State
	LastFailTime  time.Time
}

// GetMetrics returns current metrics.
func (cb *Breaker) GetMetrics() Metrics {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	return Metrics{
		CurrentState: cb.state,
	}
}

// MultiBreaker manages multiple circuit breakers.
type MultiBreaker struct {
	breakers map[string]*Breaker
	mu       sync.RWMutex
}

// NewMultiBreaker creates a multi-breaker manager.
func NewMultiBreaker() *MultiBreaker {
	return &MultiBreaker{
		breakers: make(map[string]*Breaker),
	}
}

// Get returns or creates a circuit breaker.
func (mb *MultiBreaker) Get(name string, cfg Config) *Breaker {
	mb.mu.RLock()
	b, ok := mb.breakers[name]
	mb.mu.RUnlock()

	if ok {
		return b
	}

	mb.mu.Lock()
	defer mb.mu.Unlock()

	if b, ok = mb.breakers[name]; ok {
		return b
	}

	b = New(name, cfg)
	mb.breakers[name] = b
	return b
}

// WithCircuit wraps a function with circuit breaker.
func WithCircuit(breaker *Breaker) func(context.Context, func() error) error {
	return func(ctx context.Context, fn func() error) error {
		return breaker.Execute(ctx, fn)
	}
}
