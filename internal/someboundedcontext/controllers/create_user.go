package controller

import (
	"encoding/json/v2"
	"log/slog"
	"net/http"

	apperrors "project_template/internal/shared/errors"
	"project_template/internal/someboundedcontext/dto"
	"project_template/internal/someboundedcontext/services"
	"project_template/pkg/validation"
)

type CreateUserHandler struct {
	logger      *slog.Logger
	userService *services.UserService
	validator   *validation.Validator
}

func NewCreateUserHandler(logger *slog.Logger, userService *services.UserService, validator *validation.Validator) *CreateUserHandler {
	return &CreateUserHandler{
		logger:      logger,
		userService: userService,
		validator:   validator,
	}
}

func (*CreateUserHandler) Pattern() string {
	return "POST /api/users"
}

func (h *CreateUserHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	var user dto.CreateUserRequest
	if err := json.UnmarshalRead(r.Body, &user); err != nil {
		return apperrors.NewBadRequest("invalid request body")
	}

	if err := h.validator.Validate(user); err != nil {
		return err
	}

	createdUser, err := h.userService.CreateUser(r.Context(), user)
	if err != nil {
		h.logger.Error("failed to create user", "error", err)
		return apperrors.NewInternalError("failed to create user")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.MarshalWrite(w, createdUser)
}
