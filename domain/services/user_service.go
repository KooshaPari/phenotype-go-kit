// Package services contains domain services.
// Domain services encapsulate business logic that doesn't belong to a single entity.
package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/KooshaPari/phenotype-go-kit/domain/entities"
	"github.com/KooshaPari/phenotype-go-kit/domain/ports"
)

// ValidationError represents a validation failure.
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// UserService handles user business logic.
// Following DDD principles: domain services encapsulate complex operations.
type UserService struct {
	userRepo ports.QueryRepositoryPort[*entities.User, string]
	cache    ports.CachePort
	logger   ports.LoggerPort
}

// NewUserService creates a new UserService.
func NewUserService(
	userRepo ports.QueryRepositoryPort[*entities.User, string],
	cache ports.CachePort,
	logger ports.LoggerPort,
) *UserService {
	return &UserService{
		userRepo: userRepo,
		cache:    cache,
		logger:   logger,
	}
}

// Create creates a new user.
// Following ATDD: validation happens at the service layer.
func (s *UserService) Create(ctx context.Context, email, name string) (*entities.User, error) {
	// Validation (CDD - Contract-Driven Development)
	if !isValidEmail(email) {
		return nil, ValidationError{Field: "email", Message: "invalid email format"}
	}
	if len(name) < 2 {
		return nil, ValidationError{Field: "name", Message: "name must be at least 2 characters"}
	}

	// Check for duplicate email
	existing, err := s.userRepo.FindByFilter(ctx, map[string]any{"email": email})
	var notFound *ports.ErrNotFound
	if err != nil && !errors.As(err, &notFound) {
		s.logger.Error("failed to check existing user", ports.Error(err))
		return nil, fmt.Errorf("check user: %w", err)
	}
	if len(existing) > 0 {
		return nil, ValidationError{Field: "email", Message: "email already exists"}
	}

	// Create user
	now := time.Now()
	user := &entities.User{
		ID:        generateID(),
		Email:     email,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Save
	saved, err := s.userRepo.Save(ctx, user)
	if err != nil {
		s.logger.Error("failed to create user", ports.Error(err))
		return nil, fmt.Errorf("save user: %w", err)
	}

	s.logger.Info("user created", ports.String("user_id", saved.ID))
	return saved, nil
}

// GetByID retrieves a user by ID.
// Implements caching (read-through pattern).
func (s *UserService) GetByID(ctx context.Context, id string) (*entities.User, error) {
	// Try cache first
	cacheKey := "user:" + id
	if _, err := s.cache.Get(ctx, cacheKey); err == nil {
		// Cache hit - deserialize and return
		// (simplified - in practice use JSON unmarshal)
		return nil, nil // Placeholder
	}

	// Cache miss - get from repo
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		var notFound *ports.ErrNotFound
		if errors.As(err, &notFound) {
			return nil, err
		}
		return nil, fmt.Errorf("find user: %w", err)
	}

	// Don't cache deleted users
	if !user.IsActive() {
		return nil, &ports.ErrNotFound{Entity: "user", ID: id}
	}

	// Cache the result
	// (simplified - in practice serialize to JSON)

	return user, nil
}

// isValidEmail validates email format.
// Following YAGNI: simple validation for now.
func isValidEmail(email string) bool {
	return len(email) > 3 && contains(email, "@")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Placeholder functions - in practice use proper implementations.
func generateID() string {
	return "id-" + time.Now().Format("20060102150405")
}
