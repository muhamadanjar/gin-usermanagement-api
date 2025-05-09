package repositories

import (
	"fmt"
	"usermanagement-api/domain/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
}

type menuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{db}
}

func (r *menuRepository) Create(menu *entities.Menu) error {
	return r.db.Create(menu).Error
}

func (r *menuRepository) FindByID(id uuid.UUID) (*entities.Menu, error) {
	var menu entities.Menu
	if err := r.db.First(&menu, id).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *menuRepository) FindByName(name string) (*entities.Menu, error) {
	var menu entities.Menu
	if err := r.db.Where("name = ?", name).First(&menu).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *menuRepository) FindAll(page, pageSize int) ([]*entities.Menu, int64, error) {
	var menus []*entities.Menu
	var count int64

	offset := (page - 1) * pageSize

	if err := r.db.Model(&entities.Menu{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Preload("Children").Offset(offset).Limit(pageSize).Find(&menus).Error; err != nil {
		return nil, 0, err
	}

	return menus, count, nil
}

func (r *menuRepository) FindAllActive() ([]*entities.Menu, error) {
	var menus []*entities.Menu
	// parent_id IS NULL
	if err := r.db.Preload("Children", "active = ?", true).Where("active = ?", true).Order("\"order\" asc").Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

func (r *menuRepository) FindAllByParentID(parentID *uuid.UUID) ([]*entities.Menu, error) {
	var menus []*entities.Menu
	query := r.db.Order("\"order\" asc")

	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}

	if err := query.Find(&menus).Error; err != nil {
		return nil, err
	}

	return menus, nil
}

func (r *menuRepository) Update(menu *entities.Menu) error {
	return r.db.Save(menu).Error
}

func (r *menuRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entities.Menu{}, id).Error
}

func (r *menuRepository) FindMenusByRoleID(roleID uuid.UUID) ([]*entities.Menu, error) {
	var menus []*entities.Menu
	var permissionIDs []uuid.UUID
	perm_err := r.db.Table("role_permissions").
		Select("permission_id").
		Where("role_id = ?", roleID).
		Pluck("permission_id", &permissionIDs).Error
	if perm_err != nil {
		return nil, perm_err
	}

	fmt.Println("perms ids", permissionIDs)

	// This query assumes you have a role_menus table that connects roles to menus
	// You might need to modify this based on your actual database structure
	err := r.db.Table("menus").
		Joins("INNER JOIN model_permissions ON menus.id = CAST(model_permissions.model_id as uuid)").
		Where("model_permissions.model_type = ? AND model_permissions.permission_id IN ? AND menus.is_visible = ? AND menus.is_active = ?", "menu", permissionIDs, true, true).
		Order("menus.sequence ASC").Find(&menus).Error

	return menus, err
}

func (r *menuRepository) MenuBySuperUser() ([]*entities.Menu, error) {
	var menus []*entities.Menu
	if err := r.db.Table("menus").
		Joins("INNER JOIN model_permissions ON menus.id = CAST(model_permissions.model_id as uuid)").
		Where("menus.is_active = ? AND menus.is_visible = ?", true, true).Order("menus.sequence asc").Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}
