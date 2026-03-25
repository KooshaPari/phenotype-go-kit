package logging

import (
	"net/http"
	"time"
)

// Interceptor wraps handlers with structured logging.
type Interceptor struct {
	logger *Logger
}

// NewInterceptor creates a new logging interceptor.
func NewInterceptor(logger *Logger) *Interceptor {
	return &Interceptor{logger: logger}
}

// WrapHandler wraps an HTTP handler with logging.
func (i *Interceptor) WrapHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		traceID := r.Header.Get("X-Trace-ID")
		if traceID == "" {
			traceID = generateTraceID()
		}

		ctx := AddTraceID(r.Context(), traceID)
		r = r.WithContext(ctx)

		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		i.logger.Info("incoming request",
			"method", r.Method,
			"path", r.URL.Path,
			"query", r.URL.RawQuery,
			"remote_addr", r.RemoteAddr,
			"trace_id", traceID,
			"user_agent", r.UserAgent(),
		)

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		i.logger.Info("request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status_code", wrapped.statusCode,
			"duration_ms", duration.Milliseconds(),
			"trace_id", traceID,
		)
	})
}

// WrapFunc wraps an HTTP handler function with logging.
func (i *Interceptor) WrapFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		traceID := r.Header.Get("X-Trace-ID")
		if traceID == "" {
			traceID = generateTraceID()
		}

		ctx := AddTraceID(r.Context(), traceID)
		r = r.WithContext(ctx)

		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		i.logger.Info("incoming request",
			"method", r.Method,
			"path", r.URL.Path,
			"trace_id", traceID,
		)

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		i.logger.Info("request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status_code", wrapped.statusCode,
			"duration_ms", duration.Milliseconds(),
			"trace_id", traceID,
		)
	}
}

// MiddlewareFunc returns a middleware function for chi/echo/etc.
func (i *Interceptor) MiddlewareFunc(next http.Handler) http.Handler {
	return i.WrapHandler(next)
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func (w *responseWriter) WriteHeader(code int) {
	if !w.written {
		w.statusCode = code
		w.ResponseWriter.WriteHeader(code)
		w.written = true
	}
}

func (w *responseWriter) Write(b []byte) (int, error) {
	if !w.written {
		w.statusCode = http.StatusOK
		w.written = true
	}
	return w.ResponseWriter.Write(b)
}

// LogInterceptorFunc creates a simple middleware for standard libraries.
func LogInterceptorFunc(logger *Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			traceID := r.Header.Get("X-Trace-ID")
			if traceID == "" {
				traceID = generateTraceID()
			}

			r = r.WithContext(AddTraceID(r.Context(), traceID))
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			logger.Info("request started",
				"method", r.Method,
				"path", r.URL.Path,
				"trace_id", traceID,
			)

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)
			logger.Info("request completed",
				"method", r.Method,
				"path", r.URL.Path,
				"status_code", wrapped.statusCode,
				"duration_ms", duration.Milliseconds(),
				"trace_id", traceID,
			)
		})
	}
}

func generateTraceID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(16)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[int(time.Now().UnixNano())%len(letters)]
	}
	return string(b)
}
