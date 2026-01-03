package handlers

import (
	"context"
	"encoding/json/v2"
	"log/slog"

	"project_template/internal/shared/events"
	"project_template/pkg/messagebus"
)

// UserCreatedHandler handles UserCreatedEvent from someboundedcontext.
type UserCreatedHandler struct {
	logger *slog.Logger
}

// NewUserCreatedHandler creates a new UserCreatedHandler.
func NewUserCreatedHandler(logger *slog.Logger) messagebus.Handler {
	return &UserCreatedHandler{
		logger: logger,
	}
}

func (h *UserCreatedHandler) Topic() string {
	return events.TopicUserCreated
}

func (h *UserCreatedHandler) Handle(ctx context.Context, payload []byte) error {
	var event events.UserCreatedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		h.logger.Error("failed to unmarshal UserCreatedEvent", "error", err)
		return err
	}

	h.logger.Info("received UserCreatedEvent",
		"user_id", event.UserID,
		"name", event.Name,
		"email", event.Email,
	)

	return nil
}
