package controller

import (
	"encoding/json/v2"
	"errors"
	"log/slog"
	"net/http"

	apperrors "project_template/internal/shared/errors"
	"project_template/internal/someboundedcontext/services"
)

// UserHandler handles GET /api/users/{id}
type UserHandler struct {
	logger      *slog.Logger
	userService *services.UserService
}

func NewUserHandler(logger *slog.Logger, userService *services.UserService) *UserHandler {
	return &UserHandler{
		logger:      logger,
		userService: userService,
	}
}

func (*UserHandler) Pattern() string {
	return "GET /api/users/{id}"
}

func (h *UserHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	if id == "" {
		return apperrors.NewBadRequest("user id required")
	}

	user, err := h.userService.GetUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return apperrors.NewNotFound("user")
		}
		h.logger.Error("failed to get user", "error", err)
		return apperrors.NewInternalError("failed to get user")
	}

	w.Header().Set("Content-Type", "application/json")
	return json.MarshalWrite(w, user)
}
