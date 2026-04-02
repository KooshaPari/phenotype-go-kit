package versioning

import (
	"context"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const (
	DefaultVersion = "v1"
	HeaderVersion  = "Accept-Version"
	QueryVersion   = "version"
)

// Config holds API versioning configuration.
type Config struct {
	DefaultVersion  string
	AllowedVersions []string
	VersionHeader   string
	VersionQuery    string
	VersionRegex    string
}

// Middleware provides API versioning.
type Middleware struct {
	config     Config
	versionMap map[string]http.Handler
}

// New creates a new API versioning middleware.
func New(cfg Config) *Middleware {
	if cfg.DefaultVersion == "" {
		cfg.DefaultVersion = DefaultVersion
	}

	return &Middleware{
		config:     cfg,
		versionMap: make(map[string]http.Handler),
	}
}

// RegisterHandler registers a handler for a specific version.
func (v *Middleware) RegisterHandler(version string, handler http.Handler) {
	v.versionMap[version] = handler
}

// Middleware returns the HTTP middleware.
func (v *Middleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			version := v.extractVersion(r)

			if version == "" {
				version = v.config.DefaultVersion
			}

			if !v.isVersionAllowed(version) {
				http.Error(w, "version not supported", http.StatusNotAcceptable)
				return
			}

			r = r.WithContext(withVersion(r.Context(), version))

			w.Header().Set("API-Version", version)

			next.ServeHTTP(w, r)
		})
	}
}

func (v *Middleware) extractVersion(r *http.Request) string {
	if v.config.VersionHeader != "" {
		if version := r.Header.Get(v.config.VersionHeader); version != "" {
			return version
		}
	}

	if accept := r.Header.Get("Accept"); accept != "" {
		if version := extractFromAccept(accept, v.config.VersionRegex); version != "" {
			return version
		}
	}

	if v.config.VersionQuery != "" {
		if version := r.URL.Query().Get(v.config.VersionQuery); version != "" {
			return version
		}
	}

	return ""
}

func extractFromAccept(accept, regexPattern string) string {
	pattern := `vnd\.phenotype\.v(\d+)`
	if regexPattern != "" {
		pattern = regexPattern
	}

	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(accept)
	if len(matches) > 1 {
		return "v" + matches[1]
	}

	if strings.HasPrefix(accept, "application/vnd.phenotype.") {
		parts := strings.Split(accept, ".")
		for i, part := range parts {
			if strings.HasPrefix(part, "v") {
				if n, err := strconv.Atoi(part[1:]); err == nil {
					return "v" + string(rune(n+'0'))
				}
			}
			if i > 0 && strings.HasPrefix(parts[i-1], "vnd.phenotype") {
				if strings.HasPrefix(part, "v") {
					return part
				}
			}
		}
	}

	return ""
}

func (v *Middleware) isVersionAllowed(version string) bool {
	if len(v.config.AllowedVersions) == 0 {
		return true
	}

	for _, allowed := range v.config.AllowedVersions {
		if allowed == version {
			return true
		}
	}

	return false
}

type contextKey int

const versionKey contextKey = 0

// GetVersion extracts version from context.
func GetVersion(ctx context.Context) string {
	if v := ctx.Value(versionKey); v != nil {
		return v.(string)
	}
	return ""
}

func withVersion(ctx context.Context, version string) context.Context {
	return context.WithValue(ctx, versionKey, version)
}

// VersionHandler routes to version-specific handlers.
func VersionHandler(versions map[string]http.Handler, defaultHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		version := r.Header.Get(HeaderVersion)
		if version == "" {
			version = r.URL.Query().Get(QueryVersion)
		}
		if version == "" {
			version = DefaultVersion
		}

		if handler, ok := versions[version]; ok {
			handler.ServeHTTP(w, r)
			return
		}

		defaultHandler.ServeHTTP(w, r)
	})
}
