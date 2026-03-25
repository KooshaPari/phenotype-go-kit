package embeddings

import (
	"context"
	"fmt"
	"sync"
)

// Registry manages embeddings providers.
// Following Law of Demeter - only exposes necessary operations.
type Registry struct {
	providers map[string]Provider
	mu        sync.RWMutex
}

// NewRegistry creates a new embeddings registry.
func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]Provider),
	}
}

// Register adds a provider to the registry.
func (r *Registry) Register(provider Provider) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := provider.Name()
	if _, exists := r.providers[name]; exists {
		return fmt.Errorf("provider %q already registered", name)
	}

	r.providers[name] = provider
	return nil
}

// Unregister removes a provider from the registry.
func (r *Registry) Unregister(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.providers, name)
}

// Get retrieves a provider by name.
func (r *Registry) Get(name string) (Provider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	provider, exists := r.providers[name]
	if !exists {
		return nil, fmt.Errorf("provider %q not found", name)
	}

	return provider, nil
}

// List returns all registered provider names.
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}

	return names
}

// Embed delegates to a named provider.
func (r *Registry) Embed(ctx context.Context, providerName string, texts []string) (*EmbedResult, error) {
	provider, err := r.Get(providerName)
	if err != nil {
		return nil, err
	}

	return provider.Embed(ctx, texts)
}

// EmbedSingle delegates to a named provider for single text.
func (r *Registry) EmbedSingle(ctx context.Context, providerName string, text string) (*Embedding, error) {
	provider, err := r.Get(providerName)
	if err != nil {
		return nil, err
	}

	return provider.EmbedSingle(ctx, text)
}

// DefaultRegistry returns a registry pre-configured with default providers.
func DefaultRegistry() *Registry {
	r := NewRegistry()

	// Register OpenAI as default
	openai := NewOpenAIProvider()
	_ = r.Register(openai)

	return r
}
