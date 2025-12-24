package controller

import (
	"encoding/json/v2"
	"log/slog"
	"net/http"
	"project_template/internal/someboundedcontext/dto"

	"project_template/internal/someboundedcontext/services"
)

type CreateUserHandler struct {
	logger      *slog.Logger
	userService *services.UserService
}

func NewCreateUserHandler(logger *slog.Logger, userService *services.UserService) *CreateUserHandler {
	return &CreateUserHandler{
		logger:      logger,
		userService: userService,
	}
}

func (*CreateUserHandler) Pattern() string {
	return "POST /api/users"
}

func (h *CreateUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var user dto.CreateUserRequest
	if err := json.UnmarshalRead(r.Body, &user); err != nil {
		h.logger.Error("failed to decode user", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)

		return
	}

	createdUser, err := h.userService.CreateUser(r.Context(), user)
	if err != nil {
		h.logger.Error("failed to create user", "error", err)
		http.Error(w, "failed to create user", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.MarshalWrite(w, createdUser); err != nil {
		h.logger.Error("failed to encode user", "error", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
