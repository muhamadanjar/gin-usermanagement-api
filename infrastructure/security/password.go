package security

import (
	"usermanagement-api/domain/ports"
	"usermanagement-api/pkg/utils"
)

// passwordHasher satisfies ports.PasswordHasher on top of pkg/utils (bcrypt).
type passwordHasher struct{}

func NewPasswordHasher() ports.PasswordHasher {
	return &passwordHasher{}
}

func (h *passwordHasher) Hash(password string) (string, error) {
	return utils.HashPassword(password)
}

func (h *passwordHasher) Check(password, hash string) bool {
	return utils.CheckPasswordHash(password, hash)
}
