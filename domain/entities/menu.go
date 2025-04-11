package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Menu struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name        string         `gorm:"unique;not null" json:"name"`
	Path        string         `json:"path"`
	Icon        string         `json:"icon"`
	Description string         `json:"description"`
	ParentID    *uuid.UUID     `json:"parent_id"`
	Parent      *Menu          `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []*Menu        `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Order       int            `gorm:"default:0" json:"order"`
	Active      bool           `gorm:"default:true" json:"active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
