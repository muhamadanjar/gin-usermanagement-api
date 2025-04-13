package usecase

import (
	"errors"
	"time"
	"usermanagement-api/domain/entities"
	"usermanagement-api/domain/repositories"
	"usermanagement-api/internal/dto"
	"usermanagement-api/pkg/auth"
	"usermanagement-api/pkg/utils"

	"github.com/google/uuid"
)

type AuthUseCase interface {
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error)
	Register(req *dto.RegisterRequest) (*dto.UserResponse, error)
	GetUserPermissions(userID uuid.UUID) ([]*entities.Permission, error)
	CreateModelPermission(req *dto.ModelPermissionRequest) (*dto.ModelPermissionResponse, error)
	GetModelPermissions(modelType string, modelID uuid.UUID) ([]*dto.ModelPermissionResponse, error)
	CheckPermission(modelType string, modelID uuid.UUID, permissionID uuid.UUID) (bool, error)
}

type authUseCase struct {
	userRepo            repositories.UserRepository
	roleRepo            repositories.RoleRepository
	modelPermissionRepo repositories.ModelPermissionRepository
}

func NewAuthUseCase(
	userRepo repositories.UserRepository,
	roleRepo repositories.RoleRepository,
	modelPermissionRepo repositories.ModelPermissionRepository,
) AuthUseCase {
	return &authUseCase{
		userRepo:            userRepo,
		roleRepo:            roleRepo,
		modelPermissionRepo: modelPermissionRepo,
	}
}

func (uc *authUseCase) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// Find user by email
	user, err := uc.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check if user is active
	if !user.Active {
		return nil, errors.New("account is deactivated")
	}

	// Verify password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	authResp := &dto.AuthResponse{
		Token: token,
		Type:  "Bearer",
	}

	// Map user to response
	userResp := &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Active:    user.Active,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}

	// Map roles
	if len(user.Roles) > 0 {
		for _, role := range user.Roles {
			userResp.Roles = append(userResp.Roles, dto.RoleSimple{
				ID:   role.ID,
				Name: role.Name,
			})
		}
	}

	return &dto.LoginResponse{
		Auth: *authResp,
		User: *userResp,
	}, nil
}

func (uc *authUseCase) Register(req *dto.RegisterRequest) (*dto.UserResponse, error) {
	// Check if email already exists
	if _, err := uc.userRepo.FindByEmail(req.Email); err == nil {
		return nil, errors.New("email already exists")
	}

	// Check if username already exists
	if _, err := uc.userRepo.FindByUsername(req.Username); err == nil {
		return nil, errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &entities.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Active:    true,
	}

	// Save user
	if err := uc.userRepo.Create(user); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Active:    user.Active,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (uc *authUseCase) GetUserPermissions(userID uuid.UUID) ([]*entities.Permission, error) {
	// Get user roles
	roles, err := uc.roleRepo.FindRolesByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Get role IDs
	var roleIDs []uuid.UUID
	for _, role := range roles {
		roleIDs = append(roleIDs, role.ID)
	}

	// Get permissions for roles
	permissions, err := uc.roleRepo.FindPermissionsByRoleIDs(roleIDs)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

func (uc *authUseCase) CreateModelPermission(req *dto.ModelPermissionRequest) (*dto.ModelPermissionResponse, error) {
	// Create model permission
	modelPermission := &entities.ModelPermission{
		ModelID:      req.ModelID,
		ModelType:    req.ModelType,
		PermissionID: req.PermissionID,
	}

	// Save model permission
	if err := uc.modelPermissionRepo.Create(modelPermission); err != nil {
		return nil, err
	}

	// Get model permission with permission
	modelPermissionWithPermission, err := uc.modelPermissionRepo.FindByID(modelPermission.ID)
	if err != nil {
		return nil, err
	}

	return &dto.ModelPermissionResponse{
		ID:           modelPermissionWithPermission.ID,
		ModelID:      modelPermissionWithPermission.ModelID,
		ModelType:    modelPermissionWithPermission.ModelType,
		PermissionID: modelPermissionWithPermission.PermissionID,
		Permission: dto.PermissionSimple{
			ID:   modelPermissionWithPermission.Permission.ID,
			Name: modelPermissionWithPermission.Permission.Name,
		},
		CreatedAt: modelPermissionWithPermission.CreatedAt.Format(time.RFC3339),
		UpdatedAt: modelPermissionWithPermission.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (uc *authUseCase) GetModelPermissions(modelType string, modelID uuid.UUID) ([]*dto.ModelPermissionResponse, error) {
	// Get model permissions
	modelPermissions, err := uc.modelPermissionRepo.FindByModelTypeAndModelID(modelType, modelID)
	if err != nil {
		return nil, err
	}

	var response []*dto.ModelPermissionResponse
	for _, mp := range modelPermissions {
		response = append(response, &dto.ModelPermissionResponse{
			ID:           mp.ID,
			ModelID:      mp.ModelID,
			ModelType:    mp.ModelType,
			PermissionID: mp.PermissionID,
			Permission: dto.PermissionSimple{
				ID:   mp.Permission.ID,
				Name: mp.Permission.Name,
			},
			CreatedAt: mp.CreatedAt.Format(time.RFC3339),
			UpdatedAt: mp.UpdatedAt.Format(time.RFC3339),
		})
	}

	return response, nil
}

func (uc *authUseCase) CheckPermission(modelType string, modelID uuid.UUID, permissionID uuid.UUID) (bool, error) {
	return uc.modelPermissionRepo.CheckPermission(modelType, modelID, permissionID)
}
