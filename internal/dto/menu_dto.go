package dto

import "github.com/google/uuid"

type CreateMenuRequest struct {
	Name        string     `json:"name" binding:"required"`
	Path        string     `json:"path"`
	Icon        string     `json:"icon"`
	Description string     `json:"description"`
	ParentID    *uuid.UUID `json:"parent_id"`
	Order       int        `json:"order"`
}

type UpdateMenuRequest struct {
	Name        string     `json:"name"`
	Path        string     `json:"path"`
	Icon        string     `json:"icon"`
	Description string     `json:"description"`
	ParentID    *uuid.UUID `json:"parent_id"`
	Order       int        `json:"order"`
	Active      *bool      `json:"active"`
}

type MenuResponse struct {
	ID          uuid.UUID     `json:"id"`
	Name        string        `json:"name"`
	Path        string        `json:"path"`
	Icon        string        `json:"icon"`
	Description string        `json:"description"`
	ParentID    *uuid.UUID    `json:"parent_id"`
	Parent      *MenuSimple   `json:"parent,omitempty"`
	Children    []*MenuSimple `json:"children,omitempty"`
	Order       int           `json:"order"`
	Active      bool          `json:"active"`
	CreatedAt   string        `json:"created_at"`
	UpdatedAt   string        `json:"updated_at"`
}

type MenuSimple struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Path string    `json:"path"`
}
