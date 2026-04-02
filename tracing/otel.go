package tracing

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

// Config holds OpenTelemetry configuration.
type Config struct {
	ServiceName    string `default:"phenotype-go-kit"`
	ServiceVersion string `default:"1.0.0"`
	Environment    string `default:"development"`

	TraceExporterEndpoint string  `default:"localhost:4317"`
	TraceExporterInsecure bool    `default:"true"`
	TraceSamplingRate     float64 `default:"0.1"`
}

// Tracer provides tracing functionality.
type Tracer struct {
	provider *sdktrace.TracerProvider
	tracer   trace.Tracer
	config   Config
}

// NewTracer creates a new tracing provider.
func NewTracer(cfg Config) (*Tracer, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.ServiceVersion),
			semconv.DeploymentEnvironmentKey.String(cfg.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(cfg.TraceExporterEndpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	sampler := sdktrace.ParentBased(sdktrace.TraceIDRatioBased(cfg.TraceSamplingRate))
	spanProcessor := sdktrace.NewBatchSpanProcessor(traceExporter)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sampler),
		sdktrace.WithSpanProcessor(spanProcessor),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	t := tp.Tracer(cfg.ServiceName)

	return &Tracer{
		provider: tp,
		tracer:   t,
		config:   cfg,
	}, nil
}

// StartSpan starts a new span.
func (t *Tracer) StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, name, opts...)
}

// StartSpanWithAttributes starts a span with common attributes.
func (t *Tracer) StartSpanWithAttributes(ctx context.Context, name string, attrs map[string]interface{}) (context.Context, trace.Span) {
	sattrs := make([]attribute.KeyValue, 0, len(attrs))
	for k, v := range attrs {
		sattrs = append(sattrs, attribute.String(k, fmt.Sprintf("%v", v)))
	}
	return t.tracer.Start(ctx, name, trace.WithAttributes(sattrs...))
}

// AddEvent adds an event to the current span.
func AddEvent(ctx context.Context, name string, attrs ...attribute.KeyValue) {
	if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
		span.AddEvent(name, trace.WithAttributes(attrs...))
	}
}

// AddError records an error on the span.
func AddError(ctx context.Context, err error) {
	if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
		span.RecordError(err)
	}
}

// GetTraceID returns the trace ID from context.
func GetTraceID(ctx context.Context) string {
	if span := trace.SpanFromContext(ctx); span.SpanContext().HasTraceID() {
		return span.SpanContext().TraceID().String()
	}
	return ""
}

// GetSpanID returns the span ID from context.
func GetSpanID(ctx context.Context) string {
	if span := trace.SpanFromContext(ctx); span.SpanContext().HasSpanID() {
		return span.SpanContext().SpanID().String()
	}
	return ""
}

// InjectTraceHeaders injects trace headers into the carrier.
func (t *Tracer) InjectTraceHeaders(ctx context.Context, carrier propagation.MapCarrier) {
	otel.GetTextMapPropagator().Inject(ctx, carrier)
}

// ExtractTraceHeaders extracts trace headers from carrier into context.
func (t *Tracer) ExtractTraceHeaders(ctx context.Context, carrier propagation.MapCarrier) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, carrier)
}

// Shutdown gracefully shuts down the tracer.
func (t *Tracer) Shutdown(ctx context.Context) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return t.provider.Shutdown(timeoutCtx)
}

// InstrumentHTTP instruments an HTTP handler with tracing.
func InstrumentHTTP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		propagator := otel.GetTextMapPropagator()
		ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

		tracer := otel.Tracer("phenotype")
		ctx, span := tracer.Start(ctx, r.URL.Path,
			trace.WithAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.url", r.URL.String()),
			),
		)
		defer span.End()

		w.Header().Add("X-Trace-ID", span.SpanContext().TraceID().String())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// DBAttributes returns attributes for database tracing.
func DBAttributes(operation, dbName, query string) []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("db.system", "postgresql"),
		attribute.String("db.operation", operation),
		attribute.String("db.name", dbName),
		attribute.String("db.statement", query),
	}
}

// HTTPAttributes returns attributes for HTTP tracing.
func HTTPAttributes(method, URL, status string) []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("http.method", method),
		attribute.String("http.url", URL),
		attribute.String("http.status_code", status),
	}
}
