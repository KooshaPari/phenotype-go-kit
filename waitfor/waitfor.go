package waitfor

import (
	"context"
	"fmt"
	"time"

	"github.com/coder/quartz"
)

// WaitTimeout configures the polling behavior of WaitFor.
type WaitTimeout struct {
	Timeout     time.Duration
	MinInterval time.Duration
	MaxInterval time.Duration
	InitialWait bool
	Clock       quartz.Clock
}

// ErrTimedOut is returned when WaitFor exceeds the configured timeout.
var ErrTimedOut = fmt.Errorf("timeout waiting for condition")

// WaitFor polls condition with exponential backoff until it returns true,
// an error, or the timeout expires.
func WaitFor(ctx context.Context, timeout WaitTimeout, condition func() (bool, error)) error {
	clock := timeout.Clock
	if clock == nil {
		clock = quartz.NewReal()
	}

	minInterval := timeout.MinInterval
	maxInterval := timeout.MaxInterval
	timeoutDuration := timeout.Timeout
	if minInterval == 0 {
		minInterval = 10 * time.Millisecond
	}
	if maxInterval == 0 {
		maxInterval = 500 * time.Millisecond
	}
	if timeoutDuration == 0 {
		timeoutDuration = 10 * time.Second
	}
	if minInterval > maxInterval {
		return fmt.Errorf("minInterval is greater than maxInterval")
	}

	timeoutTimer := clock.NewTimer(timeoutDuration)
	defer timeoutTimer.Stop()

	interval := minInterval
	sleepTimer := clock.NewTimer(interval)
	defer sleepTimer.Stop()

	waitForTimer := timeout.InitialWait
	for {
		if waitForTimer {
			select {
			case <-sleepTimer.C:
			case <-ctx.Done():
				return ctx.Err()
			case <-timeoutTimer.C:
				return ErrTimedOut
			}
		}
		waitForTimer = true

		ok, err := condition()
		if err != nil {
			return err
		}
		if ok {
			return nil
		}

		interval = min(interval*2, maxInterval)
		if !sleepTimer.Stop() {
			select {
			case <-sleepTimer.C:
			default:
			}
		}
		sleepTimer.Reset(interval)
	}
}

// After returns a channel that sends the current time after duration d
// using the provided clock. If clk is nil, a real clock is used.
func After(clk quartz.Clock, d time.Duration) <-chan time.Time {
	if clk == nil {
		clk = quartz.NewReal()
	}
	timer := clk.NewTimer(d)
	ch := make(chan time.Time)
	go func() {
		defer timer.Stop()
		defer close(ch)
		ch <- <-timer.C
	}()
	return ch
}
