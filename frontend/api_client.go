package frontend

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

// Client provides HTTP API client functionality.
type Client struct {
	baseURL    string
	httpClient *http.Client
	headers    map[string]string
	logger     *slog.Logger
}

// Config holds API client configuration.
type Config struct {
	BaseURL       string
	Timeout       time.Duration
	RetryAttempts int
	RetryDelay    time.Duration
	APIKey        string
}

// NewClient creates a new API client.
func NewClient(cfg Config) *Client {
	client := &http.Client{
		Timeout: cfg.Timeout,
	}

	if cfg.Timeout == 0 {
		client.Timeout = 30 * time.Second
	}

	return &Client{
		baseURL:    cfg.BaseURL,
		httpClient: client,
		headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		logger: slog.Default(),
	}
}

// SetAuth sets authorization header.
func (c *Client) SetAuth(token string) {
	c.headers["Authorization"] = "Bearer " + token
}

// SetAPIKey sets API key header.
func (c *Client) SetAPIKey(key string) {
	c.headers["X-API-Key"] = key
}

// SetHeader sets a custom header.
func (c *Client) SetHeader(key, value string) {
	c.headers[key] = value
}

// Get performs a GET request.
func (c *Client) Get(ctx context.Context, path string, params map[string]string) (*Response, error) {
	return c.doRequest(ctx, "GET", path, params, nil)
}

// Post performs a POST request.
func (c *Client) Post(ctx context.Context, path string, body interface{}) (*Response, error) {
	return c.doRequest(ctx, "POST", path, nil, body)
}

// Put performs a PUT request.
func (c *Client) Put(ctx context.Context, path string, body interface{}) (*Response, error) {
	return c.doRequest(ctx, "PUT", path, nil, body)
}

// Patch performs a PATCH request.
func (c *Client) Patch(ctx context.Context, path string, body interface{}) (*Response, error) {
	return c.doRequest(ctx, "PATCH", path, nil, body)
}

// Delete performs a DELETE request.
func (c *Client) Delete(ctx context.Context, path string) (*Response, error) {
	return c.doRequest(ctx, "DELETE", path, nil, nil)
}

func (c *Client) doRequest(ctx context.Context, method, path string, params map[string]string, body interface{}) (*Response, error) {
	uri, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, err
	}

	if params != nil {
		query := uri.Query()
		for k, v := range params {
			query.Add(k, v)
		}
		uri.RawQuery = query.Encode()
	}

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, uri.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := &Response{
		StatusCode: resp.StatusCode,
		Body:       respBody,
		Headers:    resp.Header,
	}

	if ct := resp.Header.Get("Content-Type"); ct == "application/json" {
		_ = json.Unmarshal(respBody, &response.Data)
	}

	return response, nil
}

// Response holds the API response.
type Response struct {
	StatusCode int
	Body       []byte
	Data       interface{}
	Headers    http.Header
}

// Error checks if response is an error.
func (r *Response) Error() error {
	if r.StatusCode >= 400 {
		return fmt.Errorf("API error: %d - %s", r.StatusCode, string(r.Body))
	}
	return nil
}

// IsSuccess checks if status code is success.
func (r *Response) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}
