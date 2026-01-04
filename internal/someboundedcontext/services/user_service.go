package services

import (
	"context"
	"errors"
	"log/slog"
	"project_template/internal/shared/events"
	"project_template/internal/someboundedcontext/config"
	"project_template/internal/someboundedcontext/dto"
	"project_template/internal/someboundedcontext/entities"
	"project_template/internal/someboundedcontext/repositories"
	"project_template/pkg/messagebus"
	"project_template/pkg/telemetry"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserService struct {
	logger     *slog.Logger
	config     config.Config
	repository *repositories.UserRepository
	publisher  messagebus.Publisher
}

func NewUserService(logger *slog.Logger, config config.Config, repository *repositories.UserRepository, publisher messagebus.Publisher) *UserService {
	return &UserService{
		logger:     logger,
		config:     config,
		repository: repository,
		publisher:  publisher,
	}
}

func (s *UserService) GetUser(ctx context.Context, id string) (dto.UserResponse, error) {
	ctx, span := telemetry.StartServiceSpan(ctx, "UserService", "GetUser")
	defer span.End()
	span.SetAttributes(attribute.String("user.id", id))

	uid, err := uuid.Parse(id)
	if err != nil {
		return dto.UserResponse{}, ErrUserNotFound
	}

	user, err := s.repository.GetByID(ctx, uid)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return dto.UserResponse{}, ErrUserNotFound
		}
		telemetry.RecordError(span, err)
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
	ctx, span := telemetry.StartServiceSpan(ctx, "UserService", "ListUsers")
	defer span.End()

	users, err := s.repository.GetAll(ctx)
	if err != nil {
		telemetry.RecordError(span, err)
		s.logger.Error("failed to list users", "error", err)
		return nil, err
	}

	span.SetAttributes(attribute.Int("users.count", len(users)))
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
	ctx, span := telemetry.StartServiceSpan(ctx, "UserService", "CreateUser")
	defer span.End()
	span.SetAttributes(
		attribute.String("user.name", user.Name),
		attribute.String("user.email", user.Email),
	)

	newUser := &entities.User{
		Name:  user.Name,
		Email: user.Email,
	}

	err := s.repository.Create(ctx, newUser)
	if err != nil {
		telemetry.RecordError(span, err)
		s.logger.Error("failed to create user", "error", err)
		return dto.UserResponse{}, err
	}

	event := events.UserCreatedEvent{
		UserID: newUser.ID,
		Name:   newUser.Name,
		Email:  newUser.Email,
	}
	if err := s.publisher.Publish(ctx, event); err != nil {
		s.logger.Error("failed to publish UserCreatedEvent", "error", err, "user_id", newUser.ID)
	}

	span.SetAttributes(attribute.String("user.id", newUser.ID.String()))
	return dto.UserResponse{
		ID:    newUser.ID,
		Name:  newUser.Name,
		Email: newUser.Email,
	}, nil
}
