package entities

import (
	"time"

	"github.com/google/uuid"
)

type TokenHistory struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Token      string
	ExpiredAt  time.Time
	LastUsedAt time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
