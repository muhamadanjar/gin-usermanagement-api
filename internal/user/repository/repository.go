package repository

import (
	"context"
	"usermanagement-api/internal/user/domain"

	"gorm.io/gorm"
)

type UserRepository interface {
	RegisterUser(ctx context.Context, tx *gorm.DB, user domain.User) (domain.User, error)
	GetUserById(ctx context.Context, tx *gorm.DB, userId string) (domain.User, error)
	GetUserByEmail(ctx context.Context, tx *gorm.DB, email string) (domain.User, error)
	CheckEmail(ctx context.Context, tx *gorm.DB, email string) (domain.User, bool, error)
	UpdateUser(ctx context.Context, tx *gorm.DB, user domain.User) (domain.User, error)
	DeleteUser(ctx context.Context, tx *gorm.DB, userId string) error
}
