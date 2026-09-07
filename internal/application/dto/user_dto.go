package dto

import "github.com/google/uuid"

type CreateUserRequest struct {
	Username  string            `json:"username" binding:"required"`
	Email     string            `json:"email" binding:"required,email"`
	Password  string            `json:"password" binding:"required,min=6"`
	Name      string            `json:"name"`
	FirstName string            `json:"first_name"`
	LastName  string            `json:"last_name"`
	RoleIDs   []uuid.UUID       `json:"role_ids"`
	MetaData  map[string]string `json:"meta_data"`
}

type UpdateUserRequest struct {
	Username  string            `json:"username"`
	Email     string            `json:"email" binding:"omitempty,email"`
	Password  string            `json:"password" binding:"omitempty,min=6"`
	Name      string            `json:"name"`
	FirstName string            `json:"first_name"`
	LastName  string            `json:"last_name"`
	Active    *bool             `json:"active"`
	RoleIDs   []uuid.UUID       `json:"role_ids"`
	MetaData  map[string]string `json:"meta_data"`
}

type UserResponse struct {
	ID          uuid.UUID         `json:"id"`
	Username    string            `json:"username"`
	Email       string            `json:"email"`
	Name        string            `json:"name,omitempty"`
	FirstName   string            `json:"first_name"`
	LastName    string            `json:"last_name"`
	IsActive    bool              `json:"is_active"`
	IsSuperuser bool              `json:"is_superuser"`
	IsVerified  bool              `json:"is_verified"`
	Status      string            `json:"status,omitempty"`
	AvatarUrl   string            `json:"avatar_url"`
	Roles       []RoleSimple      `json:"roles"`
	Privileges  []MenuResponse    `json:"privileges,omitempty"`
	MetaData    map[string]string `json:"meta_data,omitempty"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
}

type UserSimple struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

type AssignRolesRequest struct {
	RoleIDs []uuid.UUID `json:"roles_ids" binding:"required"`
}

// AssignRoleRequest mirrors FastAPI's AssignRoleUserRequest (single role append).
type AssignRoleRequest struct {
	RoleID uuid.UUID `json:"role_id" binding:"required"`
}

// UserRoleSyncRequest mirrors FastAPI's UserRoleSyncRequest (replace all roles).
type UserRoleSyncRequest struct {
	RoleIDs []uuid.UUID `json:"role_ids" binding:"required"`
}

type UpdateAvatarRequest struct {
	Avatar string `json:"avatar" binding:"required"`
}

type TokenHistoryResponse struct {
	ID         uuid.UUID `json:"id"`
	Token      string    `json:"token"`
	ExpiredAt  string    `json:"expired_at"`
	LastUsedAt string    `json:"last_used_at"`
}
