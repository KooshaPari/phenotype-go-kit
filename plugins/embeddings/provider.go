package embeddings

import (
	"context"
	"time"
)

// Embedding represents a vector embedding.
type Embedding struct {
	Vector  []float32
	Model   string
	Tokens  int
	Created time.Time
}

// EmbedResult contains the result of an embedding operation.
type EmbedResult struct {
	Embeddings []Embedding
	Usage      Usage
}

// Usage represents token usage statistics.
type Usage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

// Provider defines the interface for embeddings providers.
// Following Interface Segregation Principle (ISP) - minimal interface.
//
// # Law of Demeter Compliance
//
// Implementations must:
//   - Accept dependencies via constructor (not service locator)
//   - Return value objects (not internal state)
//   - Not expose internal components via getters
type Provider interface {
	// Embed generates embeddings for the given texts.
	Embed(ctx context.Context, texts []string) (*EmbedResult, error)

	// EmbedSingle generates an embedding for a single text.
	EmbedSingle(ctx context.Context, text string) (*Embedding, error)

	// Name returns the provider name.
	Name() string

	// Model returns the default model name.
	Model() string

	// Dimensions returns the embedding dimensions.
	Dimensions() int
}

// Config holds common configuration for embeddings providers.
type Config struct {
	APIKey     string
	BaseURL    string
	Model      string
	Dimensions int
	Timeout    time.Duration
}

// Option defines functional options for provider configuration.
type Option func(*Config)

// WithAPIKey sets the API key.
func WithAPIKey(key string) Option {
	return func(c *Config) {
		c.APIKey = key
	}
}

// WithBaseURL sets the base URL.
func WithBaseURL(url string) Option {
	return func(c *Config) {
		c.BaseURL = url
	}
}

// WithModel sets the model name.
func WithModel(model string) Option {
	return func(c *Config) {
		c.Model = model
	}
}

// WithDimensions sets the embedding dimensions.
func WithDimensions(dims int) Option {
	return func(c *Config) {
		c.Dimensions = dims
	}
}

// WithTimeout sets the request timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}
