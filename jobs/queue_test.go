package jobs

import (
	"context"
	"log/slog"
	"testing"
	"time"
)

func TestQueue_Enqueue(t *testing.T) {
	handlers := Registry{
		"test": func(ctx context.Context, job *Job) error {
			return nil
		},
	}

	logger := slog.Default()
	queue := NewQueue(QueueConfig{Workers: 2}, handlers, logger)

	ctx := cancelAfter(100 * time.Millisecond)
	queue.Start(ctx)

	job := &Job{
		Type:    "test",
		Payload: []byte(`{"key":"value"}`),
	}

	err := queue.Enqueue(ctx, job)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	retrieved, ok := queue.GetJob(job.ID)
	if !ok {
		t.Fatal("job not found in queue")
	}

	if retrieved.Status != JobStatusCompleted {
		t.Errorf("expected job status completed, got %v", retrieved.Status)
	}

	queue.Stop(ctx)
}

func TestQueue_Retry(t *testing.T) {
	failCount := 0
	handlers := Registry{
		"failing": func(ctx context.Context, job *Job) error {
			failCount++
			if failCount < 2 {
				return &RetryableError{"temporary failure"}
			}
			return nil
		},
	}

	logger := slog.Default()
	queue := NewQueue(QueueConfig{Workers: 1, Interval: 10 * time.Millisecond}, handlers, logger)

	ctx := cancelAfter(500 * time.Millisecond)
	queue.Start(ctx)

	job := &Job{
		Type:       "failing",
		Payload:    []byte(`{}`),
		MaxRetries: 3,
	}

	err := queue.Enqueue(ctx, job)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	time.Sleep(400 * time.Millisecond)
	queue.Stop(ctx)

	if failCount != 2 {
		t.Errorf("expected 2 attempts, got %d", failCount)
	}
}

func TestNewEmailJob(t *testing.T) {
	payload := EmailPayload{
		To:      "test@example.com",
		From:    "noreply@phenotype.dev",
		Subject: "Test Email",
		Body:    "Hello World",
	}

	job, err := NewEmailJob(payload)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if job.Type != "email" {
		t.Errorf("expected type 'email', got %s", job.Type)
	}

	// ID is generated when job is enqueued, not in the helper
	// job.ID can be empty at this point
	if job.Payload == nil {
		t.Error("expected payload to be set")
	}
}

func TestNewSMSJob(t *testing.T) {
	payload := SMSPayload{
		To:      "+1234567890",
		Message: "Hello World",
	}

	job, err := NewSMSJob(payload)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if job.Type != "sms" {
		t.Errorf("expected type 'sms', got %s", job.Type)
	}
}

func cancelAfter(d time.Duration) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(d)
		cancel()
	}()
	return ctx
}

type RetryableError struct {
	msg string
}

func (e *RetryableError) Error() string {
	return e.msg
}
