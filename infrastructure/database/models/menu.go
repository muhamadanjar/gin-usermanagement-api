package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MenuPermissionLink struct {
	MenuID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	PermissionID uuid.UUID `gorm:"type:uuid;primaryKey"`
}

func (MenuPermissionLink) TableName() string {
	return "menu_permission"
}

type MenuModel struct {
	ID            uuid.UUID          `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name          string             `gorm:"unique;not null" json:"name"`
	Url           string             `json:"url"`
	PermissionKey string             `gorm:"unique;index" json:"permission_key"`
	Icon          string             `json:"icon"`
	Description   string             `json:"description"`
	ParentID      *uuid.UUID         `json:"parent_id"`
	Parent        *MenuModel         `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children      []*MenuModel       `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Permissions   []*PermissionModel `gorm:"many2many:menu_permission;" json:"permissions,omitempty"`
	Sequence      int                `gorm:"default:0" json:"sequence"`
	IsActive      bool               `gorm:"default:true" json:"is_active"`
	IsVisible     bool               `gorm:"default:true" json:"is_visible"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
	DeletedAt     gorm.DeletedAt     `gorm:"index" json:"-"`
}
