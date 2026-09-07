package entities

import (
	"time"

	"github.com/google/uuid"
)

type Permission struct {
	ID          uuid.UUID
	Name        string
	Description string
	Roles       []*Role
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
