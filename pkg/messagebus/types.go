package messagebus

import "context"

// Event represents a domain event that can be published and consumed.
type Event interface {
	// Topic returns the topic/channel name for this event.
	Topic() string
}

// Publisher publishes events to the message bus.
type Publisher interface {
	// Publish publishes an event to its topic.
	Publish(ctx context.Context, event Event) error
	// Close closes the publisher.
	Close() error
}

// Handler processes events from a specific topic.
type Handler interface {
	// Topic returns the topic this handler subscribes to.
	Topic() string
	// Handle processes the event payload.
	Handle(ctx context.Context, payload []byte) error
}

// HandlerFunc is a function type that implements Handler.
type HandlerFunc struct {
	topic    string
	handleFn func(ctx context.Context, payload []byte) error
}

// NewHandlerFunc creates a new HandlerFunc.
func NewHandlerFunc(topic string, fn func(ctx context.Context, payload []byte) error) *HandlerFunc {
	return &HandlerFunc{
		topic:    topic,
		handleFn: fn,
	}
}

func (h *HandlerFunc) Topic() string {
	return h.topic
}

func (h *HandlerFunc) Handle(ctx context.Context, payload []byte) error {
	return h.handleFn(ctx, payload)
}

// Subscriber manages subscriptions and runs handlers.
type Subscriber interface {
	// Subscribe registers a handler for its topic.
	Subscribe(handler Handler) error
	// Run starts processing messages. Blocks until context is cancelled.
	Run(ctx context.Context) error
	// Close closes the subscriber.
	Close() error
}

// MessageBus combines Publisher and Subscriber capabilities.
type MessageBus interface {
	Publisher
	Subscriber
}
