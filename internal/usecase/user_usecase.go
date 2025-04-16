package usecase

import (
	"errors"
	"time"
	"usermanagement-api/domain/entities"
	"usermanagement-api/domain/repositories"
	"usermanagement-api/internal/dto"
	"usermanagement-api/pkg/utils"

	"github.com/google/uuid"
)

type UserUseCase interface {
	Create(req *dto.CreateUserRequest) (*dto.UserResponse, error)
	GetByID(id uuid.UUID) (*dto.UserResponse, error)
	GetAll(page, pageSize int) ([]*dto.UserResponse, int64, error)
	Update(id uuid.UUID, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	Delete(id uuid.UUID) error
	AssignRoles(userID uuid.UUID, roleIDs []uuid.UUID) (*dto.UserResponse, error)
}

type userUseCase struct {
	userRepo repositories.UserRepository
	roleRepo repositories.RoleRepository
}

func NewUserUseCase(userRepo repositories.UserRepository, roleRepo repositories.RoleRepository) UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (uc *userUseCase) Create(req *dto.CreateUserRequest) (*dto.UserResponse, error) {
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

	// Add roles if provided
	if len(req.RoleIDs) > 0 {
		for _, roleID := range req.RoleIDs {
			user.Roles = append(user.Roles, &entities.Role{ID: roleID})
		}
	}

	// Save user
	if err := uc.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Get user with roles
	userWithRoles, err := uc.userRepo.FindByID(user.ID)
	if err != nil {
		return nil, err
	}

	// Map to response
	return uc.mapToUserResponse(userWithRoles), nil
}

func (uc *userUseCase) GetByID(id uuid.UUID) (*dto.UserResponse, error) {
	user, err := uc.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return uc.mapToUserResponse(user), nil
}

func (uc *userUseCase) GetAll(page, pageSize int) ([]*dto.UserResponse, int64, error) {
	users, total, err := uc.userRepo.FindAll(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	var response []*dto.UserResponse
	for _, user := range users {
		response = append(response, uc.mapToUserResponse(user))
	}

	return response, total, nil
}

func (uc *userUseCase) Update(id uuid.UUID, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := uc.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Username != "" && req.Username != user.Username {
		// Check if new username already exists
		if existingUser, err := uc.userRepo.FindByUsername(req.Username); err == nil && existingUser.ID != id {
			return nil, errors.New("username already exists")
		}
		user.Username = req.Username
	}

	if req.Email != "" && req.Email != user.Email {
		// Check if new email already exists
		if existingUser, err := uc.userRepo.FindByEmail(req.Email); err == nil && existingUser.ID != id {
			return nil, errors.New("email already exists")
		}
		user.Email = req.Email
	}

	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashedPassword
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}

	if req.LastName != "" {
		user.LastName = req.LastName
	}

	if req.Active != nil {
		user.IsActive = *req.Active
	}

	// Update user
	if err := uc.userRepo.Update(user); err != nil {
		return nil, err
	}

	// Update roles if provided
	if len(req.RoleIDs) > 0 {
		if err := uc.userRepo.AssignRoles(id, req.RoleIDs); err != nil {
			return nil, err
		}
		// Get updated user with roles
		user, err = uc.userRepo.FindByID(id)
		if err != nil {
			return nil, err
		}
	}

	return uc.mapToUserResponse(user), nil
}

func (uc *userUseCase) Delete(id uuid.UUID) error {
	return uc.userRepo.Delete(id)
}

func (uc *userUseCase) AssignRoles(userID uuid.UUID, roleIDs []uuid.UUID) (*dto.UserResponse, error) {
	// Check if user exists
	_, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	// Assign roles
	if err := uc.userRepo.AssignRoles(userID, roleIDs); err != nil {
		return nil, err
	}

	// Get updated user with roles
	updatedUser, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	return uc.mapToUserResponse(updatedUser), nil
}

func (uc *userUseCase) mapToUserResponse(user *entities.User) *dto.UserResponse {
	resp := &dto.UserResponse{
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
			resp.Roles = append(resp.Roles, dto.RoleSimple{
				ID:   role.ID,
				Name: role.Name,
			})
		}
	}

	return resp
}
