package repository

import (
	"usermanagement-api/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetAll([]domain.User, error) {

}

func (r *userRepository) CreateUser(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) EditUser() {

}

func (r *userRepository) ViewUser() {

}

func (r *userRepository) DeleteUser(id uuid.UUID) {

}
