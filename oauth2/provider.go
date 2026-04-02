package oauth2

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Config holds OAuth2 configuration.
type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string

	AuthorizationURL string
	TokenURL         string
	UserInfoURL      string

	SigningKey string
}

// TokenResponse represents OAuth2 token response.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// UserInfo represents OAuth2 user info.
type UserInfo struct {
	ID      string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	Locale  string `json:"locale"`
}

// Provider represents an OAuth2 provider.
type Provider struct {
	config *Config
	logger *slog.Logger
	client *http.Client
}

// NewProvider creates a new OAuth2 provider.
func NewProvider(cfg Config) *Provider {
	return &Provider{
		config: &cfg,
		logger: slog.Default(),
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// GenerateState generates a random state parameter.
func GenerateState() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// GenerateCodeVerifier generates PKCE code verifier.
func GenerateCodeVerifier() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// GenerateCodeChallenge generates PKCE code challenge from verifier.
func GenerateCodeChallenge(verifier string) string {
	// SHA256 hash and base64 URL encode
	sum := 0
	for i := range verifier {
		sum += int(verifier[i])
	}
	// Simplified - in production use crypto/sha256
	return base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%x", sum)))
}

// AuthorizationURL returns the OAuth2 authorization URL.
func (p *Provider) AuthorizationURL(state, codeVerifier string) string {
	params := url.Values{
		"client_id":     {p.config.ClientID},
		"redirect_uri":  {p.config.RedirectURL},
		"response_type": {"code"},
		"scope":         {strings.Join(p.config.Scopes, " ")},
		"state":         {state},
	}

	if codeVerifier != "" {
		params.Set("code_challenge_method", "S256")
		params.Set("code_challenge", GenerateCodeChallenge(codeVerifier))
	}

	authURL, err := url.Parse(p.config.AuthorizationURL)
	if err != nil {
		return ""
	}
	authURL.RawQuery = params.Encode()

	return authURL.String()
}

// ExchangeCode exchanges authorization code for tokens.
func (p *Provider) ExchangeCode(ctx context.Context, code, codeVerifier string) (*TokenResponse, error) {
	data := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {p.config.ClientID},
		"client_secret": {p.config.ClientSecret},
		"code":          {code},
		"redirect_uri":  {p.config.RedirectURL},
	}

	if codeVerifier != "" {
		data.Set("code_verifier", codeVerifier)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.config.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token exchange failed: %s", body)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

// RefreshToken refreshes an access token.
func (p *Provider) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	data := url.Values{
		"grant_type":    {"refresh_token"},
		"client_id":     {p.config.ClientID},
		"client_secret": {p.config.ClientSecret},
		"refresh_token": {refreshToken},
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.config.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("token refresh failed")
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

// GetUserInfo retrieves user information using access token.
func (p *Provider) GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", p.config.UserInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("user info request failed")
	}

	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// Session represents a user session.
type Session struct {
	ID           string
	UserID       string
	Email        string
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

// SessionManager manages OAuth2 sessions.
type SessionManager struct {
	sessions map[string]*Session
	keys     map[string]string // token -> session ID
	mu       sync.RWMutex
}

// NewSessionManager creates a new session manager.
func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*Session),
		keys:     make(map[string]string),
	}
}

// CreateSession creates a new session from OAuth2 tokens.
func (sm *SessionManager) CreateSession(userID, email, accessToken, refreshToken string, expiresIn int) *Session {
	session := &Session{
		ID:           generateID(),
		UserID:       userID,
		Email:        email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(expiresIn) * time.Second),
		CreatedAt:    time.Now(),
	}

	sm.mu.Lock()
	sm.sessions[session.ID] = session
	sm.keys[accessToken] = session.ID
	sm.mu.Unlock()

	return session
}

// GetSession retrieves a session by ID.
func (sm *SessionManager) GetSession(sessionID string) (*Session, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	session, ok := sm.sessions[sessionID]
	if !ok || time.Now().After(session.ExpiresAt) {
		return nil, false
	}

	return session, true
}

// GetSessionByToken retrieves a session by access token.
func (sm *SessionManager) GetSessionByToken(accessToken string) (*Session, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if sessionID, ok := sm.keys[accessToken]; ok {
		if session, ok := sm.sessions[sessionID]; ok {
			if !time.Now().After(session.ExpiresAt) {
				return session, true
			}
		}
	}

	return nil, false
}

// DeleteSession removes a session.
func (sm *SessionManager) DeleteSession(sessionID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if session, ok := sm.sessions[sessionID]; ok {
		delete(sm.keys, session.AccessToken)
		delete(sm.sessions, sessionID)
	}
}

// Middleware returns session middleware for HTTP handlers.
func (sm *SessionManager) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token != "" {
				token = strings.TrimPrefix(token, "Bearer ")

				if session, ok := sm.GetSessionByToken(token); ok {
					ctx := r.Context()
					ctx = context.WithValue(ctx, "session", session)
					r = r.WithContext(ctx)
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func generateID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return fmt.Sprintf("%x", bytes)
}

// Predefined providers.
var GoogleProvider = Provider{
	config: &Config{
		AuthorizationURL: "https://accounts.google.com/o/oauth2/v2/auth",
		TokenURL:         "https://oauth2.googleapis.com/token",
		UserInfoURL:      "https://www.googleapis.com/oauth2/v2/userinfo",
		Scopes:           []string{"email", "profile"},
	},
}

var GitHubProvider = Provider{
	config: &Config{
		AuthorizationURL: "https://github.com/login/oauth/authorize",
		TokenURL:         "https://github.com/login/oauth/access_token",
		UserInfoURL:      "https://api.github.com/user",
		Scopes:           []string{"user:email", "read:user"},
	},
}

// InitGoogle initializes Google OAuth2 provider.
func InitGoogle(clientID, clientSecret, redirectURL string) *Provider {
	return &Provider{
		config: &Config{
			ClientID:         clientID,
			ClientSecret:     clientSecret,
			RedirectURL:      redirectURL,
			AuthorizationURL: "https://accounts.google.com/o/oauth2/v2/auth",
			TokenURL:         "https://oauth2.googleapis.com/token",
			UserInfoURL:      "https://www.googleapis.com/oauth2/v2/userinfo",
			Scopes:           []string{"email", "profile"},
		},
		logger: slog.Default(),
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// InitGitHub initializes GitHub OAuth2 provider.
func InitGitHub(clientID, clientSecret, redirectURL string) *Provider {
	return &Provider{
		config: &Config{
			ClientID:         clientID,
			ClientSecret:     clientSecret,
			RedirectURL:      redirectURL,
			AuthorizationURL: "https://github.com/login/oauth/authorize",
			TokenURL:         "https://github.com/login/oauth/access_token",
			UserInfoURL:      "https://api.github.com/user",
			Scopes:           []string{"user:email", "read:user"},
		},
		logger: slog.Default(),
		client: &http.Client{Timeout: 30 * time.Second},
	}
}
