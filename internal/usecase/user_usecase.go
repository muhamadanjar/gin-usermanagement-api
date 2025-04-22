package usecase

import (
	"errors"
	"log"
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
	GetUserWithMeta(id uuid.UUID) (*dto.UserResponse, error)
}

type userUseCase struct {
	userRepo     repositories.UserRepository
	roleRepo     repositories.RoleRepository
	userMetaRepo repositories.UserMetaRepository
}

func NewUserUseCase(userRepo repositories.UserRepository, roleRepo repositories.RoleRepository, userMetaRepo repositories.UserMetaRepository) UserUseCase {
	return &userUseCase{
		userRepo:     userRepo,
		roleRepo:     roleRepo,
		userMetaRepo: userMetaRepo,
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

	if len(req.MetaData) > 0 {
		for key, value := range req.MetaData {
			userMeta := &entities.UserMeta{
				Key:    key,
				Value:  value,
				UserID: user.ID,
			}
			if err := uc.userMetaRepo.Create(userMeta); err != nil {
				// Log error but don't fail the whole operation
				log.Printf("Failed to create user meta %s: %v", key, err)
			}
		}
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

	if req.MetaData != nil {
		// Get existing meta
		existingMeta, err := uc.userMetaRepo.GetAllByUserID(user.ID)
		if err != nil {
			log.Printf("Failed to get existing user meta: %v", err)
		}

		// Update or create meta
		for key, value := range req.MetaData {
			if existingMeta[key] != "" {
				// Update existing meta
				existingUserMeta, err := uc.userMetaRepo.FindByUserIDAndKey(user.ID, key)
				if err == nil {
					existingUserMeta.Value = value
					if err := uc.userMetaRepo.Update(existingUserMeta); err != nil {
						log.Printf("Failed to update user meta %s: %v", key, err)
					}
				}
			} else {
				// Create new meta
				userMeta := &entities.UserMeta{
					Key:    key,
					Value:  value,
					UserID: user.ID,
				}
				if err := uc.userMetaRepo.Create(userMeta); err != nil {
					log.Printf("Failed to create user meta %s: %v", key, err)
				}
			}
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
		AvatarUrl: "https://gravatar.com/avatar/" + user.Email,
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
	} else {
		resp.Roles = []dto.RoleSimple{}
	}
	// Map privileges
	resp.Privileges = []dto.MenuResponse{}

	return resp
}

func (uc *userUseCase) GetUserWithMeta(id uuid.UUID) (*dto.UserResponse, error) {
	user, err := uc.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Get user meta
	userMeta, err := uc.userMetaRepo.GetAllByUserID(user.ID)
	if err != nil {
		log.Printf("Failed to get user meta: %v", err)
		userMeta = make(map[string]string) // Empty map if error
	}

	response := uc.mapToUserResponse(user)
	response.MetaData = userMeta

	return response, nil
}
