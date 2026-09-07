package repositories

import (
	"usermanagement-api/domain/entities"

	"github.com/google/uuid"
)

type MenuRepository interface {
	Create(menu *entities.Menu) error
	FindByID(id uuid.UUID) (*entities.Menu, error)
	FindByName(name string) (*entities.Menu, error)
	FindAll(page, pageSize int) ([]*entities.Menu, int64, error)
	FindAllActive() ([]*entities.Menu, error)
	FindAllByParentID(parentID *uuid.UUID) ([]*entities.Menu, error)
	Update(menu *entities.Menu) error
	Delete(id uuid.UUID) error
	FindMenusByRoleID(roleID uuid.UUID) ([]*entities.Menu, error)
	MenuBySuperUser() ([]*entities.Menu, error)
	AssignPermissionsByName(menuID uuid.UUID, names []string) error
}
