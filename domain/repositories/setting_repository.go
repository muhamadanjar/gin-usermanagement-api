package repositories

import (
	"usermanagement-api/domain/entities"
)

type SettingRepository interface {
	Create(setting *entities.Setting) error
	FindByKey(key string) (*entities.Setting, error)
	FindAll() ([]*entities.Setting, error)
	Update(setting *entities.Setting) error
	Delete(key string) error
	Upsert(setting *entities.Setting) error
}
