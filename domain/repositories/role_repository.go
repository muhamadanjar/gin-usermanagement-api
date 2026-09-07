package repositories

import (
	"usermanagement-api/domain/entities"

	"github.com/google/uuid"
)

type RoleRepository interface {
	Create(role *entities.Role) error
	FindByID(id uuid.UUID) (*entities.Role, error)
	FindByName(name string) (*entities.Role, error)
	FindAll(page, pageSize int) ([]*entities.Role, int64, error)
	Update(role *entities.Role) error
	Delete(id uuid.UUID) error
	AssignPermissions(roleID uuid.UUID, permissionIDs []uuid.UUID) error
	FindRolesByUserID(userID uuid.UUID) ([]*entities.Role, error)
	FindPermissionsByRoleIDs(roleIDs []uuid.UUID) ([]*entities.Permission, error)
}
