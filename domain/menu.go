package domain

import "github.com/google/uuid"

type Menu struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name     string    `json:"name"`
	Url      string    `json:"url"`
	IsActive bool      `json:"is_active"`
}
