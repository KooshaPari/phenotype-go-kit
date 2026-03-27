package embeddings

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/KooshaPari/phenotype-go-kit/contracts/plugins"
)

// OllamaProvider implements Provider for Ollama API.
type OllamaProvider struct {
	config   Config
	client   *http.Client
	endpoint string
	model    string
}

// Ensure OllamaProvider implements Provider.
var _ Provider = (*OllamaProvider)(nil)

// OllamaConfig extends Config with Ollama-specific options.
type OllamaConfig struct {
	Config
	KeepAlive string // How long to keep model loaded (e.g., "5m", "0")
}

// NewOllamaProvider creates a new Ollama embeddings provider.
func NewOllamaProvider(opts ...Option) *OllamaProvider {
	cfg := Config{
		BaseURL:    "http://localhost:11434",
		Model:      "nomic-embed-text",
		Dimensions: 768,
		Timeout:    60 * time.Second,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	return &OllamaProvider{
		config:   cfg,
		client:   &http.Client{Timeout: cfg.Timeout},
		endpoint: cfg.BaseURL + "/api/embeddings",
		model:    cfg.Model,
	}
}

// Manifest returns the plugin manifest.
func (p *OllamaProvider) Manifest() *plugins.Manifest {
	return plugins.DefaultManifest()
}

// Name returns the provider name.
func (p *OllamaProvider) Name() string {
	return "ollama"
}

// Model returns the model name.
func (p *OllamaProvider) Model() string {
	return p.model
}

// Dimensions returns the embedding dimensions.
func (p *OllamaProvider) Dimensions() int {
	return p.config.Dimensions
}

// OllamaEmbedRequest represents the Ollama embeddings API request.
type OllamaEmbedRequest struct {
	Model     string `json:"model"`
	Prompt    string `json:"prompt"`
	KeepAlive string `json:"keep_alive,omitempty"`
}

// OllamaEmbedResponse represents the Ollama embeddings API response.
type OllamaEmbedResponse struct {
	Embedding []float32 `json:"embedding"`
}

// Embed generates embeddings for the given texts.
func (p *OllamaProvider) Embed(ctx context.Context, texts []string) (*EmbedResult, error) {
	embeddings := make([]Embedding, len(texts))

	for i, text := range texts {
		reqBody := OllamaEmbedRequest{
			Model:     p.model,
			Prompt:    text,
			KeepAlive: "5m",
		}

		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request: %w", err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.endpoint, bytes.NewReader(jsonBody))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := p.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to send request: %w", err)
		}

		var embedResp OllamaEmbedResponse
		if err := json.NewDecoder(resp.Body).Decode(&embedResp); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		resp.Body.Close()

		embeddings[i] = Embedding{
			Vector:  embedResp.Embedding,
			Model:   p.model,
			Tokens:  0,
			Created: time.Now(),
		}
	}

	return &EmbedResult{
		Embeddings: embeddings,
		Usage:      Usage{},
	}, nil
}

// EmbedSingle generates an embedding for a single text.
func (p *OllamaProvider) EmbedSingle(ctx context.Context, text string) (*Embedding, error) {
	result, err := p.Embed(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	if len(result.Embeddings) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}
	return &result.Embeddings[0], nil
}
