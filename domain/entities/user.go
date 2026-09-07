package entities

import (
	"time"

	"github.com/google/uuid"
)

// User mirrors the FastAPI UserModel columns (plus FirstName/LastName kept for
// API backward compatibility).
type User struct {
	ID                  uuid.UUID
	Username            string
	Email               string
	Name                string
	Password            string
	FirstName           string
	LastName            string
	IsVerified          bool
	IsSuperuser         bool
	IsActive            bool
	Avatar              string
	EmailVerifiedAt     time.Time
	FailedLoginAttempts int
	LockedUntil         time.Time
	LastLogin           time.Time
	Status              string
	Roles               []*Role
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
