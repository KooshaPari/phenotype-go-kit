package jobs

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

// Queue implements an in-memory job queue with worker pool.
type Queue struct {
	mu       sync.RWMutex
	jobs     map[string]*Job
	handlers Registry
	workers  int
	interval time.Duration

	// In production, replace with Redis/DB-backed queue
	jobChan  chan *Job
	stopChan chan struct{}
	wg       sync.WaitGroup
	logger   *slog.Logger
}

// QueueConfig holds configuration for the job queue.
type QueueConfig struct {
	Workers  int           `default:"5"`
	Interval time.Duration `default:"100ms"`
}

// NewQueue creates a new job queue with the given configuration.
func NewQueue(cfg QueueConfig, handlers Registry, logger *slog.Logger) *Queue {
	if cfg.Workers <= 0 {
		cfg.Workers = 5
	}
	if cfg.Interval <= 0 {
		cfg.Interval = 100 * time.Millisecond
	}
	return &Queue{
		jobs:     make(map[string]*Job),
		handlers: handlers,
		workers:  cfg.Workers,
		interval: cfg.Interval,
		jobChan:  make(chan *Job, 1000),
		stopChan: make(chan struct{}),
		logger:   logger,
	}
}

// Enqueue adds a job to the queue.
func (q *Queue) Enqueue(ctx context.Context, job *Job) error {
	if job.ID == "" {
		job.ID = generateJobID()
	}
	if job.CreatedAt.IsZero() {
		job.CreatedAt = time.Now()
	}
	if job.MaxRetries == 0 {
		job.MaxRetries = 3
	}
	job.Status = JobStatusPending

	q.mu.Lock()
	q.jobs[job.ID] = job
	q.mu.Unlock()

	select {
	case q.jobChan <- job:
		q.logger.Info("job enqueued", "job_id", job.ID, "type", job.Type)
		return nil
	default:
		return ErrQueueFull
	}
}

// Start begins the worker pool to process jobs.
func (q *Queue) Start(ctx context.Context) {
	for i := 0; i < q.workers; i++ {
		q.wg.Add(1)
		go q.worker(ctx, i)
	}
	q.logger.Info("job queue started", "workers", q.workers)
}

// Stop gracefully shuts down the worker pool.
func (q *Queue) Stop(ctx context.Context) error {
	close(q.stopChan)
	q.wg.Wait()
	q.logger.Info("job queue stopped")
	return nil
}

func (q *Queue) worker(ctx context.Context, id int) {
	defer q.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case <-q.stopChan:
			return
		case job := <-q.jobChan:
			q.processJob(ctx, job)
		}
	}
}

func (q *Queue) processJob(ctx context.Context, job *Job) {
	handler, ok := q.handlers[job.Type]
	if !ok {
		job.Status = JobStatusFailed
		job.Error = "no handler registered for job type"
		q.logger.Error("job handler not found", "job_id", job.ID, "type", job.Type)
		q.updateJob(job)
		return
	}

	now := time.Now()
	job.Status = JobStatusRunning
	job.StartedAt = &now
	q.updateJob(job)

	q.logger.Info("processing job", "job_id", job.ID, "worker", "worker-0")

	err := handler(ctx, job)
	if err != nil {
		q.handleFailure(ctx, job, err)
		return
	}

	job.Status = JobStatusCompleted
	q.updateJob(job)
	q.logger.Info("job completed", "job_id", job.ID)
}

func (q *Queue) handleFailure(ctx context.Context, job *Job, err error) {
	job.Error = err.Error()
	job.Retries++

	if job.Retries < job.MaxRetries {
		job.Status = JobStatusRetryable
		q.logger.Warn("job failed, will retry", "job_id", job.ID, "retries", job.Retries, "error", err)
		// Re-enqueue for retry
		time.AfterFunc(q.interval*time.Duration(job.Retries), func() {
			select {
			case q.jobChan <- job:
			default:
			}
		})
	} else {
		now := time.Now()
		job.Status = JobStatusFailed
		job.FailedAt = &now
		q.logger.Error("job failed permanently", "job_id", job.ID, "error", err)
	}

	q.updateJob(job)
}

func (q *Queue) updateJob(job *Job) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.jobs[job.ID] = job
}

// GetJob retrieves a job by ID.
func (q *Queue) GetJob(id string) (*Job, bool) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	job, ok := q.jobs[id]
	return job, ok
}

// ListJobs returns all jobs in the queue.
func (q *Queue) ListJobs() []*Job {
	q.mu.RLock()
	defer q.mu.RUnlock()
	jobs := make([]*Job, 0, len(q.jobs))
	for _, job := range q.jobs {
		jobs = append(jobs, job)
	}
	return jobs
}

// ErrQueueFull is returned when the job queue is full.
var ErrQueueFull = &QueueError{"queue is full"}

type QueueError struct {
	msg string
}

func (e *QueueError) Error() string { return e.msg }

func generateJobID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
		time.Sleep(time.Nanosecond) // ensure different values
	}
	return string(b)
}
