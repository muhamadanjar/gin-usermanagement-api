package domain

import "github.com/google/uuid"

type User struct {
	ID         uuid.UUID
	Name       string
	Email      string
	Username   string
	IsActive   bool
	IsVerified bool `json:"is_verified"`
}
