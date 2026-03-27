package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/KooshaPari/phenotype-go-kit/contracts/ports/outbound"
)

// AWSSecretsAdapter implements outbound.SecretPort for AWS Secrets Manager.
type AWSSecretsAdapter struct {
	region string
}

// NewAWSSecretsAdapter creates a new AWS Secrets Manager adapter.
func NewAWSSecretsAdapter(region string) *AWSSecretsAdapter {
	return &AWSSecretsAdapter{
		region: region,
	}
}

// Get implements outbound.SecretPort.
func (a *AWSSecretsAdapter) Get(ctx context.Context, key string) (*outbound.SecretValue, error) {
	cmd := exec.CommandContext(ctx, "aws", "secretsmanager", "get-secret-value",
		"--secret-id", key,
		"--region", a.region,
		"--query", "SecretString",
		"--output", "text",
	)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("aws get-secret-value failed: %w", err)
	}

	return &outbound.SecretValue{
		Key:   key,
		Value: strings.TrimSpace(string(output)),
	}, nil
}

// Set implements outbound.SecretPort.
func (a *AWSSecretsAdapter) Set(ctx context.Context, key string, value string, opts ...outbound.SecretOption) error {
	cmd := exec.CommandContext(ctx, "aws", "secretsmanager", "put-secret-value",
		"--secret-id", key,
		"--secret-string", value,
		"--region", a.region,
	)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("aws put-secret-value failed: %w", err)
	}
	return nil
}

// Delete implements outbound.SecretPort.
func (a *AWSSecretsAdapter) Delete(ctx context.Context, key string) error {
	cmd := exec.CommandContext(ctx, "aws", "secretsmanager", "delete-secret",
		"--secret-id", key,
		"--region", a.region,
		"--force-delete-without-recovery",
	)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("aws delete-secret failed: %w", err)
	}
	return nil
}

// List implements outbound.SecretPort.
func (a *AWSSecretsAdapter) List(ctx context.Context, path string) ([]string, error) {
	cmd := exec.CommandContext(ctx, "aws", "secretsmanager", "list-secrets",
		"--filter", "key=NAME,values="+path,
		"--region", a.region,
		"--output", "json",
	)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("aws list-secrets failed: %w", err)
	}

	var result struct {
		SecretList []struct {
			Name string `json:"Name"`
		} `json:"SecretList"`
	}
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse list output: %w", err)
	}

	names := make([]string, len(result.SecretList))
	for i, s := range result.SecretList {
		names[i] = s.Name
	}
	return names, nil
}

// Exists implements outbound.SecretPort.
func (a *AWSSecretsAdapter) Exists(ctx context.Context, key string) (bool, error) {
	_, err := a.Get(ctx, key)
	if err != nil {
		return false, nil
	}
	return true, nil
}
