package repositories

import (
	"usermanagement-api/domain/entities"
	domainrepos "usermanagement-api/domain/repositories"
	"usermanagement-api/infrastructure/database/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domainrepos.UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *entities.User) error {
	m := toModelUser(user)
	if err := r.db.Create(m).Error; err != nil {
		return err
	}
	user.ID = m.ID
	user.CreatedAt = m.CreatedAt
	user.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *userRepository) FindByID(id uuid.UUID) (*entities.User, error) {
	var m models.UserModel
	if err := r.db.Preload("Roles").First(&m, id).Error; err != nil {
		return nil, err
	}
	return toEntityUser(&m), nil
}

func (r *userRepository) FindByEmail(email string) (*entities.User, error) {
	var m models.UserModel
	if err := r.db.Preload("Roles").Where("email = ?", email).First(&m).Error; err != nil {
		return nil, err
	}
	return toEntityUser(&m), nil
}

func (r *userRepository) FindByUsername(username string) (*entities.User, error) {
	var m models.UserModel
	if err := r.db.Preload("Roles").Where("username = ?", username).First(&m).Error; err != nil {
		return nil, err
	}
	return toEntityUser(&m), nil
}

func (r *userRepository) FindAll(page, pageSize int) ([]*entities.User, int64, error) {
	var rows []*models.UserModel
	var count int64

	offset := (page - 1) * pageSize

	if err := r.db.Model(&models.UserModel{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Preload("Roles").Offset(offset).Limit(pageSize).Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	users := make([]*entities.User, 0, len(rows))
	for _, m := range rows {
		users = append(users, toEntityUser(m))
	}
	return users, count, nil
}

func (r *userRepository) Update(user *entities.User) error {
	return r.db.Save(toModelUser(user)).Error
}

func (r *userRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.UserModel{}, id).Error
}

func (r *userRepository) AssignRoles(userID uuid.UUID, roleIDs []uuid.UUID) error {
	tx := r.db.Begin()

	if err := tx.Model(&models.UserModel{ID: userID}).Association("Roles").Clear(); err != nil {
		tx.Rollback()
		return err
	}

	var roles []*models.RoleModel
	for _, roleID := range roleIDs {
		roles = append(roles, &models.RoleModel{ID: roleID})
	}

	if err := tx.Model(&models.UserModel{ID: userID}).Association("Roles").Append(roles); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
