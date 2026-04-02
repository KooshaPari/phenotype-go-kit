package jobs

import (
	"context"
	"encoding/json"
	"time"
)

// Job represents a background job in the queue.
type Job struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"`
	Payload    json.RawMessage `json:"payload"`
	Retries    int             `json:"retries"`
	MaxRetries int             `json:"max_retries"`
	Status     JobStatus       `json:"status"`
	CreatedAt  time.Time       `json:"created_at"`
	StartedAt  *time.Time      `json:"started_at,omitempty"`
	FailedAt   *time.Time      `json:"failed_at,omitempty"`
	Error      string          `json:"error,omitempty"`
}

// JobStatus represents the current state of a job.
type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
	JobStatusRetryable JobStatus = "retryable"
)

// JobHandler is a function that processes a job.
type JobHandler func(ctx context.Context, job *Job) error

// Registry holds job type to handler mappings.
type Registry map[string]JobHandler
