package repositories

import (
	"usermanagement-api/domain/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *entities.User) error
	FindByID(id uuid.UUID) (*entities.User, error)
	FindByEmail(email string) (*entities.User, error)
	FindByUsername(username string) (*entities.User, error)
	FindAll(page, pageSize int) ([]*entities.User, int64, error)
	Update(user *entities.User) error
	Delete(id uuid.UUID) error
	AssignRoles(userID uuid.UUID, roleIDs []uuid.UUID) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *entities.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uuid.UUID) (*entities.User, error) {
	var user entities.User
	if err := r.db.Preload("Roles").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*entities.User, error) {
	var user entities.User
	if err := r.db.Preload("Roles").Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(username string) (*entities.User, error) {
	var user entities.User
	if err := r.db.Preload("Roles").Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAll(page, pageSize int) ([]*entities.User, int64, error) {
	var users []*entities.User
	var count int64

	offset := (page - 1) * pageSize

	if err := r.db.Model(&entities.User{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Preload("Roles").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

func (r *userRepository) Update(user *entities.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entities.User{}, id).Error
}

func (r *userRepository) AssignRoles(userID uuid.UUID, roleIDs []uuid.UUID) error {
	tx := r.db.Begin()

	// Remove existing roles
	if err := tx.Model(&entities.User{ID: userID}).Association("Roles").Clear(); err != nil {
		tx.Rollback()
		return err
	}

	// Add new roles
	var roles []*entities.Role
	for _, roleID := range roleIDs {
		roles = append(roles, &entities.Role{ID: roleID})
	}

	if err := tx.Model(&entities.User{ID: userID}).Association("Roles").Append(roles); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
