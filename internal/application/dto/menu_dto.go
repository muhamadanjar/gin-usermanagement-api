package dto

import "github.com/google/uuid"

// MenuAssignPermissionRequest mirrors FastAPI's MenuAssignPermissionRequest.
type MenuAssignPermissionRequest struct {
	PermissionInputs []string `json:"permission_inputs" binding:"required"`
}

type CreateMenuRequest struct {
	Name          string     `json:"name" binding:"required"`
	Url           string     `json:"url"`
	PermissionKey string     `json:"permission_key"`
	Icon          string     `json:"icon"`
	Description   string     `json:"description"`
	ParentID      *uuid.UUID `json:"parent_id"`
	Sequence      int        `json:"sequence"`
}

type UpdateMenuRequest struct {
	Name          string     `json:"name"`
	Url           string     `json:"url"`
	PermissionKey string     `json:"permission_key"`
	Icon          string     `json:"icon"`
	Description   string     `json:"description"`
	ParentID      *uuid.UUID `json:"parent_id"`
	Sequence      int        `json:"sequence"`
	IsActive      *bool      `json:"is_active"`
}

type MenuResponse struct {
	ID            uuid.UUID          `json:"id"`
	Name          string             `json:"name"`
	Url           string             `json:"url"`
	PermissionKey string             `json:"permission_key,omitempty"`
	Icon          string             `json:"icon"`
	Description   string             `json:"description"`
	ParentID      *uuid.UUID         `json:"parent_id"`
	Parent        *MenuSimple        `json:"parent,omitempty"`
	Children      []*MenuSimple      `json:"children,omitempty"`
	Permissions   []PermissionSimple `json:"permissions,omitempty"`
	Sequence      int                `json:"sequence"`
	IsActive      bool               `json:"active"`
	IsVisible     bool               `json:"is_visible"`
	CreatedAt     string             `json:"created_at"`
	UpdatedAt     string             `json:"updated_at"`
	DeletedAt     string             `json:"delete_at"`
}

type MenuSimple struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Url  string    `json:"url"`
}
