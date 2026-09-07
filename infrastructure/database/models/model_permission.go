package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ModelPermissionModel struct {
	ID           uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ModelID      uuid.UUID       `gorm:"not null" json:"model_id"`
	ModelType    string          `gorm:"not null" json:"model_type"` // may be "role" or "menu" or other types
	PermissionID uuid.UUID       `gorm:"not null" json:"permission_id"`
	Permission   PermissionModel `gorm:"foreignKey:PermissionID" json:"permission"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	DeletedAt    gorm.DeletedAt  `gorm:"index" json:"-"`
}
