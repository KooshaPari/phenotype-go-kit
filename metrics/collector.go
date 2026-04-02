package metrics

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	// Metric namespace for the application
	Namespace = "phenotype"
	// Subsystem for HTTP metrics
	SubsystemHTTP = "http"
	// Subsystem for business metrics
	SubsystemBusiness = "business"
	// Subsystem for system metrics
	SubsystemSystem = "system"
)

// Metrics holds all application metrics.
type Metrics struct {
	httpRequestCount    *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
	httpResponseSize    *prometheus.HistogramVec

	jobQueueDepth     *prometheus.GaugeVec
	jobProcessingTime *prometheus.HistogramVec
	jobRetries        *prometheus.CounterVec

	dbQueryDuration *prometheus.HistogramVec
	dbQueryErrors   *prometheus.CounterVec

	businessMetrics map[string]*prometheus.CounterVec

	mu sync.RWMutex
}

// NewMetrics creates a new metrics instance.
func NewMetrics() *Metrics {
	m := &Metrics{
		businessMetrics: make(map[string]*prometheus.CounterVec),
	}

	// HTTP metrics
	m.httpRequestCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Subsystem: SubsystemHTTP,
			Name:      "requests_total",
			Help:      "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	m.httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: Namespace,
			Subsystem: SubsystemHTTP,
			Name:      "request_duration_seconds",
			Help:      "HTTP request duration in seconds",
			Buckets:   []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5},
		},
		[]string{"method", "path"},
	)

	m.httpResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: Namespace,
			Subsystem: SubsystemHTTP,
			Name:      "response_size_bytes",
			Help:      "HTTP response size in bytes",
			Buckets:   []float64{100, 1000, 10000, 100000, 1000000},
		},
		[]string{"method", "path"},
	)

	// Job queue metrics
	m.jobQueueDepth = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "job_queue_depth",
			Help:      "Current depth of job queue",
		},
		[]string{"job_type"},
	)

	m.jobProcessingTime = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: Namespace,
			Subsystem: "jobs",
			Name:      "processing_duration_seconds",
			Help:      "Job processing duration in seconds",
			Buckets:   []float64{.1, .5, 1, 5, 10, 30, 60},
		},
		[]string{"job_type", "status"},
	)

	m.jobRetries = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Subsystem: "jobs",
			Name:      "retries_total",
			Help:      "Total number of job retries",
		},
		[]string{"job_type", "attempt"},
	)

	// Database metrics
	m.dbQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: Namespace,
			Subsystem: "db",
			Name:      "query_duration_seconds",
			Help:      "Database query duration in seconds",
			Buckets:   []float64{.001, .005, .01, .05, .1, .5, 1, 5},
		},
		[]string{"query_type", "table"},
	)

	m.dbQueryErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Subsystem: "db",
			Name:      "query_errors_total",
			Help:      "Total number of database query errors",
		},
		[]string{"query_type", "table", "error_code"},
	)

	return m
}

// RecordHTTPRequest records an HTTP request.
func (m *Metrics) RecordHTTPRequest(method, path string, status int, duration time.Duration, size int64) {
	m.httpRequestCount.WithLabelValues(method, path, fmt.Sprintf("%d", status)).Inc()
	m.httpRequestDuration.WithLabelValues(method, path).Observe(duration.Seconds())
	m.httpResponseSize.WithLabelValues(method, path).Observe(float64(size))
}

// RecordJobProcessing records job processing metrics.
func (m *Metrics) RecordJobProcessing(jobType, status string, duration time.Duration) {
	m.jobProcessingTime.WithLabelValues(jobType, status).Observe(duration.Seconds())
}

// RecordJobRetry records a job retry.
func (m *Metrics) RecordJobRetry(jobType string, attempt int) {
	m.jobRetries.WithLabelValues(jobType, fmt.Sprintf("%d", attempt)).Inc()
}

// RecordJobQueueDepth records the current job queue depth.
func (m *Metrics) RecordJobQueueDepth(jobType string, depth int) {
	m.jobQueueDepth.WithLabelValues(jobType).Set(float64(depth))
}

// RecordDBQuery records database query metrics.
func (m *Metrics) RecordDBQuery(queryType, table string, duration time.Duration, err error) {
	m.dbQueryDuration.WithLabelValues(queryType, table).Observe(duration.Seconds())
	if err != nil {
		m.dbQueryErrors.WithLabelValues(queryType, table, "error").Inc()
	}
}

// RecordBusinessMetric records a custom business metric.
func (m *Metrics) RecordBusinessMetric(name string, value int64, labels map[string]string) {
	m.mu.RLock()
	vec, ok := m.businessMetrics[name]
	m.mu.RUnlock()

	labelValues := make([]string, 0, len(labels))
	labelNames := make([]string, 0, len(labels))
	for k := range labels {
		labelNames = append(labelNames, k)
		labelValues = append(labelValues, labels[k])
	}

	if !ok {
		vec = promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Subsystem: SubsystemBusiness,
				Name:      name,
				Help:      fmt.Sprintf("Business metric: %s", name),
			},
			labelNames,
		)
		m.mu.Lock()
		m.businessMetrics[name] = vec
		m.mu.Unlock()
	}

	vec.WithLabelValues(labelValues...).Add(float64(value))
}

// MetricsMiddleware returns HTTP middleware for recording metrics.
func MetricsMiddleware(m *Metrics) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Use WrapHeader to get status code
			wrapped := &statusRecorder{ResponseWriter: w, code: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)
			path := r.URL.Path

			// Don't record metrics for health check endpoints
			if path != "/health" && path != "/ready" {
				m.RecordHTTPRequest(r.Method, path, wrapped.code, duration, 0)
			}
		})
	}
}

type statusRecorder struct {
	http.ResponseWriter
	code int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.code = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if r.code == 0 {
		r.code = http.StatusOK
	}
	return r.ResponseWriter.Write(b)
}
