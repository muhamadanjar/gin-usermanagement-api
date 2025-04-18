package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Menu struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name        string         `gorm:"unique;not null" json:"name"`
	Url         string         `json:"url"`
	Icon        string         `json:"icon"`
	Description string         `json:"description"`
	ParentID    *uuid.UUID     `json:"parent_id"`
	Parent      *Menu          `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []*Menu        `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Sequence    int            `gorm:"default:0" json:"sequence"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	IsVisible   bool           `gorm:"default:true" json:"is_visible"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
