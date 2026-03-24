package webhook

import (
	"crypto/hmac"
	"errors"
	"net/http"
	"strings"
)

// Signature verification errors.
var (
	ErrInvalidSignature = errors.New("invalid webhook signature")
	ErrMissingSignature = errors.New("missing webhook signature header")
)

// VerifyRequest verifies the signature of an incoming webhook request.
func VerifyRequest(req *http.Request, payload []byte, secret string) error {
	signature := req.Header.Get(SignatureHeaderName)
	if signature == "" {
		return ErrMissingSignature
	}

	// Support multiple algorithms
	if strings.HasPrefix(signature, "sha256=") {
		signature = strings.TrimPrefix(signature, "sha256=")
	}

	expected := computeSignature(payload, secret)
	if !hmac.Equal([]byte(signature), []byte(expected)) {
		return ErrInvalidSignature
	}

	return nil
}

// VerifyPayload verifies a raw payload against a signature.
func VerifyPayload(payload []byte, signature string, secret string) error {
	if signature == "" {
		return ErrMissingSignature
	}

	if strings.HasPrefix(signature, "sha256=") {
		signature = strings.TrimPrefix(signature, "sha256=")
	}

	if !VerifySignature(payload, signature, secret) {
		return ErrInvalidSignature
	}

	return nil
}

// GetDeliveryID extracts the delivery ID from the request header.
func GetDeliveryID(req *http.Request) string {
	return req.Header.Get("X-Webhook-Delivery")
}

// GetEventType extracts the event type from the request header.
func GetEventType(req *http.Request) string {
	return req.Header.Get("X-Webhook-Event")
}
