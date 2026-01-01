package telemetry

import (
	"net/http"
	"project_template/pkg/webserver"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// NewHTTPMiddleware returns the tracing middleware as a provider for fx injection.
func NewHTTPMiddleware(_ *Telemetry) webserver.Middleware {
	return HTTPMiddleware
}

// NewHTTPMetricsMiddleware returns the metrics middleware as a provider for fx injection.
func NewHTTPMetricsMiddleware(_ *Telemetry) webserver.Middleware {
	return HTTPMetricsMiddleware
}

// HTTPMiddleware creates a middleware that traces HTTP requests
func HTTPMiddleware(next http.Handler) http.Handler {
	tracer := otel.Tracer("http-server")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract trace context from incoming request
		ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

		// Start a new span
		spanName := r.Method + " " + r.URL.Path
		ctx, span := tracer.Start(ctx, spanName,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("http.request.method", r.Method),
				attribute.String("url.path", r.URL.Path),
				attribute.String("url.scheme", r.URL.Scheme),
				attribute.String("server.address", r.Host),
				attribute.String("user_agent.original", r.UserAgent()),
				attribute.String("client.address", r.RemoteAddr),
			),
		)
		defer span.End()

		// Wrap response writer to capture status code
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler with the traced context
		next.ServeHTTP(rw, r.WithContext(ctx))

		// Record response attributes
		span.SetAttributes(attribute.Int("http.response.status_code", rw.statusCode))

		// Mark span as error if status code indicates failure
		if rw.statusCode >= 400 {
			span.SetAttributes(attribute.Bool("error", true))
		}
	})
}

// responseWriter wraps http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
