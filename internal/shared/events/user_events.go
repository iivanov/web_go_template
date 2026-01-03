package events

import "github.com/google/uuid"

const TopicUserCreated = "user.created"

// UserCreatedEvent is published when a new user is created.
type UserCreatedEvent struct {
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
	Email  string    `json:"email"`
}

func (e UserCreatedEvent) Topic() string {
	return TopicUserCreated
}
