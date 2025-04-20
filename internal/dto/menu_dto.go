package dto

import "github.com/google/uuid"

type CreateMenuRequest struct {
	Name        string     `json:"name" binding:"required"`
	Url         string     `json:"url"`
	Icon        string     `json:"icon"`
	Description string     `json:"description"`
	ParentID    *uuid.UUID `json:"parent_id"`
	Sequence    int        `json:"sequence"`
}

type UpdateMenuRequest struct {
	Name        string     `json:"name"`
	Url         string     `json:"url"`
	Icon        string     `json:"icon"`
	Description string     `json:"description"`
	ParentID    *uuid.UUID `json:"parent_id"`
	Sequence    int        `json:"sequence"`
	IsActive    *bool      `json:"is_active"`
}

type MenuResponse struct {
	ID          uuid.UUID     `json:"id"`
	Name        string        `json:"name"`
	Url         string        `json:"url"`
	Icon        string        `json:"icon"`
	Description string        `json:"description"`
	ParentID    *uuid.UUID    `json:"parent_id"`
	Parent      *MenuSimple   `json:"parent,omitempty"`
	Children    []*MenuSimple `json:"children,omitempty"`
	Sequence    int           `json:"sequence"`
	IsActive    bool          `json:"active"`
	IsVisible   bool          `json:"is_visible"`
	CreatedAt   string        `json:"created_at"`
	UpdatedAt   string        `json:"updated_at"`
	DeletedAt   string        `json:"delete_at"`
}

type MenuSimple struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Url  string    `json:"url"`
}
