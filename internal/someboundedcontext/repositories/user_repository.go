package repositories

import (
	"context"
	"errors"

	"project_template/internal/someboundedcontext/entities"
	"project_template/pkg/telemetry"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *entities.User) error {
	ctx, span := telemetry.StartRepositorySpan(ctx, "UserRepository", "Create")
	defer span.End()

	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		telemetry.RecordError(span, err)
	}
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	ctx, span := telemetry.StartRepositorySpan(ctx, "UserRepository", "GetByID")
	defer span.End()
	span.SetAttributes(attribute.String("user.id", id.String()))

	var user entities.User
	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		telemetry.RecordError(span, err)
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]*entities.User, error) {
	ctx, span := telemetry.StartRepositorySpan(ctx, "UserRepository", "GetAll")
	defer span.End()

	var users []*entities.User
	err := r.db.WithContext(ctx).Find(&users).Error
	if err != nil {
		telemetry.RecordError(span, err)
	}
	span.SetAttributes(attribute.Int("users.count", len(users)))
	return users, err
}

func (r *UserRepository) Update(ctx context.Context, user *entities.User) error {
	ctx, span := telemetry.StartRepositorySpan(ctx, "UserRepository", "Update")
	defer span.End()
	span.SetAttributes(attribute.String("user.id", user.ID.String()))

	err := r.db.WithContext(ctx).Save(user).Error
	if err != nil {
		telemetry.RecordError(span, err)
	}
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, span := telemetry.StartRepositorySpan(ctx, "UserRepository", "Delete")
	defer span.End()
	span.SetAttributes(attribute.String("user.id", id.String()))

	err := r.db.WithContext(ctx).Delete(&entities.User{}, "id = ?", id).Error
	if err != nil {
		telemetry.RecordError(span, err)
	}
	return err
}
