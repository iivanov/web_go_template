package dto

import (
	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type UsersListResponse []UserResponse
