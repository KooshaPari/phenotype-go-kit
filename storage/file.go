package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/url"
	"path"
	"time"
)

var (
	ErrFileNotFound = errors.New("file not found")
	ErrInvalidPath  = errors.New("invalid file path")
)

// ObjectStorage defines the interface for cloud object storage.
type ObjectStorage interface {
	// Upload uploads a file to the storage.
	Upload(ctx context.Context, key string, data []byte, contentType string) error

	// Download downloads a file from the storage.
	Download(ctx context.Context, key string) ([]byte, error)

	// Delete deletes a file from the storage.
	Delete(ctx context.Context, key string) error

	// GetURL returns a signed URL for the file.
	GetURL(ctx context.Context, key string, expiry time.Duration) (string, error)

	// List lists files with the given prefix.
	List(ctx context.Context, prefix string) ([]string, error)
}

// S3Config holds configuration for AWS S3.
type S3Config struct {
	Bucket    string
	Region    string
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
}

// FileService handles file operations.
type FileService struct {
	storage ObjectStorage
	bucket  string
	baseURL string
	logger  *slog.Logger
}

// NewFileService creates a new file service.
func NewFileService(storage ObjectStorage, bucket, baseURL string, logger *slog.Logger) *FileService {
	return &FileService{
		storage: storage,
		bucket:  bucket,
		baseURL: baseURL,
		logger:  logger,
	}
}

// UploadRequest represents a file upload request.
type UploadRequest struct {
	Filename    string
	ContentType string
	Data        []byte
}

// UploadResponse represents the response from an upload.
type UploadResponse struct {
	Key        string
	URL        string
	Size       int64
	UploadedAt time.Time
}

// Upload handles file upload.
func (s *FileService) Upload(ctx context.Context, req UploadRequest) (*UploadResponse, error) {
	if req.Filename == "" {
		return nil, ErrInvalidPath
	}

	key := s.generateKey(req.Filename)

	err := s.storage.Upload(ctx, key, req.Data, req.ContentType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload: %w", err)
	}

	s.logger.Info("file uploaded", "key", key, "size", len(req.Data))

	url, err := s.storage.GetURL(ctx, key, 24*time.Hour)
	if err != nil {
		s.logger.Warn("failed to get signed URL", "key", key)
	}

	return &UploadResponse{
		Key:        key,
		URL:        url,
		Size:       int64(len(req.Data)),
		UploadedAt: time.Now(),
	}, nil
}

// Download handles file download.
func (s *FileService) Download(ctx context.Context, key string) ([]byte, error) {
	data, err := s.storage.Download(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to download: %w", err)
	}

	return data, nil
}

// Delete handles file deletion.
func (s *FileService) Delete(ctx context.Context, key string) error {
	err := s.storage.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}

	s.logger.Info("file deleted", "key", key)
	return nil
}

// GetDownloadURL returns a signed download URL.
func (s *FileService) GetDownloadURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	return s.storage.GetURL(ctx, key, expiry)
}

// ParseMultipartForm parses a multipart form and extracts files.
func ParseMultipartForm(form *multipart.Form) ([]UploadRequest, error) {
	var uploads []UploadRequest

	for filename, headers := range form.File {
		for _, header := range headers {
			file, err := header.Open()
			if err != nil {
				return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
			}

			data, err := io.ReadAll(file)
			if err != nil {
				file.Close()
				return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
			}
			file.Close()

			uploads = append(uploads, UploadRequest{
				Filename:    header.Filename,
				ContentType: header.Header.Get("Content-Type"),
				Data:        data,
			})
		}
	}

	return uploads, nil
}

func (s *FileService) generateKey(filename string) string {
	// Generate a unique key with timestamp prefix
	timestamp := time.Now().Format("2006/01/02")
	ext := path.Ext(filename)
	name := filename[:len(filename)-len(ext)]
	return fmt.Sprintf("%s/%s-%s%s", timestamp, name, randomID(8), ext)
}

func randomID(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[int(time.Now().UnixNano())%len(letters)]
	}
	return string(b)
}

// MockStorage implements ObjectStorage for testing.
type MockStorage struct {
	files map[string][]byte
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		files: make(map[string][]byte),
	}
}

func (m *MockStorage) Upload(ctx context.Context, key string, data []byte, contentType string) error {
	m.files[key] = data
	return nil
}

func (m *MockStorage) Download(ctx context.Context, key string) ([]byte, error) {
	data, ok := m.files[key]
	if !ok {
		return nil, ErrFileNotFound
	}
	return data, nil
}

func (m *MockStorage) Delete(ctx context.Context, key string) error {
	delete(m.files, key)
	return nil
}

func (m *MockStorage) GetURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	// Return a mock URL
	parsed, _ := url.Parse("/download")
	query := parsed.Query()
	query.Set("key", key)
	parsed.RawQuery = query.Encode()
	return parsed.String(), nil
}

func (m *MockStorage) List(ctx context.Context, prefix string) ([]string, error) {
	var keys []string
	for k := range m.files {
		if len(prefix) == 0 || len(k) >= len(prefix) && k[:len(prefix)] == prefix {
			keys = append(keys, k)
		}
	}
	return keys, nil
}

// BytesReader wraps bytes as an io.Reader for testing.
type BytesReader struct {
	data   []byte
	offset int
}

func NewBytesReader(data []byte) *BytesReader {
	return &BytesReader{data: data}
}

func (r *BytesReader) Read(p []byte) (n int, err error) {
	if r.offset >= len(r.data) {
		return 0, io.EOF
	}
	n = copy(p, r.data[r.offset:])
	r.offset += n
	return n, nil
}

var _ io.Reader = (*BytesReader)(nil)

// BytesBuffer wraps a bytes.Buffer for testing.
type BytesBuffer struct {
	buf bytes.Buffer
}

func NewBytesBuffer() *BytesBuffer {
	return &BytesBuffer{}
}

func (b *BytesBuffer) Write(p []byte) (int, error) {
	return b.buf.Write(p)
}

func (b *BytesBuffer) Bytes() []byte {
	return b.buf.Bytes()
}
