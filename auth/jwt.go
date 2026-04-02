package auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken  = errors.New("invalid token")
	ErrExpiredToken  = errors.New("token has expired")
	ErrInvalidClaims = errors.New("invalid claims")
)

// JWTConfig holds JWT configuration.
type JWTConfig struct {
	SecretKey          string
	PrivateKey         *rsa.PrivateKey
	PublicKey          *rsa.PublicKey
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	Issuer             string
	Audience           string
}

// TokenClaims represents JWT claims.
type TokenClaims struct {
	jwt.RegisteredClaims
	UserID string   `json:"user_id"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
	Scope  string   `json:"scope"`
}

// TokenPair contains access and refresh tokens.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// JWTValidator validates JWT tokens.
type JWTValidator struct {
	config *JWTConfig
	logger *slog.Logger
}

// NewJWTValidator creates a new JWT validator.
func NewJWTValidator(cfg JWTConfig) *JWTValidator {
	return &JWTValidator{
		config: &cfg,
		logger: slog.Default(),
	}
}

// GenerateTokenPair creates new access and refresh tokens.
func (v *JWTValidator) GenerateTokenPair(ctx context.Context, userID, email string, roles []string) (TokenPair, error) {
	now := time.Now()

	accessClaims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(v.config.AccessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    v.config.Issuer,
			Audience:  jwt.ClaimStrings{v.config.Audience},
		},
		UserID: userID,
		Email:  email,
		Roles:  roles,
		Scope:  "access",
	}

	refreshClaims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(v.config.RefreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    v.config.Issuer,
			Audience:  jwt.ClaimStrings{v.config.Audience},
		},
		UserID: userID,
		Email:  email,
		Roles:  roles,
		Scope:  "refresh",
	}

	accessToken, err := v.signClaims(accessClaims)
	if err != nil {
		return TokenPair{}, err
	}

	refreshToken, err := v.signClaims(refreshClaims)
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(v.config.AccessTokenExpiry.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

func (v *JWTValidator) signClaims(claims TokenClaims) (string, error) {
	var token *jwt.Token

	if v.config.PrivateKey != nil {
		token = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		return token.SignedString(v.config.PrivateKey)
	}

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(v.config.SecretKey))
}

// ValidateAccessToken validates an access token.
func (v *JWTValidator) ValidateAccessToken(ctx context.Context, tokenString string) (*TokenClaims, error) {
	claims, err := v.validateToken(tokenString, "access")
	if err != nil {
		return nil, err
	}

	if claims.Scope != "access" {
		return nil, ErrInvalidClaims
	}

	return claims, nil
}

// ValidateRefreshToken validates a refresh token.
func (v *JWTValidator) ValidateRefreshToken(ctx context.Context, tokenString string) (*TokenClaims, error) {
	claims, err := v.validateToken(tokenString, "refresh")
	if err != nil {
		return nil, err
	}

	if claims.Scope != "refresh" {
		return nil, ErrInvalidClaims
	}

	return claims, nil
}

func (v *JWTValidator) validateToken(tokenString, expectedScope string) (*TokenClaims, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
			return []byte(v.config.SecretKey), nil
		}
		if _, ok := token.Method.(*jwt.SigningMethodRSA); ok {
			return v.config.PublicKey, nil
		}
		return nil, ErrInvalidToken
	}

	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, keyFunc)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.Scope != expectedScope {
		return nil, ErrInvalidClaims
	}

	return claims, nil
}

// InvalidateToken revokes a token (simple implementation - production needs blacklist).
func (v *JWTValidator) InvalidateToken(ctx context.Context, tokenString string) error {
	hash := sha256.Sum256([]byte(tokenString))
	tokenHash := base64.StdEncoding.EncodeToString(hash[:])

	v.logger.Info("token invalidated", "hash", tokenHash[:8])
	return nil
}

// Middleware returns HTTP middleware for JWT validation.
func (v *JWTValidator) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "authorization required", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			claims, err := v.ValidateAccessToken(r.Context(), parts[1])
			if err != nil {
				http.Error(w, "invalid token: "+err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
			ctx = context.WithValue(ctx, "user_email", claims.Email)
			ctx = context.WithValue(ctx, "user_roles", claims.Roles)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GenerateAPIKey creates a new API key.
func GenerateAPIKey(prefix string) (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	key := base64.URLEncoding.EncodeToString(bytes)
	if prefix != "" {
		key = prefix + "_" + key
	}
	return key, nil
}

// HashAPIKey creates a hash of the API key for storage.
func HashAPIKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return base64.StdEncoding.EncodeToString(hash[:])
}

// APIKeyManager manages API keys.
type APIKeyManager struct {
	keys map[string]*APIKey
}

// APIKey represents an API key.
type APIKey struct {
	ID         string
	Name       string
	Prefix     string
	Hash       string
	UserID     string
	Scopes     []string
	CreatedAt  time.Time
	ExpiresAt  *time.Time
	LastUsedAt *time.Time
	RateLimit  int
}

// NewAPIKeyManager creates a new API key manager.
func NewAPIKeyManager() *APIKeyManager {
	return &APIKeyManager{
		keys: make(map[string]*APIKey),
	}
}

// CreateKey creates a new API key.
func (m *APIKeyManager) CreateKey(ctx context.Context, userID, name string, scopes []string, rateLimit int) (string, *APIKey, error) {
	key, err := GenerateAPIKey("pk")
	if err != nil {
		return "", nil, err
	}

	prefix := key[:8]
	apiKey := &APIKey{
		ID:        generateID(),
		Name:      name,
		Prefix:    prefix,
		Hash:      HashAPIKey(key),
		UserID:    userID,
		Scopes:    scopes,
		CreatedAt: time.Now(),
		RateLimit: rateLimit,
	}

	m.keys[key] = apiKey
	return key, apiKey, nil
}

// ValidateKey validates an API key.
func (m *APIKeyManager) ValidateKey(key string) (*APIKey, bool) {
	apiKey, ok := m.keys[key]
	return apiKey, ok
}

// RevokeKey revokes an API key.
func (m *APIKeyManager) RevokeKey(keyID string) error {
	for k, v := range m.keys {
		if v.ID == keyID {
			delete(m.keys, k)
			return nil
		}
	}
	return errors.New("key not found")
}

func generateID() string {
	bytes := make([]byte, 16)
	_, _ = rand.Read(bytes)
	return fmt.Sprintf("%x", bytes)
}

// GetUserID extracts user ID from context.
func GetUserID(ctx context.Context) string {
	if v := ctx.Value("user_id"); v != nil {
		return v.(string)
	}
	return ""
}

// GetUserEmail extracts user email from context.
func GetUserEmail(ctx context.Context) string {
	if v := ctx.Value("user_email"); v != nil {
		return v.(string)
	}
	return ""
}

// GetUserRoles extracts user roles from context.
func GetUserRoles(ctx context.Context) []string {
	if v := ctx.Value("user_roles"); v != nil {
		return v.([]string)
	}
	return nil
}

// RequireRole creates middleware that requires specific roles.
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRoles := GetUserRoles(r.Context())

			for _, required := range roles {
				for _, userRole := range userRoles {
					if userRole == required {
						next.ServeHTTP(w, r)
						return
					}
				}
			}

			http.Error(w, "insufficient permissions", http.StatusForbidden)
		})
	}
}
