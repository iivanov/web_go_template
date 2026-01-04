package controller

import (
	"encoding/json/v2"
	"log/slog"
	"net/http"

	apperrors "project_template/internal/shared/errors"
	"project_template/internal/someboundedcontext/services"
)

type UsersHandler struct {
	logger      *slog.Logger
	userService *services.UserService
}

func NewUsersHandler(logger *slog.Logger, userService *services.UserService) *UsersHandler {
	return &UsersHandler{
		logger:      logger,
		userService: userService,
	}
}

func (*UsersHandler) Pattern() string {
	return "GET /api/users"
}

func (h *UsersHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	users, err := h.userService.ListUsers(r.Context())
	if err != nil {
		h.logger.Error("failed to list users", "error", err)
		return apperrors.NewInternalError("failed to list users")
	}

	w.Header().Set("Content-Type", "application/json")
	return json.MarshalWrite(w, users)
}
