package usecase

import (
	"errors"
	"fmt"
	"time"
	"usermanagement-api/domain/entities"
	"usermanagement-api/domain/repositories"
	"usermanagement-api/internal/dto"
	"usermanagement-api/pkg/auth"
	"usermanagement-api/pkg/firebase"
	"usermanagement-api/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthUseCase interface {
	Login(req *dto.LoginRequest) (*dto.AuthInfoResponse, error)
	Register(req *dto.RegisterRequest) (*dto.UserResponse, error)
	GetUserPermissions(userID uuid.UUID) ([]*entities.Permission, error)
	CreateModelPermission(req *dto.ModelPermissionRequest) (*dto.ModelPermissionResponse, error)
	GetModelPermissions(modelType string, modelID uuid.UUID) ([]*dto.ModelPermissionResponse, error)
	CheckPermission(modelType string, modelID uuid.UUID, permissionID uuid.UUID) (bool, error)
	GetUser(userID uuid.UUID, token string) (*dto.AuthInfoResponse, error)
	CreateMetaData(userID uuid.UUID, req *dto.CreateMetaDataRequest) (any, error)
	GetMetaData(userID uuid.UUID) ([]*dto.UserMetaResponse, error)
	// SendToUser(userID uuid.UUID, title string, body string, data map[string]string) (any, error)
}

type authUseCase struct {
	userRepo            repositories.UserRepository
	roleRepo            repositories.RoleRepository
	menuRepo            repositories.MenuRepository
	userMetaRepo        repositories.UserMetaRepository
	modelPermissionRepo repositories.ModelPermissionRepository
	fcmClient           firebase.FCMClient
}

func NewAuthUseCase(
	userRepo repositories.UserRepository,
	roleRepo repositories.RoleRepository,
	menuRepo repositories.MenuRepository,
	modelPermissionRepo repositories.ModelPermissionRepository,
	userMetaRepo repositories.UserMetaRepository,
	fcmClient firebase.FCMClient,

) AuthUseCase {
	return &authUseCase{
		userRepo:            userRepo,
		roleRepo:            roleRepo,
		menuRepo:            menuRepo,
		userMetaRepo:        userMetaRepo,
		modelPermissionRepo: modelPermissionRepo,
		fcmClient:           fcmClient,
	}
}

func (uc *authUseCase) Login(req *dto.LoginRequest) (*dto.AuthInfoResponse, error) {
	// Find user by email
	user, err := uc.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Verify password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := auth.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	authResp := &dto.AuthResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Type:         "Bearer",
	}

	// Map user to response
	userResp := &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsActive:  user.IsActive,
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

	return &dto.AuthInfoResponse{
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
		IsActive:  true,
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
		IsActive:  user.IsActive,
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

func (uc *authUseCase) GetUser(userID uuid.UUID, token string) (*dto.AuthInfoResponse, error) {

	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	userResp := &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Name:      user.FirstName + " " + user.LastName,
		IsActive:  user.IsActive,
		AvatarUrl: "https://gravatar.com/avatar/" + user.Email,
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

	privileges, err := uc.getPrivilegesForUser(user.ID)
	fmt.Println("Privileges:", privileges)
	if err == nil && len(privileges) > 0 {
		userResp.Privileges = privileges
	}

	return &dto.AuthInfoResponse{
		User: *userResp,
		Auth: dto.AuthResponse{
			AccessToken: token,
			Type:        "Bearer",
		},
	}, nil
}

func (uc *authUseCase) getPrivilegesForUser(userID uuid.UUID) ([]dto.MenuResponse, error) {
	// This would need to be implemented based on your menu repository and how
	// privileges are associated with users (through roles, etc.)

	// For example:
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user.IsSuperuser {
		menu, err := uc.menuRepo.MenuBySuperUser()
		if err != nil {
			return nil, err
		}
		var privileges []dto.MenuResponse
		for _, menu := range menu {
			privileges = append(privileges, dto.MenuResponse{
				ID:        menu.ID,
				Name:      menu.Name,
				Url:       menu.Url,
				Icon:      menu.Icon,
				ParentID:  menu.ParentID,
				IsActive:  menu.IsActive,
				IsVisible: menu.IsVisible, // You might want to add this field to your menu entity
				Sequence:  menu.Sequence,  // Assuming Order corresponds to Sequence
				CreatedAt: menu.CreatedAt.Format(time.RFC3339),
			})
		}
		return privileges, nil

	}
	roles, err := uc.roleRepo.FindRolesByUserID(userID)
	if err != nil {
		return nil, err
	}
	fmt.Println("User Roles:", userID)

	// Get all menu IDs accessible by these roles
	var menuIDs []uuid.UUID
	for _, role := range roles {
		// This assumes you have a way to get menu IDs by role
		// You might need to add this method to your repository

		roleMenus, err := uc.menuRepo.FindMenusByRoleID(role.ID)
		if err != nil {
			continue
		}

		for _, menu := range roleMenus {
			menuIDs = append(menuIDs, menu.ID)
		}
	}

	// Remove duplicates
	uniqueMenuIDs := make(map[uuid.UUID]bool)
	for _, id := range menuIDs {
		uniqueMenuIDs[id] = true
	}

	// Get menu details
	var privileges []dto.MenuResponse
	for id := range uniqueMenuIDs {
		menu, err := uc.menuRepo.FindByID(id)
		if err != nil {
			continue
		}

		// Map to MenuResponse
		privileges = append(privileges, dto.MenuResponse{
			ID:        menu.ID,
			Name:      menu.Name,
			Url:       menu.Url,
			Icon:      menu.Icon,
			ParentID:  menu.ParentID,
			IsActive:  menu.IsActive,
			IsVisible: menu.IsVisible, // You might want to add this field to your menu entity
			Sequence:  menu.Sequence,  // Assuming Order corresponds to Sequence
			CreatedAt: menu.CreatedAt.Format(time.RFC3339),
			// UpdatedAt: formatTimePointer(menu.UpdatedAt),
			// DeletedAt: formatTimePointer(menu.DeletedAt.Time),
		})
	}

	return privileges, nil
}

func (uc *authUseCase) CreateMetaData(userID uuid.UUID, req *dto.CreateMetaDataRequest) (any, error) {

	existingMeta, err := uc.userMetaRepo.FindByUserIDAndKey(userID, req.Key)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if existingMeta != nil {
		// Update existing
		existingMeta.Value = req.Value
		if err := uc.userMetaRepo.Update(existingMeta); err != nil {
			return nil, err
		}
	} else {
		// Create new
		userMeta := &entities.UserMeta{
			Key:    req.Key,
			Value:  req.Value,
			UserID: userID,
		}
		if err := uc.userMetaRepo.Create(userMeta); err != nil {
			return nil, err
		}
	}

	return existingMeta, nil
}

func (uc *authUseCase) GetMetaData(userID uuid.UUID) ([]*dto.UserMetaResponse, error) {
	metas, err := uc.userMetaRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	var response []*dto.UserMetaResponse
	for _, meta := range metas {
		response = append(response, &dto.UserMetaResponse{
			ID:     meta.ID,
			Key:    meta.Key,
			Value:  meta.Value,
			UserID: meta.UserID,
		})
	}
	return response, nil
}

// func (uc *authUseCase) SendToUser(userID uuid.UUID, title string, body string, data map[string]string) (any, error) {
// 	// Get user device tokens
// 	userMeta, err := uc.userMetaRepo.FindByUserIDAndKey(userID, "device_token")
// 	if err != nil {
// 		return nil, err
// 	}

// 	if userMeta == nil {
// 		return nil, errors.New("no device token found for user")
// 	}

// 	// Send notification
// 	response, err := uc.fcmClient.s(userMeta.Value, title, body, data)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return response, nil

// }

func formatTimePointer(t time.Time) *string {
	if t.IsZero() {
		return nil
	}
	formatted := t.Format(time.RFC3339)
	return &formatted
}
