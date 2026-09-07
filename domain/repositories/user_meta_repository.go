package repositories

import (
	"usermanagement-api/domain/entities"

	"github.com/google/uuid"
)

type UserMetaRepository interface {
	Create(userMeta *entities.UserMeta) error
	FindByUserID(userID uuid.UUID) ([]*entities.UserMeta, error)
	FindByUserIDAndKey(userID uuid.UUID, key string) (*entities.UserMeta, error)
	FindByKey(key string) ([]*entities.UserMeta, error)
	Update(userMeta *entities.UserMeta) error
	Delete(id uint) error
	GetAllByUserID(userID uuid.UUID) (map[string]string, error)
}
