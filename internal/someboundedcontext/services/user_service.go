package services

import (
	"context"
	"errors"
	"log/slog"
	"project_template/internal/someboundedcontext/config"
	"project_template/internal/someboundedcontext/dto"
	"project_template/internal/someboundedcontext/entities"
	"project_template/internal/someboundedcontext/repositories"
	repoErrors "project_template/internal/someboundedcontext/repositories"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserService struct {
	logger     *slog.Logger
	config     config.Config
	repository *repositories.UserRepository
}

func NewUserService(logger *slog.Logger, config config.Config, repository *repositories.UserRepository) *UserService {
	return &UserService{
		logger:     logger,
		config:     config,
		repository: repository,
	}
}

func (s *UserService) GetUser(ctx context.Context, id string) (dto.UserResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return dto.UserResponse{}, ErrUserNotFound
	}

	user, err := s.repository.GetByID(ctx, uid)
	if err != nil {
		if errors.Is(err, repoErrors.ErrUserNotFound) {
			return dto.UserResponse{}, ErrUserNotFound
		}
		s.logger.Error("failed to get user", "error", err, "id", uid)
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *UserService) ListUsers(ctx context.Context) (dto.UsersListResponse, error) {
	users, err := s.repository.GetAll(ctx)
	if err != nil {
		s.logger.Error("failed to list users", "error", err)
		return nil, err
	}

	response := make(dto.UsersListResponse, len(users))
	for i, user := range users {
		response[i] = dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		}
	}

	return response, nil
}

func (s *UserService) CreateUser(ctx context.Context, user dto.CreateUserRequest) (dto.UserResponse, error) {
	newUser := &entities.User{
		Name:  user.Name,
		Email: user.Email,
	}

	err := s.repository.Create(ctx, newUser)
	if err != nil {
		s.logger.Error("failed to create user", "error", err)
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:    newUser.ID,
		Name:  newUser.Name,
		Email: newUser.Email,
	}, nil
}
