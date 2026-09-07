package ports

// PasswordHasher is the port for password hashing/verification.
// Implementations: infrastructure/security (bcrypt adapter over pkg/utils).
type PasswordHasher interface {
	Hash(password string) (string, error)
	Check(password, hash string) bool
}
