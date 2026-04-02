package secrets

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
)

// Manager handles secret management.
type Manager struct {
	provider string
	logger   *slog.Logger
}

// NewManager creates a new secrets manager.
func NewManager(provider string) *Manager {
	return &Manager{
		provider: provider,
		logger:   slog.Default(),
	}
}

// Get retrieves a secret.
func (m *Manager) Get(ctx context.Context, key string) (string, error) {
	switch m.provider {
	case "vault":
		return m.getFromVault(ctx, key)
	case "aws":
		return m.getFromAWS(ctx, key)
	case "env":
		return os.Getenv(key), nil
	default:
		return "", fmt.Errorf("unknown provider: %s", m.provider)
	}
}

func (m *Manager) getFromVault(ctx context.Context, key string) (string, error) {
	cmd := exec.CommandContext(ctx, "vault", "kv", "get", "-field=value", "secret/data/"+key)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (m *Manager) getFromAWS(ctx context.Context, key string) (string, error) {
	cmd := exec.CommandContext(ctx, "aws", "secretsmanager", "get-secret-value", "--secret-id", key, "--query", "SecretString", "--output", "text")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// Set stores a secret.
func (m *Manager) Set(ctx context.Context, key, value string) error {
	switch m.provider {
	case "vault":
		return m.setInVault(ctx, key, value)
	case "aws":
		return m.setInAWS(ctx, key, value)
	case "env":
		os.Setenv(key, value)
		return nil
	default:
		return fmt.Errorf("unknown provider: %s", m.provider)
	}
}

func (m *Manager) setInVault(ctx context.Context, key, value string) error {
	cmd := exec.CommandContext(ctx, "vault", "kv", "put", "secret/data/"+key, "value="+value)
	return cmd.Run()
}

func (m *Manager) setInAWS(ctx context.Context, key, value string) error {
	cmd := exec.CommandContext(ctx, "aws", "secretsmanager", "put-secret-value", "--secret-id", key, "--secret-string", value)
	return cmd.Run()
}

// Delete removes a secret.
func (m *Manager) Delete(ctx context.Context, key string) error {
	switch m.provider {
	case "vault":
		cmd := exec.CommandContext(ctx, "vault", "kv", "delete", "secret/data/"+key)
		return cmd.Run()
	case "aws":
		cmd := exec.CommandContext(ctx, "aws", "secretsmanager", "delete-secret", "--secret-id", key, "--force-delete-without-recovery")
		return cmd.Run()
	default:
		return fmt.Errorf("unknown provider: %s", m.provider)
	}
}

// EnvFromSecrets generates environment variables from secrets.
func (m *Manager) EnvFromSecrets(ctx context.Context, keys []string) (map[string]string, error) {
	result := make(map[string]string)

	for _, key := range keys {
		value, err := m.Get(ctx, key)
		if err != nil {
			m.logger.Warn("failed to get secret", "key", key, "error", err)
			continue
		}

		envKey := formatEnvKey(key)
		result[envKey] = value
	}

	return result, nil
}

func formatEnvKey(key string) string {
	return key
}

// Base64Encode encodes a value to base64.
func Base64Encode(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

// Base64Decode decodes a base64 value.
func Base64Decode(encoded string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
