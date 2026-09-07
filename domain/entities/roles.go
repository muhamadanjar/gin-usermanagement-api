package entities

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID          uuid.UUID
	Name        string
	Description string
	IsActive    bool
	Users       []*User
	Permissions []*Permission
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
