package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ModelPermission struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey"`
	ModelID      uuid.UUID      `gorm:"not null" json:"model_id"`
	ModelType    string         `gorm:"not null" json:"model_type"` // Can be "role" or "menu" or other types
	PermissionID uuid.UUID      `gorm:"not null" json:"permission_id"`
	Permission   Permission     `gorm:"foreignKey:PermissionID" json:"permission"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
