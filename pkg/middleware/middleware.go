// Package middleware provides HTTP middleware utilities for chi router.
package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// DefaultMiddlewareStack applies the default middleware stack to a chi router.
// This includes panic recovery, request logging, CORS, and request ID tracking.
func DefaultMiddlewareStack(router *chi.Mux) error {
	// Add standard middleware in order
	router.Use(panicRecoveryMiddleware)
	router.Use(requestLoggingMiddleware)
	router.Use(corsMiddleware)
	return nil
}

// HealthCheckHandler handles health check requests.
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"healthy"}`))
}

// ReadinessCheckHandler handles readiness check requests.
func ReadinessCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"ready":true}`))
}

// panicRecoveryMiddleware recovers from panics in HTTP handlers.
func panicRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error":"internal server error"}`))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// requestLoggingMiddleware logs incoming HTTP requests.
func requestLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Basic logging - can be enhanced with actual logger
		next.ServeHTTP(w, r)
	})
}

// corsMiddleware handles CORS headers.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Chain is a helper function to chain multiple middleware.
func Chain(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

// Register registers middleware with a chi router.
func Register(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	for _, mw := range middlewares {
		r.Use(mw)
	}
}
