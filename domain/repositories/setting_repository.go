// internal/domain/repositories/setting_repository.go
package repositories

import (
	"usermanagement-api/domain/entities"

	"gorm.io/gorm"
)

type SettingRepository interface {
	Create(setting *entities.Setting) error
	FindByKey(key string) (*entities.Setting, error)
	FindAll() ([]*entities.Setting, error)
	Update(setting *entities.Setting) error
	Delete(key string) error
	Upsert(setting *entities.Setting) error
}

type settingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) SettingRepository {
	return &settingRepository{db}
}

func (r *settingRepository) Create(setting *entities.Setting) error {
	return r.db.Create(setting).Error
}

func (r *settingRepository) FindByKey(key string) (*entities.Setting, error) {
	var setting entities.Setting
	if err := r.db.Where("key = ?", key).First(&setting).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *settingRepository) FindAll() ([]*entities.Setting, error) {
	var settings []*entities.Setting
	if err := r.db.Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}

func (r *settingRepository) Update(setting *entities.Setting) error {
	return r.db.Save(setting).Error
}

func (r *settingRepository) Delete(key string) error {
	return r.db.Where("key = ?", key).Delete(&entities.Setting{}).Error
}

func (r *settingRepository) Upsert(setting *entities.Setting) error {
	return r.db.Save(setting).Error
}
