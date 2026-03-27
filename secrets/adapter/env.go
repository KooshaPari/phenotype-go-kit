package adapter

import (
	"context"
	"fmt"
	"os"

	"github.com/KooshaPari/phenotype-go-kit/contracts/ports/outbound"
)

// EnvAdapter implements outbound.SecretPort using environment variables.
// This is intended for development and testing only.
type EnvAdapter struct{}

// NewEnvAdapter creates a new environment variable adapter.
func NewEnvAdapter() *EnvAdapter {
	return &EnvAdapter{}
}

// Get implements outbound.SecretPort.
func (a *EnvAdapter) Get(ctx context.Context, key string) (*outbound.SecretValue, error) {
	value := os.Getenv(key)
	if value == "" {
		return nil, fmt.Errorf("environment variable not set: %s", key)
	}
	return &outbound.SecretValue{
		Key:   key,
		Value: value,
	}, nil
}

// Set implements outbound.SecretPort.
func (a *EnvAdapter) Set(ctx context.Context, key string, value string, opts ...outbound.SecretOption) error {
	if err := os.Setenv(key, value); err != nil {
		return fmt.Errorf("failed to set environment variable: %w", err)
	}
	return nil
}

// Delete implements outbound.SecretPort.
func (a *EnvAdapter) Delete(ctx context.Context, key string) error {
	if err := os.Unsetenv(key); err != nil {
		return fmt.Errorf("failed to unset environment variable: %w", err)
	}
	return nil
}

// List implements outbound.SecretPort.
func (a *EnvAdapter) List(ctx context.Context, path string) ([]string, error) {
	env := os.Environ()
	prefix := path + "_"
	result := make([]string, 0)

	for _, e := range env {
		if len(e) > len(prefix) && e[:len(prefix)] == prefix {
			key := e[:len(prefix)-1]
			result = append(result, key)
		}
	}
	return result, nil
}

// Exists implements outbound.SecretPort.
func (a *EnvAdapter) Exists(ctx context.Context, key string) (bool, error) {
	value := os.Getenv(key)
	return value != "", nil
}
