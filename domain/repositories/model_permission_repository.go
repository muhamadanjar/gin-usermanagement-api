package repositories

import (
	"usermanagement-api/domain/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

type modelPermissionRepository struct {
	db *gorm.DB
}

func NewModelPermissionRepository(db *gorm.DB) ModelPermissionRepository {
	return &modelPermissionRepository{db}
}

func (r *modelPermissionRepository) Create(modelPermission *entities.ModelPermission) error {
	return r.db.Create(modelPermission).Error
}

func (r *modelPermissionRepository) FindByID(id uuid.UUID) (*entities.ModelPermission, error) {
	var modelPermission entities.ModelPermission
	if err := r.db.Preload("Permission").First(&modelPermission, id).Error; err != nil {
		return nil, err
	}
	return &modelPermission, nil
}

func (r *modelPermissionRepository) FindByModelTypeAndModelID(modelType string, modelID uuid.UUID) ([]*entities.ModelPermission, error) {
	var modelPermissions []*entities.ModelPermission
	if err := r.db.Preload("Permission").Where("model_type = ? AND model_id = ?", modelType, modelID).Find(&modelPermissions).Error; err != nil {
		return nil, err
	}
	return modelPermissions, nil
}

func (r *modelPermissionRepository) FindAll(page, pageSize int) ([]*entities.ModelPermission, int64, error) {
	var modelPermissions []*entities.ModelPermission
	var count int64

	offset := (page - 1) * pageSize

	if err := r.db.Model(&entities.ModelPermission{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Preload("Permission").Offset(offset).Limit(pageSize).Find(&modelPermissions).Error; err != nil {
		return nil, 0, err
	}

	return modelPermissions, count, nil
}

func (r *modelPermissionRepository) Update(modelPermission *entities.ModelPermission) error {
	return r.db.Save(modelPermission).Error
}

func (r *modelPermissionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entities.ModelPermission{}, id).Error
}

func (r *modelPermissionRepository) CheckPermission(modelType string, modelID uuid.UUID, permissionID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&entities.ModelPermission{}).
		Where("model_type = ? AND model_id = ? AND permission_id = ?", modelType, modelID, permissionID).
		Count(&count).Error

	return count > 0, err
}
