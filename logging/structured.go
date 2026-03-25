package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
)

// JSONLogSchema defines the structured log format.
type JSONLogSchema struct {
	Timestamp  string                 `json:"timestamp"`
	Level      string                 `json:"level"`
	Message    string                 `json:"message"`
	Caller     string                 `json:"caller,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
	TraceID    string                 `json:"trace_id,omitempty"`
	SpanID     string                 `json:"span_id,omitempty"`
}

// LogConfig holds logging configuration.
type LogConfig struct {
	Level      slog.Level `default:"info"`
	Output     string     `default:"stdout"`
	Directory  string     `default:"logs"`
	MaxSizeMB  int        `default:"100"`
	MaxAgeDays int        `default:"30"`
	MaxBackups int        `default:"10"`
}

// Logger wraps slog with additional functionality.
type Logger struct {
	*slog.Logger
	config     LogConfig
	traceIDKey string
	spanIDKey  string
}

// NewLogger creates a new structured logger.
func NewLogger(cfg LogConfig) *Logger {
	var output io.Writer

	switch cfg.Output {
	case "stdout":
		output = os.Stdout
	case "stderr":
		output = os.Stderr
	default:
		output = os.Stdout
	}

	handler := slog.NewJSONHandler(output, &slog.HandlerOptions{
		Level:     cfg.Level,
		AddSource: true,
	})

	logger := &Logger{
		Logger:     slog.New(handler),
		config:     cfg,
		traceIDKey: "trace_id",
		spanIDKey:  "span_id",
	}

	return logger
}

// WithTraceID returns a logger with trace ID.
func (l *Logger) WithTraceID(traceID string) *Logger {
	return &Logger{
		Logger:     l.With("trace_id", traceID),
		config:     l.config,
		traceIDKey: l.traceIDKey,
		spanIDKey:  l.spanIDKey,
	}
}

// WithAttrs returns a logger with additional attributes.
func (l *Logger) WithAttrs(attrs map[string]interface{}) *Logger {
	var args []any
	for k, v := range attrs {
		args = append(args, k, v)
	}

	return &Logger{
		Logger:     l.With(args...),
		config:     l.config,
		traceIDKey: l.traceIDKey,
		spanIDKey:  l.spanIDKey,
	}
}

// LogJSON marshals the log entry to JSON.
func (l *JSONLogSchema) LogJSON() ([]byte, error) {
	return json.Marshal(l)
}

// ContextKey type for log context values.
type ContextKey string

const (
	TraceIDKey ContextKey = "trace_id"
	SpanIDKey  ContextKey = "span_id"
)

// AddTraceID adds trace ID to context.
func AddTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// AddSpanID adds span ID to context.
func AddSpanID(ctx context.Context, spanID string) context.Context {
	return context.WithValue(ctx, SpanIDKey, spanID)
}

// GetTraceID extracts trace ID from context.
func GetTraceID(ctx context.Context) string {
	if v := ctx.Value(TraceIDKey); v != nil {
		return v.(string)
	}
	return ""
}

// GetSpanID extracts span ID from context.
func GetSpanID(ctx context.Context) string {
	if v := ctx.Value(SpanIDKey); v != nil {
		return v.(string)
	}
	return ""
}

// LogCallerInfo captures caller information.
func LogCallerInfo(depth int) (string, string) {
	pc, file, line, ok := runtime.Caller(depth)
	if !ok {
		return "", ""
	}
	fn := runtime.FuncForPC(pc)
	return fmt.Sprintf("%s:%d", filepath.Base(file), line), fn.Name()
}
