package models

import "github.com/google/uuid"

type UserMetaModel struct {
	ID     uint      `gorm:"primaryKey" json:"id"`
	Key    string    `json:"key"`
	Value  string    `json:"value"`
	UserID uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
}
