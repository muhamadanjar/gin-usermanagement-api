package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string         `gorm:"unique;not null" json:"name"`
	Description string         `json:"description"`
	Users       []*User        `gorm:"many2many:user_roles;" json:"users,omitempty"`
	Permissions []*Permission  `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
