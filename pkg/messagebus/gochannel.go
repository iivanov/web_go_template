package messagebus

import (
	"context"
	"encoding/json/v2"
	"fmt"
	"log/slog"
	"sync"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

// goChannelBus implements MessageBus using Watermill's GoChannel (in-memory).
type goChannelBus struct {
	logger   *slog.Logger
	pubSub   *gochannel.GoChannel
	router   *message.Router
	handlers []Handler
	mu       sync.Mutex
	running  bool
}

// NewGoChannelBus creates a new in-memory message bus.
func NewGoChannelBus(logger *slog.Logger) (MessageBus, error) {
	wmLogger := watermill.NewSlogLogger(logger)

	pubSub := gochannel.NewGoChannel(
		gochannel.Config{
			Persistent: true,
		},
		wmLogger,
	)

	router, err := message.NewRouter(message.RouterConfig{}, wmLogger)
	if err != nil {
		return nil, fmt.Errorf("failed to create router: %w", err)
	}

	router.AddMiddleware(
		middleware.CorrelationID,
		middleware.Recoverer,
	)

	return &goChannelBus{
		logger:   logger,
		pubSub:   pubSub,
		router:   router,
		handlers: make([]Handler, 0),
	}, nil
}

func (b *goChannelBus) Publish(ctx context.Context, event Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := message.NewMessage(watermill.NewUUID(), payload)
	msg.SetContext(ctx)

	if err := b.pubSub.Publish(event.Topic(), msg); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	b.logger.Debug("published event", "topic", event.Topic(), "uuid", msg.UUID)
	return nil
}

func (b *goChannelBus) Subscribe(handler Handler) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.running {
		return fmt.Errorf("cannot subscribe after router has started")
	}

	handlerName := fmt.Sprintf("handler_%s_%d", handler.Topic(), len(b.handlers))

	b.router.AddConsumerHandler(
		handlerName,
		handler.Topic(),
		b.pubSub,
		func(msg *message.Message) error {
			ctx := msg.Context()
			if err := handler.Handle(ctx, msg.Payload); err != nil {
				b.logger.Error("handler error", "topic", handler.Topic(), "error", err)
				return err
			}
			return nil
		},
	)

	b.handlers = append(b.handlers, handler)
	b.logger.Debug("subscribed handler", "topic", handler.Topic(), "name", handlerName)
	return nil
}

func (b *goChannelBus) Run(ctx context.Context) error {
	b.mu.Lock()
	b.running = true
	b.mu.Unlock()

	b.logger.Info("starting message bus router")
	return b.router.Run(ctx)
}

func (b *goChannelBus) Close() error {
	if err := b.router.Close(); err != nil {
		return fmt.Errorf("failed to close router: %w", err)
	}
	if err := b.pubSub.Close(); err != nil {
		return fmt.Errorf("failed to close pubsub: %w", err)
	}
	return nil
}
