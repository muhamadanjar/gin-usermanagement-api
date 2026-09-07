package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserModel struct {
	ID                  uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Username            string         `gorm:"unique;not null" json:"username"`
	Email               string         `gorm:"unique;not null" json:"email"`
	Name                string         `json:"name"`
	Password            string         `gorm:"not null" json:"-"`
	FirstName           string         `json:"first_name"`
	LastName            string         `json:"last_name"`
	IsVerified          bool           `gorm:"not null;default:false" json:"is_verified"`
	IsSuperuser         bool           `gorm:"not null;default:false" json:"is_superuser"`
	IsActive            bool           `gorm:"default:true" json:"is_active"`
	Avatar              string         `json:"avatar"`
	EmailVerifiedAt     time.Time      `json:"email_verified_at"`
	FailedLoginAttempts int            `gorm:"default:0" json:"failed_login_attempts"`
	LockedUntil         time.Time      `json:"locked_until"`
	LastLogin           time.Time      `json:"last_login"`
	Status              string         `gorm:"default:pending" json:"status"`
	Roles               []*RoleModel   `gorm:"many2many:user_roles;" json:"roles"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`
}

type RoleModel struct {
	ID          uuid.UUID          `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string             `gorm:"unique;not null" json:"name"`
	Description string             `json:"description"`
	IsActive    bool               `gorm:"default:true" json:"is_active"`
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
	Menus       []*MenuModel   `gorm:"many2many:menu_permission;" json:"menus,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
