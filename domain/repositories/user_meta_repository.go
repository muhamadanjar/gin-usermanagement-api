// internal/domain/repositories/user_meta_repository.go
package repositories

import (
	"usermanagement-api/domain/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserMetaRepository interface {
	Create(userMeta *entities.UserMeta) error
	FindByUserID(userID uuid.UUID) ([]*entities.UserMeta, error)
	FindByUserIDAndKey(userID uuid.UUID, key string) (*entities.UserMeta, error)
	Update(userMeta *entities.UserMeta) error
	Delete(id uint) error
	GetAllByUserID(userID uuid.UUID) (map[string]string, error)
}

type userMetaRepository struct {
	db *gorm.DB
}

func NewUserMetaRepository(db *gorm.DB) UserMetaRepository {
	return &userMetaRepository{db}
}

func (r *userMetaRepository) Create(userMeta *entities.UserMeta) error {
	return r.db.Create(userMeta).Error
}

func (r *userMetaRepository) FindByUserID(userID uuid.UUID) ([]*entities.UserMeta, error) {
	var userMetas []*entities.UserMeta
	if err := r.db.Where("user_id = ?", userID).Find(&userMetas).Error; err != nil {
		return nil, err
	}
	return userMetas, nil
}

func (r *userMetaRepository) FindByUserIDAndKey(userID uuid.UUID, key string) (*entities.UserMeta, error) {
	var userMeta entities.UserMeta
	if err := r.db.Where("user_id = ? AND key = ?", userID, key).First(&userMeta).Error; err != nil {
		return nil, err
	}
	return &userMeta, nil
}

func (r *userMetaRepository) Update(userMeta *entities.UserMeta) error {
	return r.db.Save(userMeta).Error
}

func (r *userMetaRepository) Delete(id uint) error {
	return r.db.Delete(&entities.UserMeta{}, id).Error
}

func (r *userMetaRepository) GetAllByUserID(userID uuid.UUID) (map[string]string, error) {
	var userMetas []*entities.UserMeta
	if err := r.db.Where("user_id = ?", userID).Find(&userMetas).Error; err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, meta := range userMetas {
		result[meta.Key] = meta.Value
	}
	return result, nil
}
