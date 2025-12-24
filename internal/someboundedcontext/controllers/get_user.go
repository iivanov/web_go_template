package controller

import (
	"encoding/json/v2"
	"log/slog"
	"net/http"

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

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		http.Error(w, "user id required", http.StatusBadRequest)

		return
	}

	user, err := h.userService.GetUser(r.Context(), id)
	if err != nil {
		h.logger.Error("failed to get user", "error", err)
		http.Error(w, "failed to get user", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.MarshalWrite(w, user); err != nil {
		h.logger.Error("failed to encode user", "error", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
