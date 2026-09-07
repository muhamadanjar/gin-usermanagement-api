package entities

import (
	"time"

	"github.com/google/uuid"
)

// User is a pure domain entity. Persistence concerns (GORM tags, soft delete)
// live in infrastructure/database/models and are mapped in the repository layer.
type User struct {
	ID          uuid.UUID
	Username    string
	Email       string
	Password    string
	FirstName   string
	LastName    string
	IsSuperuser bool
	IsActive    bool
	Roles       []*Role
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
