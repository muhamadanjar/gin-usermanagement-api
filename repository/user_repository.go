package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser()
	EditUser()
	ViewUser()
	DeleteUser(id uuid.UUID)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateUser() {

}

func (r *userRepository) EditUser() {

}

func (r *userRepository) ViewUser() {

}

func (r *userRepository) DeleteUser(id uuid.UUID) {

}
