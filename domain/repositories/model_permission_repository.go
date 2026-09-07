package repositories

import (
	"usermanagement-api/domain/entities"

	"github.com/google/uuid"
)

type ModelPermissionRepository interface {
	Create(modelPermission *entities.ModelPermission) error
	FindByID(id uuid.UUID) (*entities.ModelPermission, error)
	FindByModelTypeAndModelID(modelType string, modelID uuid.UUID) ([]*entities.ModelPermission, error)
	FindAll(page, pageSize int) ([]*entities.ModelPermission, int64, error)
	Update(modelPermission *entities.ModelPermission) error
	Delete(id uuid.UUID) error
	CheckPermission(modelType string, modelID uuid.UUID, permissionID uuid.UUID) (bool, error)
}
