package telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// StartSpan starts a new span with the given name and returns the context and span
func StartSpan(ctx context.Context, tracerName, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	tracer := otel.Tracer(tracerName)
	return tracer.Start(ctx, spanName, opts...)
}

// StartServiceSpan starts a span for service layer operations
func StartServiceSpan(ctx context.Context, serviceName, operation string) (context.Context, trace.Span) {
	return StartSpan(ctx, serviceName, operation,
		trace.WithSpanKind(trace.SpanKindInternal),
		trace.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("service.operation", operation),
		),
	)
}

// StartRepositorySpan starts a span for repository/database operations
func StartRepositorySpan(ctx context.Context, repoName, operation string) (context.Context, trace.Span) {
	return StartSpan(ctx, repoName, operation,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("db.system", "postgresql"),
			attribute.String("db.operation", operation),
			attribute.String("repository.name", repoName),
		),
	)
}

// RecordError records an error on the span and sets the status to error
func RecordError(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
}

// SetSpanAttributes sets attributes on the span
func SetSpanAttributes(span trace.Span, attrs ...attribute.KeyValue) {
	span.SetAttributes(attrs...)
}

// SpanFromContext returns the current span from the context
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}
