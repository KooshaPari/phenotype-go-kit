package models

import (
	"encoding/json"
	"time"
)

// DomainEvent represents a domain event following Event Sourcing patterns.
type DomainEvent struct {
	// ID is the unique identifier of the event.
	ID string `json:"id"`

	// AggregateID is the ID of the aggregate that generated this event.
	AggregateID string `json:"aggregate_id"`

	// AggregateType is the type of the aggregate.
	AggregateType string `json:"aggregate_type"`

	// EventType is the type of the event.
	EventType string `json:"event_type"`

	// EventData contains the event payload.
	EventData json.RawMessage `json:"event_data"`

	// Version is the version of the aggregate at the time of the event.
	Version int64 `json:"version"`

	// OccurredAt is when the event occurred.
	OccurredAt time.Time `json:"occurred_at"`

	// Metadata contains additional event metadata.
	Metadata map[string]string `json:"metadata,omitempty"`
}

// EventHandler is a function type for handling events.
type EventHandler func(ctx interface{}, event *DomainEvent) error

// EventFilter contains filtering criteria for event queries.
type EventFilter struct {
	// AggregateID filters by aggregate ID.
	AggregateID string

	// AggregateType filters by aggregate type.
	AggregateType string

	// EventType filters by event type.
	EventType string

	// From filters events after this time.
	From *time.Time

	// To filters events before this time.
	To *time.Time

	// Limit limits the number of events returned.
	Limit int

	// Offset offsets the results.
	Offset int
}

// CommandResult represents the result of a command execution.
type CommandResult struct {
	// Success indicates if the command succeeded.
	Success bool `json:"success"`

	// Error contains error information if the command failed.
	Error *CommandError `json:"error,omitempty"`

	// Data contains command result data.
	Data any `json:"data,omitempty"`

	// Metadata contains additional result metadata.
	Metadata map[string]string `json:"metadata,omitempty"`
}

// CommandError represents a command execution error.
type CommandError struct {
	// Code is the error code.
	Code string `json:"code"`

	// Message is the human-readable error message.
	Message string `json:"message"`

	// Details contains additional error details.
	Details map[string]any `json:"details,omitempty"`
}

// QueryFilter contains filtering criteria for entity queries.
type QueryFilter struct {
	// IDs filters by entity IDs.
	IDs []string

	// Limit limits the number of results.
	Limit int

	// Offset offsets the results.
	Offset int

	// SortBy specifies the field to sort by.
	SortBy string

	// SortOrder specifies the sort order (asc/desc).
	SortOrder string
}

// QueryCriteria contains advanced query criteria.
type QueryCriteria struct {
	// Filter contains field-value filters.
	Filter map[string]any

	// Pagination for the query.
	Limit  int
	Offset int

	// Sort specifies sorting criteria.
	Sort []SortCriterion

	// Include specifies related entities to include.
	Include []string

	// Exclude specifies fields to exclude.
	Exclude []string
}

// SortCriterion represents a sorting criterion.
type SortCriterion struct {
	// Field is the field to sort by.
	Field string

	// Order is the sort order (asc/desc).
	Order string
}

// AggregationResult represents the result of an aggregation query.
type AggregationResult struct {
	// Count is the count of matching entities.
	Count int64 `json:"count,omitempty"`

	// Sum is the sum of a numeric field.
	Sum float64 `json:"sum,omitempty"`

	// Average is the average of a numeric field.
	Average float64 `json:"average,omitempty"`

	// Min is the minimum value.
	Min float64 `json:"min,omitempty"`

	// Max is the maximum value.
	Max float64 `json:"max,omitempty"`

	// Groups contains grouped aggregation results.
	Groups []AggregationGroup `json:"groups,omitempty"`
}

// AggregationGroup represents a group in aggregation results.
type AggregationGroup struct {
	// Key is the group key value.
	Key any `json:"key"`

	// Count is the count of items in the group.
	Count int64 `json:"count"`

	// Sum is the sum of values in the group.
	Sum float64 `json:"sum,omitempty"`

	// Average is the average of values in the group.
	Average float64 `json:"average,omitempty"`
}

// Duration represents a time duration with JSON support.
type Duration struct {
	// Duration is the duration value.
	time.Duration
}

// MarshalJSON serializes the duration to JSON.
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// UnmarshalJSON deserializes the duration from JSON.
func (d *Duration) UnmarshalJSON(data []byte) error {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch val := v.(type) {
	case string:
		p, err := time.ParseDuration(val)
		if err != nil {
			return err
		}
		d.Duration = p
	case float64:
		d.Duration = time.Duration(val)
	case int64:
		d.Duration = time.Duration(val)
	default:
		return nil
	}
	return nil
}

// ExternalRequest represents an HTTP request to an external service.
type ExternalRequest struct {
	// Method is the HTTP method.
	Method string

	// URL is the request URL.
	URL string

	// Headers are the request headers.
	Headers map[string]string

	// Body is the request body.
	Body []byte

	// Timeout is the request timeout.
	Timeout time.Duration

	// RetryConfig contains retry configuration.
	RetryConfig *RetryConfig
}

// ExternalResponse represents an HTTP response from an external service.
type ExternalResponse struct {
	// StatusCode is the HTTP status code.
	StatusCode int

	// Headers are the response headers.
	Headers map[string]string

	// Body is the response body.
	Body []byte

	// Error contains error information if the request failed.
	Error string
}

// RetryConfig contains configuration for retry logic.
type RetryConfig struct {
	// MaxAttempts is the maximum number of retry attempts.
	MaxAttempts int

	// InitialInterval is the initial retry interval.
	InitialInterval time.Duration

	// MaxInterval is the maximum retry interval.
	MaxInterval time.Duration

	// Multiplier is the backoff multiplier.
	Multiplier float64

	// Jitter adds randomness to retry intervals.
	Jitter bool
}

// NewDefaultRetryConfig returns a default retry configuration.
func NewDefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxAttempts:     3,
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     30 * time.Second,
		Multiplier:      2.0,
		Jitter:          true,
	}
}
