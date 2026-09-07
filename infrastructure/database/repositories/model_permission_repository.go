package repositories

import (
	"usermanagement-api/domain/entities"
	domainrepos "usermanagement-api/domain/repositories"
	"usermanagement-api/infrastructure/database/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type modelPermissionRepository struct {
	db *gorm.DB
}

func NewModelPermissionRepository(db *gorm.DB) domainrepos.ModelPermissionRepository {
	return &modelPermissionRepository{db}
}

func (r *modelPermissionRepository) Create(modelPermission *entities.ModelPermission) error {
	m := toModelModelPermission(modelPermission)
	if err := r.db.Create(m).Error; err != nil {
		return err
	}
	modelPermission.ID = m.ID
	modelPermission.CreatedAt = m.CreatedAt
	modelPermission.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *modelPermissionRepository) FindByID(id uuid.UUID) (*entities.ModelPermission, error) {
	var m models.ModelPermissionModel
	if err := r.db.Preload("Permission").First(&m, id).Error; err != nil {
		return nil, err
	}
	return toEntityModelPermission(&m), nil
}

func (r *modelPermissionRepository) FindByModelTypeAndModelID(modelType string, modelID uuid.UUID) ([]*entities.ModelPermission, error) {
	var rows []*models.ModelPermissionModel
	if err := r.db.Preload("Permission").Where("model_type = ? AND model_id = ?", modelType, modelID).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*entities.ModelPermission, 0, len(rows))
	for _, m := range rows {
		result = append(result, toEntityModelPermission(m))
	}
	return result, nil
}

func (r *modelPermissionRepository) FindAll(page, pageSize int) ([]*entities.ModelPermission, int64, error) {
	var rows []*models.ModelPermissionModel
	var count int64

	offset := (page - 1) * pageSize

	if err := r.db.Model(&models.ModelPermissionModel{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Preload("Permission").Offset(offset).Limit(pageSize).Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*entities.ModelPermission, 0, len(rows))
	for _, m := range rows {
		result = append(result, toEntityModelPermission(m))
	}
	return result, count, nil
}

func (r *modelPermissionRepository) Update(modelPermission *entities.ModelPermission) error {
	return r.db.Save(toModelModelPermission(modelPermission)).Error
}

func (r *modelPermissionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.ModelPermissionModel{}, id).Error
}

func (r *modelPermissionRepository) CheckPermission(modelType string, modelID uuid.UUID, permissionID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.ModelPermissionModel{}).
		Where("model_type = ? AND model_id = ? AND permission_id = ?", modelType, modelID, permissionID).
		Count(&count).Error

	return count > 0, err
}
