package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserModel struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Username    string         `gorm:"unique;not null" json:"username"`
	Email       string         `gorm:"unique;not null" json:"email"`
	Password    string         `gorm:"not null" json:"-"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	IsSuperuser bool           `gorm:"not null;default:false" json:"is_superuser"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	Roles       []*RoleModel   `gorm:"many2many:user_roles;" json:"roles"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type RoleModel struct {
	ID          uuid.UUID          `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string             `gorm:"unique;not null" json:"name"`
	Description string             `json:"description"`
	Users       []*UserModel       `gorm:"many2many:user_roles;" json:"users,omitempty"`
	Permissions []*PermissionModel `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	DeletedAt   gorm.DeletedAt     `gorm:"index" json:"-"`
}

type PermissionModel struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name        string         `gorm:"unique;not null" json:"name"`
	Description string         `json:"description"`
	Roles       []*RoleModel   `gorm:"many2many:role_permissions;" json:"roles,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
