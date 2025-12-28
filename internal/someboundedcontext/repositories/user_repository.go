package repositories

import (
	"context"
	"errors"

	"project_template/internal/someboundedcontext/entities"

	"github.com/google/uuid"
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
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]*entities.User, error) {
	var users []*entities.User
	err := r.db.WithContext(ctx).Find(&users).Error
	return users, err
}

func (r *UserRepository) Update(ctx context.Context, user *entities.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.User{}, "id = ?", id).Error
}
