package telemetry

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Telemetry struct {
	TracerProvider *sdktrace.TracerProvider
	MeterProvider  *sdkmetric.MeterProvider
	Tracer         trace.Tracer
	Meter          metric.Meter
	config         Config
	logger         *slog.Logger
}

func NewTelemetry(lc fx.Lifecycle, cfg Config, logger *slog.Logger) (*Telemetry, error) {
	if !cfg.Enabled {
		logger.Info("Telemetry disabled")
		return &Telemetry{
			Tracer: otel.Tracer(cfg.ServiceName),
			Meter:  otel.Meter(cfg.ServiceName),
			config: cfg,
			logger: logger,
		}, nil
	}

	ctx := context.Background()

	// Create OTLP exporter
	var opts []otlptracegrpc.Option
	opts = append(opts, otlptracegrpc.WithEndpoint(cfg.OTLPEndpoint))
	if cfg.Insecure {
		opts = append(opts, otlptracegrpc.WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())))
		opts = append(opts, otlptracegrpc.WithInsecure())
	}

	exporter, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	// Create resource without merging to avoid schema URL conflicts
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
		),
		resource.WithHost(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// Create metrics exporter
	var metricOpts []otlpmetricgrpc.Option
	metricOpts = append(metricOpts, otlpmetricgrpc.WithEndpoint(cfg.OTLPEndpoint))
	if cfg.Insecure {
		metricOpts = append(metricOpts, otlpmetricgrpc.WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())))
		metricOpts = append(metricOpts, otlpmetricgrpc.WithInsecure())
	}

	metricExporter, err := otlpmetricgrpc.New(ctx, metricOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP metric exporter: %w", err)
	}

	// Create MeterProvider
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter,
			sdkmetric.WithInterval(15*time.Second),
		)),
	)

	// Set global TracerProvider, MeterProvider and propagator
	otel.SetTracerProvider(tp)
	otel.SetMeterProvider(mp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	t := &Telemetry{
		TracerProvider: tp,
		MeterProvider:  mp,
		Tracer:         tp.Tracer(cfg.ServiceName),
		Meter:          mp.Meter(cfg.ServiceName),
		config:         cfg,
		logger:         logger,
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down telemetry")
			if err := mp.Shutdown(ctx); err != nil {
				logger.Error("Failed to shutdown MeterProvider", "error", err)
			}
			return tp.Shutdown(ctx)
		},
	})

	logger.Info("Telemetry initialized", "endpoint", cfg.OTLPEndpoint, "service", cfg.ServiceName)
	return t, nil
}

// Tracer returns a tracer for the given instrumentation name
func (t *Telemetry) GetTracer(name string) trace.Tracer {
	if t.TracerProvider != nil {
		return t.TracerProvider.Tracer(name)
	}
	return otel.Tracer(name)
}
