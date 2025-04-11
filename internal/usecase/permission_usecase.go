package usecase

import (
	"errors"
	"time"
	"usermanagement-api/domain/entities"
	"usermanagement-api/domain/repositories"
	"usermanagement-api/internal/dto"

	"github.com/google/uuid"
)

type PermissionUseCase interface {
	Create(req *dto.CreatePermissionRequest) (*dto.PermissionResponse, error)
	GetByID(id uuid.UUID) (*dto.PermissionResponse, error)
	GetAll(page, pageSize int) ([]*dto.PermissionResponse, int64, error)
	Update(id uuid.UUID, req *dto.UpdatePermissionRequest) (*dto.PermissionResponse, error)
	Delete(id uuid.UUID) error
}

type permissionUseCase struct {
	permissionRepo repositories.PermissionRepository
}

func NewPermissionUseCase(permissionRepo repositories.PermissionRepository) PermissionUseCase {
	return &permissionUseCase{
		permissionRepo: permissionRepo,
	}
}

func (uc *permissionUseCase) Create(req *dto.CreatePermissionRequest) (*dto.PermissionResponse, error) {
	// Check if permission name already exists
	if _, err := uc.permissionRepo.FindByName(req.Name); err == nil {
		return nil, errors.New("permission name already exists")
	}

	// Create permission
	permission := &entities.Permission{
		Name:        req.Name,
		Description: req.Description,
	}

	// Save permission
	if err := uc.permissionRepo.Create(permission); err != nil {
		return nil, err
	}

	return uc.mapToPermissionResponse(permission), nil
}

func (uc *permissionUseCase) GetByID(id uuid.UUID) (*dto.PermissionResponse, error) {
	permission, err := uc.permissionRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return uc.mapToPermissionResponse(permission), nil
}

func (uc *permissionUseCase) GetAll(page, pageSize int) ([]*dto.PermissionResponse, int64, error) {
	permissions, total, err := uc.permissionRepo.FindAll(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	var response []*dto.PermissionResponse
	for _, permission := range permissions {
		response = append(response, uc.mapToPermissionResponse(permission))
	}

	return response, total, nil
}

func (uc *permissionUseCase) Update(id uuid.UUID, req *dto.UpdatePermissionRequest) (*dto.PermissionResponse, error) {
	permission, err := uc.permissionRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" && req.Name != permission.Name {
		// Check if new name already exists
		if existingPermission, err := uc.permissionRepo.FindByName(req.Name); err == nil && existingPermission.ID != id {
			return nil, errors.New("permission name already exists")
		}
		permission.Name = req.Name
	}

	if req.Description != "" {
		permission.Description = req.Description
	}

	// Update permission
	if err := uc.permissionRepo.Update(permission); err != nil {
		return nil, err
	}

	return uc.mapToPermissionResponse(permission), nil
}

func (uc *permissionUseCase) Delete(id uuid.UUID) error {
	return uc.permissionRepo.Delete(id)
}

func (uc *permissionUseCase) mapToPermissionResponse(permission *entities.Permission) *dto.PermissionResponse {
	return &dto.PermissionResponse{
		ID:          permission.ID,
		Name:        permission.Name,
		Description: permission.Description,
		CreatedAt:   permission.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   permission.UpdatedAt.Format(time.RFC3339),
	}
}
