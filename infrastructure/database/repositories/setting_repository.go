package repositories

import (
	"usermanagement-api/domain/entities"
	domainrepos "usermanagement-api/domain/repositories"
	"usermanagement-api/infrastructure/database/models"

	"gorm.io/gorm"
)

type settingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) domainrepos.SettingRepository {
	return &settingRepository{db}
}

func (r *settingRepository) Create(setting *entities.Setting) error {
	return r.db.Create(toModelSetting(setting)).Error
}

func (r *settingRepository) FindByKey(key string) (*entities.Setting, error) {
	var m models.SettingModel
	if err := r.db.Where("key = ?", key).First(&m).Error; err != nil {
		return nil, err
	}
	return toEntitySetting(&m), nil
}

func (r *settingRepository) FindAll() ([]*entities.Setting, error) {
	var rows []*models.SettingModel
	if err := r.db.Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*entities.Setting, 0, len(rows))
	for _, m := range rows {
		result = append(result, toEntitySetting(m))
	}
	return result, nil
}

func (r *settingRepository) Update(setting *entities.Setting) error {
	return r.db.Save(toModelSetting(setting)).Error
}

func (r *settingRepository) Delete(key string) error {
	return r.db.Where("key = ?", key).Delete(&models.SettingModel{}).Error
}

func (r *settingRepository) Upsert(setting *entities.Setting) error {
	return r.db.Save(toModelSetting(setting)).Error
}
