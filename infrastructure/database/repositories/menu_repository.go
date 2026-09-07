package repositories

import (
	"fmt"
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
	if err := r.db.First(&m, id).Error; err != nil {
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

	if err := r.db.Preload("Children").Offset(offset).Limit(pageSize).Find(&rows).Error; err != nil {
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
	if err := r.db.Preload("Children", "active = ?", true).Where("active = ?", true).Order("\"order\" asc").Find(&rows).Error; err != nil {
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
	query := r.db.Order("\"order\" asc")

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

func (r *menuRepository) FindMenusByRoleID(roleID uuid.UUID) ([]*entities.Menu, error) {
	var rows []*models.MenuModel
	var permissionIDs []uuid.UUID
	permErr := r.db.Table("role_permissions").
		Select("permission_id").
		Where("role_id = ?", roleID).
		Pluck("permission_id", &permissionIDs).Error
	if permErr != nil {
		return nil, permErr
	}

	fmt.Println("perms ids", permissionIDs)

	err := r.db.Table("menus").
		Joins("INNER JOIN model_permissions ON menus.id = CAST(model_permissions.model_id as uuid)").
		Where("model_permissions.model_type = ? AND model_permissions.permission_id IN ? AND menus.is_visible = ? AND menus.is_active = ?", "menu", permissionIDs, true, true).
		Order("menus.sequence ASC").Find(&rows).Error

	menus := make([]*entities.Menu, 0, len(rows))
	for _, m := range rows {
		menus = append(menus, toEntityMenu(m))
	}
	return menus, err
}

func (r *menuRepository) MenuBySuperUser() ([]*entities.Menu, error) {
	var rows []*models.MenuModel
	if err := r.db.Table("menus").
		Joins("INNER JOIN model_permissions ON menus.id = CAST(model_permissions.model_id as uuid)").
		Where("menus.is_active = ? AND menus.is_visible = ?", true, true).Order("menus.sequence asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	menus := make([]*entities.Menu, 0, len(rows))
	for _, m := range rows {
		menus = append(menus, toEntityMenu(m))
	}
	return menus, nil
}
