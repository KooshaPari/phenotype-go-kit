package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// Delivery represents a webhook delivery attempt.
type Delivery struct {
	ID        string          `json:"id"`
	EventType string          `json:"event_type"`
	URL       string          `json:"url"`
	Payload   json.RawMessage `json:"payload"`
	Headers   http.Header     `json:"headers"`

	StatusCode int           `json:"status_code"`
	Response   string        `json:"response,omitempty"`
	Duration   time.Duration `json:"duration"`
	Success    bool          `json:"success"`
	Retries    int           `json:"retries"`
	CreatedAt  time.Time     `json:"created_at"`
	SentAt     *time.Time    `json:"sent_at,omitempty"`
}

// SignatureHeaderName is the header used for webhook signatures.
const SignatureHeaderName = "X-Webhook-Signature"

// DeliveryConfig holds configuration for webhook delivery.
type DeliveryConfig struct {
	Timeout       time.Duration `default:"30s"`
	MaxRetries    int           `default:"5"`
	RetryBackoff  time.Duration `default:"1s"`
	SignAlgorithm string        `default:"hmac-sha256"`
}

// DeliveryService handles webhook deliveries.
type DeliveryService struct {
	client *http.Client
	config DeliveryConfig
	secret []byte
	logger *slog.Logger
}

// NewDeliveryService creates a new webhook delivery service.
func NewDeliveryService(secret string, config DeliveryConfig, logger *slog.Logger) *DeliveryService {
	if config.Timeout <= 0 {
		config.Timeout = 30 * time.Second
	}
	if config.MaxRetries <= 0 {
		config.MaxRetries = 5
	}
	if config.RetryBackoff <= 0 {
		config.RetryBackoff = time.Second
	}

	return &DeliveryService{
		client: &http.Client{
			Timeout: config.Timeout,
		},
		config: config,
		secret: []byte(secret),
		logger: logger,
	}
}

// Deliver sends a webhook to the specified URL.
func (s *DeliveryService) Deliver(ctx context.Context, eventType string, url string, payload interface{}) (*Delivery, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	signature := s.sign(payloadBytes)
	deliveryID := generateDeliveryID()

	delivery := &Delivery{
		ID:        deliveryID,
		EventType: eventType,
		URL:       url,
		Payload:   payloadBytes,
		Headers: http.Header{
			"Content-Type":       []string{"application/json"},
			"X-Webhook-Event":    []string{eventType},
			"X-Webhook-Delivery": []string{deliveryID},
			SignatureHeaderName:  []string{signature},
		},
		CreatedAt: time.Now(),
	}

	return s.sendWithRetry(ctx, delivery)
}

func (s *DeliveryService) sendWithRetry(ctx context.Context, delivery *Delivery) (*Delivery, error) {
	var lastErr error

	for attempt := 0; attempt <= s.config.MaxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			backoff := s.config.RetryBackoff * time.Duration(1<<uint(attempt-1))
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
			delivery.Retries = attempt
		}

		err := s.send(ctx, delivery)
		if err == nil {
			return delivery, nil
		}

		lastErr = err
		s.logger.Warn("webhook delivery failed",
			"attempt", attempt+1,
			"url", delivery.URL,
			"error", err)
	}

	delivery.Success = false
	delivery.Response = lastErr.Error()
	return delivery, lastErr
}

func (s *DeliveryService) send(ctx context.Context, delivery *Delivery) error {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, delivery.URL, bytes.NewReader(delivery.Payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	for k, v := range delivery.Headers {
		req.Header[k] = v
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	delivery.Duration = time.Since(start)
	delivery.StatusCode = resp.StatusCode

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		delivery.Response = fmt.Sprintf("failed to read response: %v", err)
	} else {
		delivery.Response = string(body)
	}

	now := time.Now()
	delivery.SentAt = &now

	// Consider 2xx as success
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		delivery.Success = true
		return nil
	}

	return fmt.Errorf("delivery failed with status %d", resp.StatusCode)
}

// sign generates an HMAC-SHA256 signature for the payload.
func (s *DeliveryService) sign(payload []byte) string {
	mac := hmac.New(sha256.New, s.secret)
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

// VerifySignature verifies that the webhook signature is valid.
func VerifySignature(payload []byte, signature string, secret string) bool {
	expected := computeSignature(payload, secret)
	return hmac.Equal([]byte(expected), []byte(signature))
}

func computeSignature(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

func generateDeliveryID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(12)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[int(time.Now().UnixNano())%len(letters)]
	}
	return string(b)
}
