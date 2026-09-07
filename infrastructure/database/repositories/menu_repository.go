package repositories

import (
	"strings"
	"usermanagement-api/domain/entities"
	domainrepos "usermanagement-api/domain/repositories"
	"usermanagement-api/infrastructure/database/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type menuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) domainrepos.MenuRepository {
	return &menuRepository{db}
}

func (r *menuRepository) Create(menu *entities.Menu) error {
	m := toModelMenu(menu)
	if err := r.db.Create(m).Error; err != nil {
		return err
	}
	menu.ID = m.ID
	menu.CreatedAt = m.CreatedAt
	menu.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *menuRepository) FindByID(id uuid.UUID) (*entities.Menu, error) {
	var m models.MenuModel
	if err := r.db.Preload("Permissions").First(&m, id).Error; err != nil {
		return nil, err
	}
	return toEntityMenu(&m), nil
}

func (r *menuRepository) FindByName(name string) (*entities.Menu, error) {
	var m models.MenuModel
	if err := r.db.Where("name = ?", name).First(&m).Error; err != nil {
		return nil, err
	}
	return toEntityMenu(&m), nil
}

func (r *menuRepository) FindAll(page, pageSize int) ([]*entities.Menu, int64, error) {
	var rows []*models.MenuModel
	var count int64

	offset := (page - 1) * pageSize

	if err := r.db.Model(&models.MenuModel{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Preload("Children").Preload("Permissions").Offset(offset).Limit(pageSize).Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	menus := make([]*entities.Menu, 0, len(rows))
	for _, m := range rows {
		menus = append(menus, toEntityMenu(m))
	}
	return menus, count, nil
}

func (r *menuRepository) FindAllActive() ([]*entities.Menu, error) {
	var rows []*models.MenuModel
	if err := r.db.Preload("Children", "is_active = ?", true).Preload("Permissions").Where("is_active = ?", true).Order("sequence asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	menus := make([]*entities.Menu, 0, len(rows))
	for _, m := range rows {
		menus = append(menus, toEntityMenu(m))
	}
	return menus, nil
}

func (r *menuRepository) FindAllByParentID(parentID *uuid.UUID) ([]*entities.Menu, error) {
	var rows []*models.MenuModel
	query := r.db.Order("sequence asc")

	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}

	if err := query.Find(&rows).Error; err != nil {
		return nil, err
	}

	menus := make([]*entities.Menu, 0, len(rows))
	for _, m := range rows {
		menus = append(menus, toEntityMenu(m))
	}
	return menus, nil
}

func (r *menuRepository) Update(menu *entities.Menu) error {
	return r.db.Save(toModelMenu(menu)).Error
}

func (r *menuRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.MenuModel{}, id).Error
}

// AssignPermissionsByName appends permissions (by name, created if missing) to a
// menu via the menu_permission join table — mirrors FastAPI's
// assign_permissions_by_names.
func (r *menuRepository) AssignPermissionsByName(menuID uuid.UUID, names []string) error {
	var menu models.MenuModel
	if err := r.db.Preload("Permissions").First(&menu, menuID).Error; err != nil {
		return err
	}
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		var perm models.PermissionModel
		err := r.db.Where("name = ?", name).First(&perm).Error
		if err == gorm.ErrRecordNotFound {
			perm = models.PermissionModel{Name: name}
			if err := r.db.Create(&perm).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
		if err := r.db.Model(&menu).Association("Permissions").Append(&perm); err != nil {
			return err
		}
	}
	return nil
}

func (r *menuRepository) FindMenusByRoleID(roleID uuid.UUID) ([]*entities.Menu, error) {
	var rows []*models.MenuModel

	// role → permissions → menus, via the menu_permission join table (mirrors FastAPI).
	permSub := r.db.Table("role_permissions").
		Select("permission_id").
		Where("role_id = ?", roleID)
	permIDs := []uuid.UUID{}
	// Distinct permission ids assigned to the role.
	if err := r.db.Model(&models.PermissionModel{}).
		Where("id IN (?)", permSub).
		Where("is_active = ?", true).
		Distinct("id").
		Pluck("id", &permIDs).Error; err != nil {
		return nil, err
	}
	if len(permIDs) == 0 {
		return []*entities.Menu{}, nil
	}

	err := r.db.Preload("Permissions").
		Joins("JOIN menu_permission mp ON mp.menu_id = menus.id").
		Where("mp.permission_id IN ?", permIDs).
		Where("menus.is_visible = ? AND menus.is_active = ?", true, true).
		Order("menus.sequence ASC").Find(&rows).Error

	menus := make([]*entities.Menu, 0, len(rows))
	for _, m := range rows {
		menus = append(menus, toEntityMenu(m))
	}
	return menus, err
}

func (r *menuRepository) MenuBySuperUser() ([]*entities.Menu, error) {
	var rows []*models.MenuModel
	if err := r.db.Preload("Permissions").
		Where("menus.is_active = ? AND menus.is_visible = ?", true, true).
		Order("menus.sequence asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	menus := make([]*entities.Menu, 0, len(rows))
	for _, m := range rows {
		menus = append(menus, toEntityMenu(m))
	}
	return menus, nil
}
