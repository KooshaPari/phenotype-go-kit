package inbound

import (
	"time"
)

// AuthCommands defines CQRS commands for authentication operations.
type AuthCommands struct{}

// NewAuthCommands creates new auth command handlers.
func NewAuthCommands() *AuthCommands {
	return &AuthCommands{}
}

// LoginCommand represents user login.
type LoginCommand struct {
	Email    string
	Password string
}

// LoginHandler handles LoginCommand.
type LoginHandler func(cmd LoginCommand) (*AuthResponse, error)

// RefreshTokenCommand represents token refresh.
type RefreshTokenCommand struct {
	RefreshToken string
}

// RefreshTokenHandler handles RefreshTokenCommand.
type RefreshTokenHandler func(cmd RefreshTokenCommand) (*AuthResponse, error)

// LogoutCommand represents user logout.
type LogoutCommand struct {
	AccessToken string
}

// LogoutHandler handles LogoutCommand.
type LogoutHandler func(cmd LogoutCommand) error

// CreateAPIKeyCommand represents API key creation.
type CreateAPIKeyCommand struct {
	UserID    string
	Name      string
	Scopes    []string
	RateLimit int
}

// CreateAPIKeyHandler handles CreateAPIKeyCommand.
type CreateAPIKeyHandler func(cmd CreateAPIKeyCommand) (*CreateAPIKeyResponse, error)

// RevokeAPIKeyCommand represents API key revocation.
type RevokeAPIKeyCommand struct {
	KeyID string
}

// RevokeAPIKeyHandler handles RevokeAPIKeyCommand.
type RevokeAPIKeyHandler func(cmd RevokeAPIKeyCommand) error

// AuthResponse contains the response for auth operations.
type AuthResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token,omitempty"`
	ExpiresIn    int64    `json:"expires_in"`
	TokenType    string   `json:"token_type"`
	UserID       string   `json:"user_id"`
	Email        string   `json:"email"`
	Roles        []string `json:"roles"`
}

// CreateAPIKeyResponse contains the response for API key creation.
type CreateAPIKeyResponse struct {
	Key       string    `json:"key"`
	KeyID     string    `json:"key_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// AuthQueries defines CQRS queries for authentication.
type AuthQueries struct{}

// NewAuthQueries creates new auth query handlers.
func NewAuthQueries() *AuthQueries {
	return &AuthQueries{}
}

// ValidateTokenQuery represents token validation.
type ValidateTokenQuery struct {
	Token string
}

// ValidateTokenHandler handles ValidateTokenQuery.
type ValidateTokenHandler func(query ValidateTokenQuery) (*ValidateTokenResponse, error)

// ValidateTokenResponse contains token validation result.
type ValidateTokenResponse struct {
	Valid  bool
	UserID string
	Email  string
	Roles  []string
	Scope  string
}

// GetAPIKeysQuery represents fetching user's API keys.
type GetAPIKeysQuery struct {
	UserID string
}

// GetAPIKeysHandler handles GetAPIKeysQuery.
type GetAPIKeysHandler func(query GetAPIKeysQuery) ([]*APIKeyInfo, error)

// APIKeyInfo contains API key metadata.
type APIKeyInfo struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Prefix     string     `json:"prefix"`
	Scopes     []string   `json:"scopes"`
	CreatedAt  time.Time  `json:"created_at"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	RateLimit  int        `json:"rate_limit"`
}
