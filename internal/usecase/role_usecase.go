package usecase

import (
	"errors"
	"time"
	"usermanagement-api/domain/entities"
	"usermanagement-api/domain/repositories"
	"usermanagement-api/internal/dto"

	"github.com/google/uuid"
)

type RoleUseCase interface {
	Create(req *dto.CreateRoleRequest) (*dto.RoleResponse, error)
	GetByID(id uuid.UUID) (*dto.RoleResponse, error)
	GetAll(page, pageSize int) ([]*dto.RoleResponse, int64, error)
	Update(id uuid.UUID, req *dto.UpdateRoleRequest) (*dto.RoleResponse, error)
	Delete(id uuid.UUID) error
	AssignPermissions(roleID uuid.UUID, permissionIDs []uuid.UUID) (*dto.RoleResponse, error)
	GetUserRoles(userID uuid.UUID) ([]*dto.RoleResponse, error)
}

type roleUseCase struct {
	roleRepo       repositories.RoleRepository
	permissionRepo repositories.PermissionRepository
}

func NewRoleUseCase(roleRepo repositories.RoleRepository, permissionRepo repositories.PermissionRepository) RoleUseCase {
	return &roleUseCase{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
	}
}

func (uc *roleUseCase) Create(req *dto.CreateRoleRequest) (*dto.RoleResponse, error) {
	// Check if role name already exists
	if _, err := uc.roleRepo.FindByName(req.Name); err == nil {
		return nil, errors.New("role name already exists")
	}

	// Create role
	role := &entities.Role{
		Name:        req.Name,
		Description: req.Description,
	}

	// Save role
	if err := uc.roleRepo.Create(role); err != nil {
		return nil, err
	}

	return uc.mapToRoleResponse(role), nil
}

func (uc *roleUseCase) GetByID(id uuid.UUID) (*dto.RoleResponse, error) {
	role, err := uc.roleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return uc.mapToRoleResponse(role), nil
}

func (uc *roleUseCase) GetAll(page, pageSize int) ([]*dto.RoleResponse, int64, error) {
	roles, total, err := uc.roleRepo.FindAll(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	var response []*dto.RoleResponse
	for _, role := range roles {
		response = append(response, uc.mapToRoleResponse(role))
	}

	return response, total, nil
}

func (uc *roleUseCase) Update(id uuid.UUID, req *dto.UpdateRoleRequest) (*dto.RoleResponse, error) {
	role, err := uc.roleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" && req.Name != role.Name {
		// Check if new name already exists
		if existingRole, err := uc.roleRepo.FindByName(req.Name); err == nil && existingRole.ID != id {
			return nil, errors.New("role name already exists")
		}
		role.Name = req.Name
	}

	if req.Description != "" {
		role.Description = req.Description
	}

	// Update role
	if err := uc.roleRepo.Update(role); err != nil {
		return nil, err
	}

	return uc.mapToRoleResponse(role), nil
}

func (uc *roleUseCase) Delete(id uuid.UUID) error {
	return uc.roleRepo.Delete(id)
}

func (uc *roleUseCase) AssignPermissions(roleID uuid.UUID, permissionIDs []uuid.UUID) (*dto.RoleResponse, error) {
	// Check if role exists
	_, err := uc.roleRepo.FindByID(roleID)
	if err != nil {
		return nil, err
	}

	// Assign permissions
	if err := uc.roleRepo.AssignPermissions(roleID, permissionIDs); err != nil {
		return nil, err
	}

	// Get updated role with permissions
	updatedRole, err := uc.roleRepo.FindByID(roleID)
	if err != nil {
		return nil, err
	}

	return uc.mapToRoleResponse(updatedRole), nil
}

func (uc *roleUseCase) GetUserRoles(userID uuid.UUID) ([]*dto.RoleResponse, error) {
	roles, err := uc.roleRepo.FindRolesByUserID(userID)
	if err != nil {
		return nil, err
	}

	var response []*dto.RoleResponse
	for _, role := range roles {
		response = append(response, uc.mapToRoleResponse(role))
	}

	return response, nil
}

func (uc *roleUseCase) mapToRoleResponse(role *entities.Role) *dto.RoleResponse {
	resp := &dto.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   role.UpdatedAt.Format(time.RFC3339),
	}

	// Map users
	if role.Users != nil {
		for _, user := range role.Users {
			resp.Users = append(resp.Users, dto.UserSimple{
				ID:       user.ID,
				Username: user.Username,
				Email:    user.Email,
			})
		}
	}

	// Map permissions
	if role.Permissions != nil {
		for _, permission := range role.Permissions {
			resp.Permissions = append(resp.Permissions, dto.PermissionSimple{
				ID:   permission.ID,
				Name: permission.Name,
			})
		}
	}

	return resp
}
