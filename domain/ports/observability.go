// Package ports defines interfaces for observability.
// Following OpenTelemetry and SRE best practices.
package ports

import (
	"context"
)

// LoggerPort defines the interface for structured logging.
// This allows swapping logging implementations (slog, zerolog, zap, etc.).
type LoggerPort interface {
	// Debug logs a debug message.
	Debug(msg string, fields ...Field)

	// Info logs an info message.
	Info(msg string, fields ...Field)

	// Warn logs a warning message.
	Warn(msg string, fields ...Field)

	// Error logs an error message.
	Error(msg string, fields ...Field)

	// With returns a new logger with additional fields.
	With(fields ...Field) LoggerPort
}

// Field represents a key-value log field.
type Field struct {
	Key   string
	Value any
}

// String creates a string field.
func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

// Int creates an int field.
func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

// Int64 creates an int64 field.
func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

// Float64 creates a float64 field.
func Float64(key string, value float64) Field {
	return Field{Key: key, Value: value}
}

// Bool creates a bool field.
func Bool(key string, value bool) Field {
	return Field{Key: key, Value: value}
}

// Error creates an error field.
func Error(err error) Field {
	return Field{Key: "error", Value: err}
}

// Duration creates a duration field.
func Duration(key string, value any) Field {
	return Field{Key: key, Value: value}
}

// Any creates a field with any value.
func Any(key string, value any) Field {
	return Field{Key: key, Value: value}
}

// MetricsPort defines the interface for metrics collection.
// This allows swapping metrics backends (Prometheus, StatsD, etc.).
type MetricsPort interface {
	// Counter records a counter metric.
	Counter(name string, labels map[string]string) CounterMetric

	// Gauge records a gauge metric.
	Gauge(name string, labels map[string]string) GaugeMetric

	// Histogram records a histogram metric.
	Histogram(name string, labels map[string]string) HistogramMetric
}

// CounterMetric is a counter that can be incremented.
type CounterMetric interface {
	// Inc increments the counter by 1.
	Inc()

	// Add adds the given value to the counter.
	Add(v float64)
}

// GaugeMetric is a gauge that can be set.
type GaugeMetric interface {
	// Set sets the gauge to the given value.
	Set(v float64)

	// Inc increments the gauge by 1.
	Inc()

	// Dec decrements the gauge by 1.
	Dec()
}

// HistogramMetric is a histogram that can record observations.
type HistogramMetric interface {
	// Observe records an observation.
	Observe(v float64)
}

// TracerPort defines the interface for distributed tracing.
// This allows swapping tracing implementations (OTEL, Jaeger, Zipkin, etc.).
type TracerPort interface {
	// StartSpan starts a new span.
	StartSpan(ctx context.Context, name string) (context.Context, SpanPort)

	// SpanFromContext returns the current span from context.
	SpanFromContext(ctx context.Context) SpanPort
}

// SpanPort represents a trace span.
type SpanPort interface {
	// End ends the span.
	End()

	// AddEvent adds an event to the span.
	AddEvent(name string, attrs map[string]any)

	// SetStatus sets the span status.
	SetStatus(code int, msg string)

	// RecordError records an error on the span.
	RecordError(err error)

	// SetAttributes sets attributes on the span.
	SetAttributes(attrs map[string]any)
}
