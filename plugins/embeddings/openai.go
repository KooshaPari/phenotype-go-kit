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

// OpenAIProvider implements Provider for OpenAI API.
type OpenAIProvider struct {
	config   Config
	client   *http.Client
	endpoint string
}

// Ensure OpenAIProvider implements Provider.
var _ Provider = (*OpenAIProvider)(nil)

// NewOpenAIProvider creates a new OpenAI embeddings provider.
func NewOpenAIProvider(opts ...Option) *OpenAIProvider {
	cfg := Config{
		BaseURL:    "https://api.openai.com/v1",
		Model:      "text-embedding-3-small",
		Dimensions: 1536,
		Timeout:    30 * time.Second,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	return &OpenAIProvider{
		config:   cfg,
		client:   &http.Client{Timeout: cfg.Timeout},
		endpoint: cfg.BaseURL + "/embeddings",
	}
}

// Manifest returns the plugin manifest.
func (p *OpenAIProvider) Manifest() *plugins.Manifest {
	return plugins.DefaultManifest()
}

// Name returns the provider name.
func (p *OpenAIProvider) Name() string {
	return "openai"
}

// Model returns the model name.
func (p *OpenAIProvider) Model() string {
	return p.config.Model
}

// Dimensions returns the embedding dimensions.
func (p *OpenAIProvider) Dimensions() int {
	return p.config.Dimensions
}

// EmbedRequest represents the OpenAI embeddings API request.
type EmbedRequest struct {
	Model     string   `json:"model"`
	Input     []string `json:"input"`
	Dimension int      `json:"dimensions,omitempty"`
}

// EmbedResponse represents the OpenAI embeddings API response.
type EmbedResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Embedding []float32 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// Embed generates embeddings for the given texts.
func (p *OpenAIProvider) Embed(ctx context.Context, texts []string) (*EmbedResult, error) {
	reqBody := EmbedRequest{
		Model:     p.config.Model,
		Input:     texts,
		Dimension: p.config.Dimensions,
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
	req.Header.Set("Authorization", "Bearer "+p.config.APIKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	var embedResp EmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&embedResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	result := &EmbedResult{
		Embeddings: make([]Embedding, len(embedResp.Data)),
		Usage: Usage{
			PromptTokens:     embedResp.Usage.PromptTokens,
			CompletionTokens: embedResp.Usage.CompletionTokens,
			TotalTokens:      embedResp.Usage.TotalTokens,
		},
	}

	for i, data := range embedResp.Data {
		result.Embeddings[i] = Embedding{
			Vector:  data.Embedding,
			Model:   p.config.Model,
			Tokens:  0, // Not provided per-embedding
			Created: time.Now(),
		}
	}

	return result, nil
}

// EmbedSingle generates an embedding for a single text.
func (p *OpenAIProvider) EmbedSingle(ctx context.Context, text string) (*Embedding, error) {
	result, err := p.Embed(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	if len(result.Embeddings) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}
	return &result.Embeddings[0], nil
}
