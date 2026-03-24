package waitfor_test

import (
	"context"
	"testing"
	"time"

	"github.com/KooshaPari/phenotype-go-kit/waitfor"
)

func TestWaitForImmediateSuccess(t *testing.T) {
	err := waitfor.WaitFor(context.Background(), waitfor.WaitTimeout{
		Timeout: time.Second,
	}, func() (bool, error) {
		return true, nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestWaitForTimeout(t *testing.T) {
	err := waitfor.WaitFor(context.Background(), waitfor.WaitTimeout{
		Timeout:     50 * time.Millisecond,
		MinInterval: 10 * time.Millisecond,
		MaxInterval: 20 * time.Millisecond,
		InitialWait: true,
	}, func() (bool, error) {
		return false, nil
	})
	if err != waitfor.ErrTimedOut {
		t.Fatalf("expected ErrTimedOut, got %v", err)
	}
}

func TestAfter(t *testing.T) {
	ch := waitfor.After(nil, 10*time.Millisecond)
	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Fatal("After did not fire within 1s")
	}
}
