package outbound

import (
	"context"

	"github.com/KooshaPari/phenotype-go-kit/contracts/models"
)

// Repository defines the interface for data persistence.
// Following DIP: Depend on abstraction, not concrete implementation.
type Repository[T any, ID any] interface {
	// Create persists a new entity and returns the created entity with ID.
	Create(ctx context.Context, entity T) (T, error)

	// GetByID retrieves an entity by its ID.
	GetByID(ctx context.Context, id ID) (T, error)

	// Update modifies an existing entity.
	Update(ctx context.Context, entity T) (T, error)

	// Delete removes an entity by its ID.
	Delete(ctx context.Context, id ID) error

	// List returns all entities with optional filtering.
	List(ctx context.Context, filter *models.QueryFilter) ([]T, error)

	// Exists checks if an entity with the given ID exists.
	Exists(ctx context.Context, id ID) (bool, error)

	// Count returns the total number of entities.
	Count(ctx context.Context, filter *models.QueryFilter) (int64, error)
}

// QueryRepository defines read-only repository operations.
// Following CQRS: Separate read and write models for optimized queries.
type QueryRepository[T any] interface {
	// Find retrieves entities based on query criteria.
	Find(ctx context.Context, criteria *models.QueryCriteria) ([]T, error)

	// FindOne retrieves a single entity based on criteria.
	FindOne(ctx context.Context, criteria *models.QueryCriteria) (*T, error)

	// Aggregate performs aggregation operations (count, sum, avg, etc.).
	Aggregate(ctx context.Context, criteria *models.QueryCriteria) (*models.AggregationResult, error)
}

// EventStore defines interface for event persistence.
// Following Event Sourcing: Store events, derive state.
type EventStore interface {
	// Append adds events to the store.
	Append(ctx context.Context, aggregateID string, events []models.DomainEvent) error

	// GetEvents retrieves all events for an aggregate.
	GetEvents(ctx context.Context, aggregateID string) ([]models.DomainEvent, error)

	// GetEventsSince retrieves events since a specific version.
	GetEventsSince(ctx context.Context, aggregateID string, version int64) ([]models.DomainEvent, error)

	// GetAllEvents retrieves all events with optional filtering.
	GetAllEvents(ctx context.Context, filter *models.EventFilter) ([]models.DomainEvent, error)
}

// EventPublisher defines interface for publishing domain events.
type EventPublisher interface {
	// Publish sends events to the message bus.
	Publish(ctx context.Context, topic string, events []models.DomainEvent) error

	// PublishAsync sends events asynchronously.
	PublishAsync(ctx context.Context, topic string, events []models.DomainEvent) (<-chan error, error)

	// Subscribe registers a handler for events on a topic.
	Subscribe(ctx context.Context, topic string, handler models.EventHandler) error
}

// EventBusPort defines the minimal event bus contract used by application services.
type EventBusPort interface {
	Publish(ctx context.Context, topic string, event any) error
}

// Cache defines interface for caching operations.
type Cache interface {
	// Get retrieves a value from cache.
	Get(ctx context.Context, key string) ([]byte, error)

	// Set stores a value in cache with optional TTL.
	Set(ctx context.Context, key string, value []byte, ttl *models.Duration) error

	// Delete removes a key from cache.
	Delete(ctx context.Context, key string) error

	// Exists checks if a key exists in cache.
	Exists(ctx context.Context, key string) (bool, error)

	// Clear removes all entries from cache.
	Clear(ctx context.Context) error

	// GetOrSet retrieves from cache or compute and store.
	GetOrSet(ctx context.Context, key string, compute func() ([]byte, error), ttl *models.Duration) ([]byte, error)
}

// ExternalService defines interface for calling external HTTP services.
type ExternalService interface {
	// Call makes an HTTP request to an external service.
	Call(ctx context.Context, request *models.ExternalRequest) (*models.ExternalResponse, error)

	// CallWithRetry makes an HTTP request with retry logic.
	CallWithRetry(ctx context.Context, request *models.ExternalRequest, retryConfig *models.RetryConfig) (*models.ExternalResponse, error)
}

// SecretStore defines interface for secret management.
type SecretStore interface {
	// Get retrieves a secret value.
	Get(ctx context.Context, key string) (string, error)

	// Set stores a secret value.
	Set(ctx context.Context, key string, value string) error

	// Delete removes a secret.
	Delete(ctx context.Context, key string) error

	// List returns all secret keys.
	List(ctx context.Context, path string) ([]string, error)
}

// MetricsCollector defines interface for collecting metrics.
type MetricsCollector interface {
	// Counter records a counter metric.
	Counter(ctx context.Context, name string, value float64, labels map[string]string) error

	// Gauge records a gauge metric.
	Gauge(ctx context.Context, name string, value float64, labels map[string]string) error

	// Histogram records a histogram metric.
	Histogram(ctx context.Context, name string, value float64, labels map[string]string) error

	// Summary records a summary metric.
	Summary(ctx context.Context, name string, value float64, labels map[string]string) error
}

// Logger defines interface for logging operations.
type Logger interface {
	// Debug logs a debug message.
	Debug(ctx context.Context, msg string, args ...any)

	// Info logs an info message.
	Info(ctx context.Context, msg string, args ...any)

	// Warn logs a warning message.
	Warn(ctx context.Context, msg string, args ...any)

	// Error logs an error message.
	Error(ctx context.Context, msg string, args ...any)

	// Fatal logs a fatal message and exits.
	Fatal(ctx context.Context, msg string, args ...any)
}

// ObservabilityPort records errors and operational signals.
type ObservabilityPort interface {
	RecordError(ctx context.Context, operation string, err error)
}

// ConfigProvider defines interface for configuration management.
type ConfigProvider interface {
	// Get retrieves a configuration value.
	Get(ctx context.Context, key string) (any, error)

	// Set sets a configuration value.
	Set(ctx context.Context, key string, value any) error

	// GetAll retrieves all configuration.
	GetAll(ctx context.Context) (map[string]any, error)

	// Subscribe registers a callback for configuration changes.
	Subscribe(ctx context.Context, key string, callback func(any)) error
}

// AdapterManifest provides metadata about an adapter implementation.
// Following PoLA: Adapters declare their capabilities explicitly.
type AdapterManifest struct {
	// Name is the unique name of the adapter.
	Name string

	// Version is the semantic version of the adapter.
	Version string

	// Description describes what the adapter does.
	Description string

	// Provides lists the port interfaces this adapter implements.
	Provides []string
}

// Adapter is implemented by all adapters to provide self-description.
// Following SoC: Adapters self-document their capabilities.
type Adapter interface {
	// Manifest returns the adapter's manifest.
	Manifest() *AdapterManifest
}
