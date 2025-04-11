package repositories

import (
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

	if err := r.db.Preload("Children").Where("parent_id IS NULL").Offset(offset).Limit(pageSize).Find(&menus).Error; err != nil {
		return nil, 0, err
	}

	return menus, count, nil
}

func (r *menuRepository) FindAllActive() ([]*entities.Menu, error) {
	var menus []*entities.Menu
	if err := r.db.Preload("Children", "active = ?", true).Where("active = ? AND parent_id IS NULL", true).Order("\"order\" asc").Find(&menus).Error; err != nil {
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
