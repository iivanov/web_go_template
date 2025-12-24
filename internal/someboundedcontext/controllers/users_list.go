package controller

import (
	"encoding/json/v2"
	"log/slog"
	"net/http"

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

func (h *UsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.ListUsers(r.Context())
	if err != nil {
		h.logger.Error("failed to list users", "error", err)
		http.Error(w, "failed to list users", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.MarshalWrite(w, users); err != nil {
		h.logger.Error("failed to encode users", "error", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
