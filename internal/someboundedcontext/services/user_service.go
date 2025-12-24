package services

import (
	"context"
	"log/slog"
	"project_template/internal/someboundedcontext/config"

	"project_template/internal/someboundedcontext/dto"
)

type UserService struct {
	logger *slog.Logger
	config config.Config
}

func NewUserService(logger *slog.Logger) *UserService {
	return &UserService{
		logger: logger,
	}
}

func (s *UserService) GetUser(_ context.Context, id string) (dto.UserResponse, error) {
	// TODO: fetch user from database
	user := dto.UserResponse{
		UID:   s.config.UIDPrefix + id,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	return user, nil
}

func (s *UserService) ListUsers(_ context.Context) (dto.UsersListResponse, error) {
	// TODO: fetch users from database
	users := dto.UsersListResponse{
		{UID: "some_uid_1", Name: "John Doe", Email: "john@example.com"},
		{UID: "some_uid_1", Name: "Jane Smith", Email: "jane@example.com"},
	}

	return users, nil
}

func (s *UserService) CreateUser(_ context.Context, user dto.CreateUserRequest) (dto.UserResponse, error) {
	// TODO: save user to database

	return dto.UserResponse{
		UID:   "asdf",
		Name:  user.Name,
		Email: "john@example.com",
	}, nil
}
