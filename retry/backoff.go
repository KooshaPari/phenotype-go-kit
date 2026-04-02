package retry

import (
	"context"
	"errors"
	"log/slog"
	"time"
)

// Config holds retry configuration.
type Config struct {
	MaxAttempts  int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
	Jitter       bool
}

// DefaultConfig returns default retry configuration.
func DefaultConfig() Config {
	return Config{
		MaxAttempts:  3,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     30 * time.Second,
		Multiplier:   2.0,
		Jitter:       true,
	}
}

// RetryFunc is a function that can be retried.
type RetryFunc func(ctx context.Context) error

// Do retries a function with exponential backoff.
func Do(ctx context.Context, cfg Config, fn RetryFunc) error {
	var lastErr error
	delay := cfg.InitialDelay

	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		if err := fn(ctx); err != nil {
			lastErr = err

			// Check if we should stop
			if attempt >= cfg.MaxAttempts {
				break
			}

			// Check for context cancellation
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			// Log retry attempt
			slog.Default().Debug("retry attempt",
				"attempt", attempt,
				"max", cfg.MaxAttempts,
				"delay", delay,
				"error", err,
			)

			// Wait before next attempt
			waitDuration := delay
			if cfg.Jitter {
				waitDuration = withJitter(delay)
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(waitDuration):
			}

			// Increase delay
			delay = time.Duration(float64(delay) * cfg.Multiplier)
			if delay > cfg.MaxDelay {
				delay = cfg.MaxDelay
			}
		} else {
			return nil
		}
	}

	return errors.New("max retries exceeded: " + lastErr.Error())
}

// WithRetry wraps a function with retry logic.
func WithRetry(cfg Config) func(RetryFunc) RetryFunc {
	return func(fn RetryFunc) RetryFunc {
		return func(ctx context.Context) error {
			return Do(ctx, cfg, fn)
		}
	}
}

// DoWithResult retries a function that returns a result.
func DoWithResult(ctx context.Context, cfg Config, fn func() (interface{}, error)) (interface{}, error) {
	var lastErr error
	delay := cfg.InitialDelay

	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		result, err := fn()
		if err != nil {
			lastErr = err

			if attempt >= cfg.MaxAttempts {
				break
			}

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}

			waitDuration := delay
			if cfg.Jitter {
				waitDuration = withJitter(delay)
			}

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(waitDuration):
			}

			delay = time.Duration(float64(delay) * cfg.Multiplier)
			if delay > cfg.MaxDelay {
				delay = cfg.MaxDelay
			}
		} else {
			return result, nil
		}
	}

	return nil, errors.New("max retries exceeded: " + lastErr.Error())
}

func withJitter(d time.Duration) time.Duration {
	// Add +/- 25% jitter
	jitter := int64(float64(d) * 0.25)
	offset := (time.Now().UnixNano() % (jitter * 2)) - jitter
	return d + time.Duration(offset)
}

// Backoff calculates the next backoff delay.
func Backoff(attempt int, cfg Config) time.Duration {
	delay := cfg.InitialDelay
	for i := 1; i < attempt; i++ {
		delay = time.Duration(float64(delay) * cfg.Multiplier)
		if delay > cfg.MaxDelay {
			delay = cfg.MaxDelay
			break
		}
	}
	return delay
}

// PermanentError marks an error as non-retryable.
type PermanentError struct {
	Err error
}

func (e *PermanentError) Error() string {
	return e.Err.Error()
}

func (e *PermanentError) Unwrap() error {
	return e.Err
}

// IsPermanent checks if an error is permanent.
func IsPermanent(err error) bool {
	var pe *PermanentError
	return errors.As(err, &pe)
}
