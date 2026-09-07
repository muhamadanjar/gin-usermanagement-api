package repositories

import (
	"usermanagement-api/domain/entities"
	domainrepos "usermanagement-api/domain/repositories"
	"usermanagement-api/infrastructure/database/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userMetaRepository struct {
	db *gorm.DB
}

func NewUserMetaRepository(db *gorm.DB) domainrepos.UserMetaRepository {
	return &userMetaRepository{db}
}

func (r *userMetaRepository) Create(userMeta *entities.UserMeta) error {
	m := toModelUserMeta(userMeta)
	if err := r.db.Create(m).Error; err != nil {
		return err
	}
	userMeta.ID = m.ID
	return nil
}

func (r *userMetaRepository) FindByUserID(userID uuid.UUID) ([]*entities.UserMeta, error) {
	var rows []*models.UserMetaModel
	if err := r.db.Where("user_id = ?", userID).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*entities.UserMeta, 0, len(rows))
	for _, m := range rows {
		result = append(result, toEntityUserMeta(m))
	}
	return result, nil
}

func (r *userMetaRepository) FindByUserIDAndKey(userID uuid.UUID, key string) (*entities.UserMeta, error) {
	var m models.UserMetaModel
	if err := r.db.Where("user_id = ? AND key = ?", userID, key).First(&m).Error; err != nil {
		return nil, err
	}
	return toEntityUserMeta(&m), nil
}

func (r *userMetaRepository) Update(userMeta *entities.UserMeta) error {
	return r.db.Save(toModelUserMeta(userMeta)).Error
}

func (r *userMetaRepository) Delete(id uint) error {
	return r.db.Delete(&models.UserMetaModel{}, id).Error
}

func (r *userMetaRepository) GetAllByUserID(userID uuid.UUID) (map[string]string, error) {
	var rows []*models.UserMetaModel
	if err := r.db.Where("user_id = ?", userID).Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, m := range rows {
		result[m.Key] = m.Value
	}
	return result, nil
}

func (r *userMetaRepository) FindByKey(key string) ([]*entities.UserMeta, error) {
	var rows []*models.UserMetaModel
	if err := r.db.Where("key = ?", key).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*entities.UserMeta, 0, len(rows))
	for _, m := range rows {
		result = append(result, toEntityUserMeta(m))
	}
	return result, nil
}
