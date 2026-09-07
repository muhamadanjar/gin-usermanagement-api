package usecase

import (
	"errors"
	"time"
	"usermanagement-api/domain/entities"
	"usermanagement-api/domain/ports"
	"usermanagement-api/domain/repositories"
	"usermanagement-api/internal/application/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthUseCase interface {
	Login(req *dto.LoginRequest) (*dto.AuthInfoResponse, error)
	Register(req *dto.RegisterRequest) (*dto.UserResponse, error)
	GetUserPermissions(userID uuid.UUID) ([]*entities.Permission, error)
	GetUser(userID uuid.UUID, token string) (*dto.AuthInfoResponse, error)
	CreateMetaData(userID uuid.UUID, req *dto.CreateMetaDataRequest) (any, error)
	GetMetaData(userID uuid.UUID) ([]*dto.UserMetaResponse, error)
	Refresh(refreshToken string) (*dto.AuthInfoResponse, error)
	Logout(token string) error
	ChangePassword(userID uuid.UUID, req *dto.ChangePasswordRequest) error
	UpdateProfile(userID uuid.UUID, req *dto.ProfileUpdateRequest) (*dto.UserResponse, error)
	GetTokenHistory(userID uuid.UUID) ([]*dto.TokenHistoryResponse, error)
	// SendToUser(userID uuid.UUID, title string, body string, data map[string]string) (any, error)
}

type authUseCase struct {
	userRepo     repositories.UserRepository
	roleRepo     repositories.RoleRepository
	menuRepo     repositories.MenuRepository
	userMetaRepo repositories.UserMetaRepository
	notifier     ports.Notifier
	tokenManager ports.TokenManager
	hasher       ports.PasswordHasher
	log          ports.Logger
}

func NewAuthUseCase(
	userRepo repositories.UserRepository,
	roleRepo repositories.RoleRepository,
	menuRepo repositories.MenuRepository,
	userMetaRepo repositories.UserMetaRepository,
	notifier ports.Notifier,
	tokenManager ports.TokenManager,
	hasher ports.PasswordHasher,
	log ports.Logger,
) AuthUseCase {
	return &authUseCase{
		userRepo:     userRepo,
		roleRepo:     roleRepo,
		menuRepo:     menuRepo,
		userMetaRepo: userMetaRepo,
		notifier:     notifier,
		tokenManager: tokenManager,
		hasher:       hasher,
		log:          log,
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
	if !uc.hasher.Check(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := uc.tokenManager.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	// Record token history (mirrors FastAPI: created on every login).
	_ = uc.userRepo.AddTokenHistory(&entities.TokenHistory{
		UserID:     user.ID,
		Token:      token.AccessToken,
		ExpiredAt:  time.Now().Add(time.Duration(token.ExpiresIn) * time.Second),
		LastUsedAt: time.Now(),
	})

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
	hashedPassword, err := uc.hasher.Hash(req.Password)
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
	uc.log.Info("Roles found")

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

	uc.log.Info("Permissions found", "permissions", permissions)

	return permissions, nil
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

func formatTimePointer(t time.Time) *string {
	if t.IsZero() {
		return nil
	}
	formatted := t.Format(time.RFC3339)
	return &formatted
}

func (uc *authUseCase) Refresh(refreshToken string) (*dto.AuthInfoResponse, error) {
	userID, err := uc.tokenManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}
	pair, err := uc.tokenManager.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		return nil, err
	}
	_ = uc.userRepo.AddTokenHistory(&entities.TokenHistory{
		UserID:     user.ID,
		Token:      pair.AccessToken,
		ExpiredAt:  time.Now().Add(time.Duration(pair.ExpiresIn) * time.Second),
		LastUsedAt: time.Now(),
	})
	// Rotate: revoke the used refresh token.
	_ = uc.userRepo.RemoveTokenHistoryByToken(refreshToken)

	userResp := &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Name:      user.Name,
		IsActive:  user.IsActive,
		AvatarUrl: "https://gravatar.com/avatar/" + user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
	return &dto.AuthInfoResponse{
		Auth: dto.AuthResponse{
			AccessToken:  pair.AccessToken,
			RefreshToken: pair.RefreshToken,
			Type:         "bearer",
		},
		User: *userResp,
	}, nil
}

func (uc *authUseCase) Logout(token string) error {
	return uc.userRepo.RemoveTokenHistoryByToken(token)
}

func (uc *authUseCase) ChangePassword(userID uuid.UUID, req *dto.ChangePasswordRequest) error {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("invalid credentials")
	}
	if !uc.hasher.Check(req.OldPassword, user.Password) {
		return errors.New("old password is incorrect")
	}
	hashed, err := uc.hasher.Hash(req.NewPassword)
	if err != nil {
		return err
	}
	user.Password = hashed
	return uc.userRepo.Update(user)
}

func (uc *authUseCase) UpdateProfile(userID uuid.UUID, req *dto.ProfileUpdateRequest) (*dto.UserResponse, error) {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Email != "" && req.Email != user.Email {
		if existing, err := uc.userRepo.FindByEmail(req.Email); err == nil && existing.ID != userID {
			return nil, errors.New("email already exists")
		}
		user.Email = req.Email
		user.Username = req.Email
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if err := uc.userRepo.Update(user); err != nil {
		return nil, err
	}
	updated, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	resp := &dto.UserResponse{
		ID:          updated.ID,
		Username:    updated.Username,
		Email:       updated.Email,
		Name:        updated.Name,
		FirstName:   updated.FirstName,
		LastName:    updated.LastName,
		IsActive:    updated.IsActive,
		IsSuperuser: updated.IsSuperuser,
		Status:      updated.Status,
		AvatarUrl:   "https://gravatar.com/avatar/" + updated.Email,
		CreatedAt:   updated.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   updated.UpdatedAt.Format(time.RFC3339),
	}
	for _, r := range updated.Roles {
		resp.Roles = append(resp.Roles, dto.RoleSimple{ID: r.ID, Name: r.Name})
	}
	return resp, nil
}

func (uc *authUseCase) GetTokenHistory(userID uuid.UUID) ([]*dto.TokenHistoryResponse, error) {
	histories, err := uc.userRepo.FindTokenHistory(userID)
	if err != nil {
		return nil, err
	}
	resp := make([]*dto.TokenHistoryResponse, 0, len(histories))
	for _, h := range histories {
		resp = append(resp, &dto.TokenHistoryResponse{
			ID:         h.ID,
			Token:      h.Token,
			ExpiredAt:  h.ExpiredAt.Format(time.RFC3339),
			LastUsedAt: h.LastUsedAt.Format(time.RFC3339),
		})
	}
	return resp, nil
}
