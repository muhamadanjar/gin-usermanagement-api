package domain

import "github.com/google/uuid"

type User struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name       string
	Email      string
	Username   string
	IsActive   bool
	IsVerified bool   `json:"is_verified"`
	Role       []Role `gorm:"many2many:user_roles;" json:"roles"`
}
