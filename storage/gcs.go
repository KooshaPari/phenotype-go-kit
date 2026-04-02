package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// GCSConfig holds configuration for Google Cloud Storage.
type GCSConfig struct {
	Bucket      string
	ProjectID   string
	Credentials string // Path to JSON credentials file
}

// GCSStorage implements ObjectStorage using Google Cloud Storage.
type GCSStorage struct {
	client *storage.Client
	bucket string
}

// NewGCSStorage creates a new GCS storage client.
func NewGCSStorage(ctx context.Context, cfg GCSConfig) (*GCSStorage, error) {
	var opts []option.ClientOption
	if cfg.Credentials != "" {
		opts = append(opts, option.WithCredentialsFile(cfg.Credentials))
	}

	client, err := storage.NewClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS client: %w", err)
	}

	return &GCSStorage{
		client: client,
		bucket: cfg.Bucket,
	}, nil
}

// Upload uploads a file to GCS.
func (s *GCSStorage) Upload(ctx context.Context, key string, data []byte, contentType string) error {
	bucket := s.client.Bucket(s.bucket)
	obj := bucket.Object(key)

	writer := obj.NewWriter(ctx)
	writer.ContentType = contentType

	if _, err := writer.Write(data); err != nil {
		writer.Close()
		return fmt.Errorf("failed to write to GCS: %w", err)
	}

	return writer.Close()
}

// Download downloads a file from GCS.
func (s *GCSStorage) Download(ctx context.Context, key string) ([]byte, error) {
	bucket := s.client.Bucket(s.bucket)
	obj := bucket.Object(key)

	reader, err := obj.NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return io.ReadAll(reader)
}

// Delete deletes a file from GCS.
func (s *GCSStorage) Delete(ctx context.Context, key string) error {
	bucket := s.client.Bucket(s.bucket)
	obj := bucket.Object(key)
	return obj.Delete(ctx)
}

// GetURL returns a signed URL for the file.
func (s *GCSStorage) GetURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	// Note: SignedURL requires service account credentials
	// In production, use proper signing with googleAccessID
	// For now, return a simple URL format
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", s.bucket, key), nil
}

// List lists files with the given prefix.
func (s *GCSStorage) List(ctx context.Context, prefix string) ([]string, error) {
	bucket := s.client.Bucket(s.bucket)
	query := &storage.Query{Prefix: prefix}

	iter := bucket.Objects(ctx, query)
	var keys []string

	for {
		obj, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		keys = append(keys, obj.Name)
	}

	return keys, nil
}

// GCSMockStorage provides a mock implementation for testing.
type GCSMockStorage struct {
	files map[string][]byte
}

// NewGCSMockStorage creates a mock GCS storage for testing.
func NewGCSMockStorage() *GCSMockStorage {
	return &GCSMockStorage{
		files: make(map[string][]byte),
	}
}

func (m *GCSMockStorage) Upload(ctx context.Context, key string, data []byte, contentType string) error {
	m.files[key] = data
	return nil
}

func (m *GCSMockStorage) Download(ctx context.Context, key string) ([]byte, error) {
	data, ok := m.files[key]
	if !ok {
		return nil, ErrFileNotFound
	}
	return data, nil
}

func (m *GCSMockStorage) Delete(ctx context.Context, key string) error {
	delete(m.files, key)
	return nil
}

func (m *GCSMockStorage) GetURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	return fmt.Sprintf("https://storage.googleapis.com/mock/%s", key), nil
}

func (m *GCSMockStorage) List(ctx context.Context, prefix string) ([]string, error) {
	var keys []string
	for k := range m.files {
		if len(prefix) == 0 || len(k) >= len(prefix) && k[:len(prefix)] == prefix {
			keys = append(keys, k)
		}
	}
	return keys, nil
}
