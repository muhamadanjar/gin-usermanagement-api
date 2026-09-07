package repositories

import (
	"usermanagement-api/domain/entities"
	domainrepos "usermanagement-api/domain/repositories"
	"usermanagement-api/infrastructure/database/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) domainrepos.PermissionRepository {
	return &permissionRepository{db}
}

func (r *permissionRepository) Create(permission *entities.Permission) error {
	m := toModelPermission(permission)
	if err := r.db.Create(m).Error; err != nil {
		return err
	}
	permission.ID = m.ID
	permission.CreatedAt = m.CreatedAt
	permission.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *permissionRepository) FindByID(id uuid.UUID) (*entities.Permission, error) {
	var m models.PermissionModel
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return toEntityPermission(&m), nil
}

func (r *permissionRepository) FindByName(name string) (*entities.Permission, error) {
	var m models.PermissionModel
	if err := r.db.Where("name = ?", name).First(&m).Error; err != nil {
		return nil, err
	}
	return toEntityPermission(&m), nil
}

func (r *permissionRepository) FindAll(page, pageSize int) ([]*entities.Permission, int64, error) {
	var rows []*models.PermissionModel
	var count int64

	offset := (page - 1) * pageSize

	if err := r.db.Model(&models.PermissionModel{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Offset(offset).Limit(pageSize).Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	permissions := make([]*entities.Permission, 0, len(rows))
	for _, m := range rows {
		permissions = append(permissions, toEntityPermission(m))
	}
	return permissions, count, nil
}

func (r *permissionRepository) Update(permission *entities.Permission) error {
	return r.db.Save(toModelPermission(permission)).Error
}

func (r *permissionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.PermissionModel{}, id).Error
}
