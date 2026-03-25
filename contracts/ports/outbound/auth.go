package outbound

import (
	"context"
	"time"
)

// TokenPair contains access and refresh tokens.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
	TokenType    string
}

// UserClaims represents authenticated user information.
type UserClaims struct {
	UserID string
	Email  string
	Roles  []string
	Scope  string
}

// AuthPort defines the interface for authentication.
// Following Interface Segregation (ISP) - focused, minimal interface.
//
// Design Principles:
//   - ISP: Minimal interface for auth operations
//   - DIP: Domain depends on abstraction
//   - Low Coupling: Single port for auth concerns
type AuthPort interface {
	// ValidateToken validates an access token and returns claims.
	ValidateToken(ctx context.Context, token string) (*UserClaims, error)

	// RefreshToken validates a refresh token and returns new token pair.
	RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error)

	// InvalidateToken revokes a token.
	InvalidateToken(ctx context.Context, token string) error
}

// TokenGenerator defines the interface for token generation.
type TokenGenerator interface {
	// GenerateTokenPair creates new access and refresh tokens.
	GenerateTokenPair(ctx context.Context, userID, email string, roles []string) (*TokenPair, error)
}

// TokenValidator defines the interface for token validation.
type TokenValidator interface {
	// ValidateAccessToken validates an access token.
	ValidateAccessToken(ctx context.Context, token string) (*UserClaims, error)

	// ValidateRefreshToken validates a refresh token.
	ValidateRefreshToken(ctx context.Context, token string) (*UserClaims, error)
}

// APIKeyPort defines the interface for API key management.
type APIKeyPort interface {
	// CreateAPIKey creates a new API key.
	CreateAPIKey(ctx context.Context, userID, name string, scopes []string) (string, error)

	// ValidateAPIKey validates an API key.
	ValidateAPIKey(ctx context.Context, key string) (*APIKeyInfo, error)

	// RevokeAPIKey revokes an API key.
	RevokeAPIKey(ctx context.Context, keyID string) error
}

// APIKeyInfo contains API key metadata (without the secret).
type APIKeyInfo struct {
	ID         string
	Name       string
	Prefix     string
	UserID     string
	Scopes     []string
	CreatedAt  time.Time
	ExpiresAt  *time.Time
	LastUsedAt *time.Time
	RateLimit  int
}
