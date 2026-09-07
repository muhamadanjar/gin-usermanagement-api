package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TokenHistoryModel struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Token      string         `gorm:"uniqueIndex" json:"token"`
	UserID     uuid.UUID      `gorm:"type:uuid;index" json:"user_id"`
	ExpiredAt  time.Time      `json:"expired_at"`
	LastUsedAt time.Time      `json:"last_used_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
