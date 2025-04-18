package dto

import "github.com/google/uuid"

type CreateUserRequest struct {
	Username  string      `json:"username" binding:"required"`
	Email     string      `json:"email" binding:"required,email"`
	Password  string      `json:"password" binding:"required,min=6"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	RoleIDs   []uuid.UUID `json:"role_ids"`
}

type UpdateUserRequest struct {
	Username  string      `json:"username"`
	Email     string      `json:"email" binding:"omitempty,email"`
	Password  string      `json:"password" binding:"omitempty,min=6"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Active    *bool       `json:"active"`
	RoleIDs   []uuid.UUID `json:"role_ids"`
}

type UserResponse struct {
	ID          uuid.UUID      `json:"id"`
	Username    string         `json:"username"`
	Email       string         `json:"email"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	IsActive    bool           `json:"is_active"`
	IsSuperuser bool           `json:"is_superuser"`
	Roles       []RoleSimple   `json:"roles,omitempty"`
	Privileges  []MenuResponse `json:"privileges,omitempty"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
}

type UserSimple struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

type AssignRolesRequest struct {
	RoleIDs []uuid.UUID `json:"roles_ids" binding:"required"`
}
