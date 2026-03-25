// Package ports defines the interfaces (ports) for the hexagonal architecture.
// These ports define the contracts that adapters must implement.
package ports

import (
	"context"
)

// Repository defines the interface for data persistence.
type Repository interface {
	// Create inserts a new entity.
	Create(ctx context.Context, entity interface{}) error
	// GetByID retrieves an entity by its ID.
	GetByID(ctx context.Context, id string) (interface{}, error)
	// Update modifies an existing entity.
	Update(ctx context.Context, entity interface{}) error
	// Delete removes an entity by its ID.
	Delete(ctx context.Context, id string) error
	// List returns all entities with optional pagination.
	List(ctx context.Context, limit, offset int) ([]interface{}, error)
}

// Cache defines the interface for caching operations.
type Cache interface {
	// Get retrieves a value by key.
	Get(ctx context.Context, key string) ([]byte, error)
	// Set stores a value with optional TTL.
	Set(ctx context.Context, key string, value []byte, ttl int) error
	// Delete removes a key.
	Delete(ctx context.Context, key string) error
	// Invalidate removes keys matching a pattern.
	Invalidate(ctx context.Context, pattern string) error
}

// EventBus defines the interface for publishing and subscribing to events.
type EventBus interface {
	// Publish sends an event to all subscribers.
	Publish(ctx context.Context, event interface{}) error
	// Subscribe registers a handler for events matching the filter.
	Subscribe(handler EventHandler, filters ...EventFilter) error
}

// EventHandler is the function signature for event handlers.
type EventHandler func(ctx context.Context, event interface{}) error

// EventFilter filters events before reaching the handler.
type EventFilter func(event interface{}) bool

// MessageQueue defines the interface for reliable message delivery.
type MessageQueue interface {
	// Enqueue adds a message to the queue.
	Enqueue(ctx context.Context, queue string, message interface{}) error
	// Dequeue retrieves and removes a message from the queue.
	Dequeue(ctx context.Context, queue string) (interface{}, error)
	// Acknowledge confirms successful processing.
	Acknowledge(ctx context.Context, messageID string) error
	// Reject returns a message to the queue.
	Reject(ctx context.Context, messageID string, requeue bool) error
}

// HTTPClient defines the interface for making HTTP requests.
type HTTPClient interface {
	// Do performs an HTTP request.
	Do(ctx context.Context, req *HTTPRequest) (*HTTPResponse, error)
}

// HTTPRequest represents an HTTP request.
type HTTPRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte
}

// HTTPResponse represents an HTTP response.
type HTTPResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}
