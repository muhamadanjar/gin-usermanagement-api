package repositories

import (
	"usermanagement-api/domain/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PermissionRepository interface {
	Create(permission *entities.Permission) error
	FindByID(id uuid.UUID) (*entities.Permission, error)
	FindByName(name string) (*entities.Permission, error)
	FindAll(page, pageSize int) ([]*entities.Permission, int64, error)
	Update(permission *entities.Permission) error
	Delete(id uuid.UUID) error
}

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db}
}

func (r *permissionRepository) Create(permission *entities.Permission) error {
	return r.db.Create(permission).Error
}

func (r *permissionRepository) FindByID(id uuid.UUID) (*entities.Permission, error) {
	var permission entities.Permission
	if err := r.db.First(&permission, id).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) FindByName(name string) (*entities.Permission, error) {
	var permission entities.Permission
	if err := r.db.Where("name = ?", name).First(&permission).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) FindAll(page, pageSize int) ([]*entities.Permission, int64, error) {
	var permissions []*entities.Permission
	var count int64

	offset := (page - 1) * pageSize

	if err := r.db.Model(&entities.Permission{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Offset(offset).Limit(pageSize).Find(&permissions).Error; err != nil {
		return nil, 0, err
	}

	return permissions, count, nil
}

func (r *permissionRepository) Update(permission *entities.Permission) error {
	return r.db.Save(permission).Error
}

func (r *permissionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entities.Permission{}, id).Error
}
