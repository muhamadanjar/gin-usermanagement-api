package ports

import "github.com/google/uuid"

// TokenPair is the domain-side token pair, free of any JWT library type.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int // seconds until the access token expires
}

// TokenManager is the port for issuing and validating bearer tokens.
// Implementations live in infrastructure (JWT adapter over pkg/auth).
type TokenManager interface {
	GenerateTokenPair(userID uuid.UUID, email string) (*TokenPair, error)
	ValidateAccessToken(tokenString string) (uuid.UUID, error)
	ValidateRefreshToken(tokenString string) (uuid.UUID, error)
}
