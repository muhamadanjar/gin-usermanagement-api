package repositories

import (
	"usermanagement-api/domain/entities"

	"github.com/google/uuid"
)

type PermissionRepository interface {
	Create(permission *entities.Permission) error
	FindByID(id uuid.UUID) (*entities.Permission, error)
	FindByName(name string) (*entities.Permission, error)
	FindAll(page, pageSize int) ([]*entities.Permission, int64, error)
	Update(permission *entities.Permission) error
	Delete(id uuid.UUID) error
}
