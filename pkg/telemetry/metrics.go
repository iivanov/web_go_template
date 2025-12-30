package telemetry

import (
	"context"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// HTTPMetrics holds HTTP-related metrics instruments
type HTTPMetrics struct {
	requestCounter  metric.Int64Counter
	requestDuration metric.Float64Histogram
	activeRequests  metric.Int64UpDownCounter
	responseSize    metric.Int64Histogram
}

// NewHTTPMetrics creates HTTP metrics instruments
func NewHTTPMetrics(meter metric.Meter) (*HTTPMetrics, error) {
	requestCounter, err := meter.Int64Counter(
		"http.server.request.total",
		metric.WithDescription("Total number of HTTP requests"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return nil, err
	}

	requestDuration, err := meter.Float64Histogram(
		"http.server.request.duration",
		metric.WithDescription("HTTP request duration in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	activeRequests, err := meter.Int64UpDownCounter(
		"http.server.active_requests",
		metric.WithDescription("Number of active HTTP requests"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return nil, err
	}

	responseSize, err := meter.Int64Histogram(
		"http.server.response.size",
		metric.WithDescription("HTTP response size in bytes"),
		metric.WithUnit("By"),
	)
	if err != nil {
		return nil, err
	}

	return &HTTPMetrics{
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
		activeRequests:  activeRequests,
		responseSize:    responseSize,
	}, nil
}

// HTTPMetricsMiddleware creates a middleware that records HTTP metrics
func HTTPMetricsMiddleware(next http.Handler) http.Handler {
	meter := otel.Meter("http-server")
	metrics, err := NewHTTPMetrics(meter)
	if err != nil {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		start := time.Now()

		attrs := []attribute.KeyValue{
			attribute.String("http.method", r.Method),
			attribute.String("http.route", r.URL.Path),
		}

		// Track active requests
		metrics.activeRequests.Add(ctx, 1, metric.WithAttributes(attrs...))
		defer metrics.activeRequests.Add(ctx, -1, metric.WithAttributes(attrs...))

		// Wrap response writer to capture status and size
		rw := &metricsResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		// Record metrics after request completes
		duration := time.Since(start).Seconds()
		statusAttrs := append(attrs, attribute.Int("http.status_code", rw.statusCode))

		metrics.requestCounter.Add(ctx, 1, metric.WithAttributes(statusAttrs...))
		metrics.requestDuration.Record(ctx, duration, metric.WithAttributes(statusAttrs...))
		metrics.responseSize.Record(ctx, int64(rw.bytesWritten), metric.WithAttributes(statusAttrs...))
	})
}

type metricsResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (rw *metricsResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *metricsResponseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.bytesWritten += n
	return n, err
}

// DatabaseMetrics holds database-related metrics instruments
type DatabaseMetrics struct {
	queryCounter  metric.Int64Counter
	queryDuration metric.Float64Histogram
}

// NewDatabaseMetrics creates database metrics instruments
func NewDatabaseMetrics(meter metric.Meter) (*DatabaseMetrics, error) {
	queryCounter, err := meter.Int64Counter(
		"db.query.total",
		metric.WithDescription("Total number of database queries"),
		metric.WithUnit("{query}"),
	)
	if err != nil {
		return nil, err
	}

	queryDuration, err := meter.Float64Histogram(
		"db.query.duration",
		metric.WithDescription("Database query duration in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	return &DatabaseMetrics{
		queryCounter:  queryCounter,
		queryDuration: queryDuration,
	}, nil
}

// RecordQuery records a database query metric
func (m *DatabaseMetrics) RecordQuery(ctx context.Context, operation string, duration time.Duration, err error) {
	attrs := []attribute.KeyValue{
		attribute.String("db.operation", operation),
		attribute.Bool("db.success", err == nil),
	}
	m.queryCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
	m.queryDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
}

// BusinessMetrics holds business-related metrics instruments
type BusinessMetrics struct {
	meter metric.Meter
}

// NewBusinessMetrics creates business metrics instruments
func NewBusinessMetrics(meter metric.Meter) *BusinessMetrics {
	return &BusinessMetrics{meter: meter}
}

// RecordUserCreated records a user creation event
func (m *BusinessMetrics) RecordUserCreated(ctx context.Context) {
	counter, err := m.meter.Int64Counter(
		"app.users.created.total",
		metric.WithDescription("Total number of users created"),
		metric.WithUnit("{user}"),
	)
	if err != nil {
		return
	}
	counter.Add(ctx, 1)
}

// RecordUserFetched records a user fetch event
func (m *BusinessMetrics) RecordUserFetched(ctx context.Context) {
	counter, err := m.meter.Int64Counter(
		"app.users.fetched.total",
		metric.WithDescription("Total number of user fetch operations"),
		metric.WithUnit("{operation}"),
	)
	if err != nil {
		return
	}
	counter.Add(ctx, 1)
}
