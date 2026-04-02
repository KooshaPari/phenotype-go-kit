package cors

import (
	"net/http"
	"strings"
)

// Config holds CORS configuration.
type Config struct {
	AllowedOrigins   []string `default:"[]"`
	AllowedMethods   []string `default:"[GET,POST,PUT,DELETE,PATCH,OPTIONS]"`
	AllowedHeaders   []string `default:"[Content-Type,Authorization,X-API-Key,X-Request-ID]"`
	ExposedHeaders   []string `default:"[X-RateLimit-Remaining,X-RateLimit-Reset,X-Trace-ID]"`
	AllowCredentials bool     `default:"true"`
	MaxAge           int      `default:"86400"`
}

// Middleware creates CORS middleware.
func Middleware(cfg Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if origin != "" && isOriginAllowed(origin, cfg.AllowedOrigins) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}

			w.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.AllowedMethods, ","))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.AllowedHeaders, ","))
			w.Header().Set("Access-Control-Expose-Headers", strings.Join(cfg.ExposedHeaders, ","))
			w.Header().Set("Access-Control-Max-Age", string(rune(cfg.MaxAge)))

			if cfg.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func isOriginAllowed(origin string, allowed []string) bool {
	// Allow all if no origins specified
	if len(allowed) == 0 {
		return true
	}

	for _, o := range allowed {
		if o == "*" || o == origin {
			return true
		}
	}

	return false
}

// PreflightHandler handles OPTIONS requests.
func PreflightHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}
