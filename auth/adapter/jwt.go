package adapter

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/KooshaPari/phenotype-go-kit/contracts/ports/outbound"
	"github.com/golang-jwt/jwt/v5"
)

var (
	_ outbound.AuthPort       = (*JWTValidatorAdapter)(nil)
	_ outbound.TokenGenerator = (*JWTValidatorAdapter)(nil)
	_ outbound.TokenValidator = (*JWTValidatorAdapter)(nil)
	_ outbound.APIKeyPort     = (*APIKeyManagerAdapter)(nil)
)

// Standard errors for auth adapters.
var (
	ErrInvalidToken  = errors.New("invalid token")
	ErrExpiredToken  = errors.New("token has expired")
	ErrInvalidClaims = errors.New("invalid claims")
)

// JWTClaims represents JWT claims following the port interface.
type JWTClaims struct {
	jwt.RegisteredClaims
	UserID string   `json:"user_id"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
	Scope  string   `json:"scope"`
}

// Config holds JWT configuration.
type Config struct {
	SecretKey          string
	PrivateKey         *rsa.PrivateKey
	PublicKey          *rsa.PublicKey
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	Issuer             string
	Audience           string
}

// JWTValidatorAdapter implements outbound.AuthPort using JWT.
type JWTValidatorAdapter struct {
	config Config
	logger *slog.Logger
}

// NewJWTValidatorAdapter creates a new JWT validator adapter.
func NewJWTValidatorAdapter(cfg Config) *JWTValidatorAdapter {
	return &JWTValidatorAdapter{
		config: cfg,
		logger: slog.Default(),
	}
}

// Manifest returns the adapter manifest.
func (a *JWTValidatorAdapter) Manifest() *outbound.AdapterManifest {
	return &outbound.AdapterManifest{
		Name:        "jwt-auth-adapter",
		Version:     "1.0.0",
		Description: "JWT authentication adapter",
		Provides:    []string{"auth-port", "token-generator", "token-validator"},
	}
}

// ValidateToken validates an access token and returns claims.
func (a *JWTValidatorAdapter) ValidateToken(ctx context.Context, token string) (*outbound.UserClaims, error) {
	claims, err := a.validateToken(token, "access")
	if err != nil {
		return nil, err
	}
	return a.toUserClaims(claims), nil
}

// RefreshToken validates a refresh token and returns new token pair.
func (a *JWTValidatorAdapter) RefreshToken(ctx context.Context, refreshToken string) (*outbound.TokenPair, error) {
	claims, err := a.validateToken(refreshToken, "refresh")
	if err != nil {
		return nil, err
	}

	return a.GenerateTokenPair(ctx, claims.UserID, claims.Email, claims.Roles)
}

// InvalidateToken revokes a token.
func (a *JWTValidatorAdapter) InvalidateToken(ctx context.Context, token string) error {
	hash := sha256.Sum256([]byte(token))
	tokenHash := base64.StdEncoding.EncodeToString(hash[:])
	a.logger.Info("token invalidated", "hash", tokenHash[:8])
	return nil
}

// GenerateTokenPair creates new access and refresh tokens.
func (a *JWTValidatorAdapter) GenerateTokenPair(ctx context.Context, userID, email string, roles []string) (*outbound.TokenPair, error) {
	now := time.Now()

	accessClaims := JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        newTokenID(),
			ExpiresAt: jwt.NewNumericDate(now.Add(a.config.AccessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    a.config.Issuer,
			Audience:  jwt.ClaimStrings{a.config.Audience},
		},
		UserID: userID,
		Email:  email,
		Roles:  roles,
		Scope:  "access",
	}

	refreshClaims := JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        newTokenID(),
			ExpiresAt: jwt.NewNumericDate(now.Add(a.config.RefreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    a.config.Issuer,
			Audience:  jwt.ClaimStrings{a.config.Audience},
		},
		UserID: userID,
		Email:  email,
		Roles:  roles,
		Scope:  "refresh",
	}

	accessToken, err := a.signClaims(accessClaims)
	if err != nil {
		return nil, err
	}

	refreshToken, err := a.signClaims(refreshClaims)
	if err != nil {
		return nil, err
	}

	return &outbound.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(a.config.AccessTokenExpiry.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// ValidateAccessToken validates an access token.
func (a *JWTValidatorAdapter) ValidateAccessToken(ctx context.Context, token string) (*outbound.UserClaims, error) {
	claims, err := a.validateToken(token, "access")
	if err != nil {
		return nil, err
	}

	if claims.Scope != "access" {
		return nil, ErrInvalidClaims
	}

	return a.toUserClaims(claims), nil
}

// ValidateRefreshToken validates a refresh token.
func (a *JWTValidatorAdapter) ValidateRefreshToken(ctx context.Context, token string) (*outbound.UserClaims, error) {
	claims, err := a.validateToken(token, "refresh")
	if err != nil {
		return nil, err
	}

	if claims.Scope != "refresh" {
		return nil, ErrInvalidClaims
	}

	return a.toUserClaims(claims), nil
}

func (a *JWTValidatorAdapter) signClaims(claims JWTClaims) (string, error) {
	var token *jwt.Token

	if a.config.PrivateKey != nil {
		token = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		return token.SignedString(a.config.PrivateKey)
	}

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.config.SecretKey))
}

func (a *JWTValidatorAdapter) validateToken(tokenString, expectedScope string) (*JWTClaims, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
			return []byte(a.config.SecretKey), nil
		}
		if _, ok := token.Method.(*jwt.SigningMethodRSA); ok {
			return a.config.PublicKey, nil
		}
		return nil, ErrInvalidToken
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, keyFunc)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.Scope != expectedScope {
		return nil, ErrInvalidClaims
	}

	return claims, nil
}

func (a *JWTValidatorAdapter) toUserClaims(claims *JWTClaims) *outbound.UserClaims {
	return &outbound.UserClaims{
		UserID: claims.UserID,
		Email:  claims.Email,
		Roles:  claims.Roles,
		Scope:  claims.Scope,
	}
}

func newTokenID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("token-%d", time.Now().UnixNano())
	}
	return base64.RawURLEncoding.EncodeToString(bytes)
}

// APIKeyManagerAdapter implements outbound.APIKeyPort.
type APIKeyManagerAdapter struct {
	keys   map[string]*APIKeyEntry
	logger *slog.Logger
}

// APIKeyEntry represents an API key entry.
type APIKeyEntry struct {
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

// NewAPIKeyManagerAdapter creates a new API key manager adapter.
func NewAPIKeyManagerAdapter() *APIKeyManagerAdapter {
	return &APIKeyManagerAdapter{
		keys:   make(map[string]*APIKeyEntry),
		logger: slog.Default(),
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

func generateID() string {
	bytes := make([]byte, 16)
	_, _ = rand.Read(bytes)
	return fmt.Sprintf("%x", bytes)
}

// CreateAPIKey creates a new API key.
func (m *APIKeyManagerAdapter) CreateAPIKey(ctx context.Context, userID, name string, scopes []string) (string, error) {
	key, err := GenerateAPIKey("pk")
	if err != nil {
		return "", err
	}

	prefix := key[:8]
	entry := &APIKeyEntry{
		ID:        generateID(),
		Name:      name,
		Prefix:    prefix,
		Hash:      HashAPIKey(key),
		UserID:    userID,
		Scopes:    scopes,
		CreatedAt: time.Now(),
		RateLimit: 1000,
	}

	m.keys[key] = entry
	return key, nil
}

// ValidateAPIKey validates an API key.
func (m *APIKeyManagerAdapter) ValidateAPIKey(ctx context.Context, key string) (*outbound.APIKeyInfo, error) {
	entry, ok := m.keys[key]
	if !ok {
		return nil, errors.New("invalid API key")
	}

	now := time.Now()
	entry.LastUsedAt = &now

	return &outbound.APIKeyInfo{
		ID:         entry.ID,
		Name:       entry.Name,
		Prefix:     entry.Prefix,
		UserID:     entry.UserID,
		Scopes:     entry.Scopes,
		CreatedAt:  entry.CreatedAt,
		ExpiresAt:  entry.ExpiresAt,
		LastUsedAt: entry.LastUsedAt,
		RateLimit:  entry.RateLimit,
	}, nil
}

// RevokeAPIKey revokes an API key.
func (m *APIKeyManagerAdapter) RevokeAPIKey(ctx context.Context, keyID string) error {
	for k, v := range m.keys {
		if v.ID == keyID {
			delete(m.keys, k)
			return nil
		}
	}
	return errors.New("key not found")
}

// ContextUser extracts user info from context.
type ContextUser struct {
	UserID string
	Email  string
	Roles  []string
}

// GetUserFromContext extracts user from context following LoD.
func GetUserFromContext(ctx context.Context) *ContextUser {
	return &ContextUser{
		UserID: getStringFromContext(ctx, "user_id"),
		Email:  getStringFromContext(ctx, "user_email"),
		Roles:  getStringSliceFromContext(ctx, "user_roles"),
	}
}

func getStringFromContext(ctx context.Context, key string) string {
	if v := ctx.Value(key); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getStringSliceFromContext(ctx context.Context, key string) []string {
	if v := ctx.Value(key); v != nil {
		if s, ok := v.([]string); ok {
			return s
		}
	}
	return nil
}

// HasRole checks if user has required role.
func (u *ContextUser) HasRole(required string) bool {
	for _, role := range u.Roles {
		if role == required {
			return true
		}
	}
	return false
}

// HasAnyRole checks if user has any of the required roles.
func (u *ContextUser) HasAnyRole(required ...string) bool {
	for _, req := range required {
		if u.HasRole(req) {
			return true
		}
	}
	return false
}

// HTTPHeaders extracts auth info for HTTP headers.
func (u *ContextUser) HTTPHeaders() map[string]string {
	return map[string]string{
		"X-User-ID":    u.UserID,
		"X-User-Email": u.Email,
		"X-User-Roles": strings.Join(u.Roles, ","),
	}
}
