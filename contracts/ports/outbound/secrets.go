package outbound

import "context"

// SecretValue represents a secret value with metadata.
type SecretValue struct {
	Key      string
	Value    string
	Version  int
	Metadata map[string]string
}

// SecretPort defines the interface for secret management operations.
// This is an outbound port implemented by infrastructure adapters.
type SecretPort interface {
	// Get retrieves a secret by key.
	Get(ctx context.Context, key string) (*SecretValue, error)

	// Set stores a secret.
	Set(ctx context.Context, key string, value string, opts ...SecretOption) error

	// Delete removes a secret.
	Delete(ctx context.Context, key string) error

	// List returns all secret keys (without values).
	List(ctx context.Context, path string) ([]string, error)

	// Exists checks if a secret exists.
	Exists(ctx context.Context, key string) (bool, error)
}

// SecretOption configures secret operations.
type SecretOption func(*secretOptions)

type secretOptions struct {
	version       int
	metadata      map[string]string
	leaseDuration int // seconds
	rotate        bool
}

// WithVersion specifies the secret version.
func WithVersion(version int) SecretOption {
	return func(o *secretOptions) {
		o.version = version
	}
}

// WithMetadata adds metadata to the secret.
func WithMetadata(metadata map[string]string) SecretOption {
	return func(o *secretOptions) {
		o.metadata = metadata
	}
}

// WithLeaseDuration sets the lease duration in seconds.
func WithLeaseDuration(seconds int) SecretOption {
	return func(o *secretOptions) {
		o.leaseDuration = seconds
	}
}

// WithRotation enables automatic rotation.
func WithRotation(enabled bool) SecretOption {
	return func(o *secretOptions) {
		o.rotate = enabled
	}
}

// SecretReader defines read-only secret access.
type SecretReader interface {
	Get(ctx context.Context, key string) (*SecretValue, error)
	List(ctx context.Context, path string) ([]string, error)
	Exists(ctx context.Context, key string) (bool, error)
}

// SecretWriter defines write-only secret access.
type SecretWriter interface {
	Set(ctx context.Context, key string, value string, opts ...SecretOption) error
	Delete(ctx context.Context, key string) error
}

// SecretManager combines read and write operations.
type SecretManager interface {
	SecretReader
	SecretWriter
}

// SecretEncryptor defines encryption/decryption operations.
type SecretEncryptor interface {
	Encrypt(ctx context.Context, plaintext []byte) ([]byte, error)
	Decrypt(ctx context.Context, ciphertext []byte) ([]byte, error)
}

// SecretRotator defines secret rotation capabilities.
type SecretRotator interface {
	Rotate(ctx context.Context, key string) error
	GetRotationSchedule(ctx context.Context, key string) (string, error)
	SetRotationSchedule(ctx context.Context, key string, schedule string) error
}
