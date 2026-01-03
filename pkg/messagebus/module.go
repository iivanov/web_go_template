package messagebus

import (
	"context"
	"log/slog"

	"go.uber.org/fx"
)

// HandlerParams collects all registered handlers via fx dependency injection.
type HandlerParams struct {
	fx.In
	Handlers []Handler `group:"messagebus_handlers"`
}

// AsHandler annotates a handler constructor to be collected by the message bus.
func AsHandler(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Handler)),
		fx.ResultTags(`group:"messagebus_handlers"`),
	)
}

// Module provides the message bus dependencies.
var Module = fx.Module("messagebus",
	fx.Provide(
		fx.Annotate(
			NewMessageBus,
			fx.As(new(MessageBus)),
			fx.As(new(Publisher)),
			fx.As(new(Subscriber)),
		),
	),
	fx.Invoke(registerHandlers),
	fx.Invoke(startRouter),
)

// NewMessageBus creates a new MessageBus based on configuration.
func NewMessageBus(logger *slog.Logger, config Config) (MessageBus, error) {
	switch config.Backend {
	case "gochannel", "":
		return NewGoChannelBus(logger)
	default:
		return NewGoChannelBus(logger)
	}
}

// registerHandlers registers all collected handlers with the message bus.
func registerHandlers(bus MessageBus, params HandlerParams) error {
	for _, handler := range params.Handlers {
		if err := bus.Subscribe(handler); err != nil {
			return err
		}
	}
	return nil
}

// startRouter starts the message bus router in a goroutine managed by fx lifecycle.
func startRouter(lc fx.Lifecycle, bus MessageBus, logger *slog.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := bus.Run(context.Background()); err != nil {
					logger.Error("message bus router error", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return bus.Close()
		},
	})
}
