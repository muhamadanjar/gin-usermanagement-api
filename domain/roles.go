package domain

import "github.com/google/uuid"

type Role struct {
	ID          uuid.UUID    `gorm:"type:uuid;primaryKey"`
	Name        string       `json:"name"`
	IsActive    bool         `json:"is_active"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions"`
}
