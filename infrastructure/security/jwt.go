package security

import (
	"usermanagement-api/domain/ports"
	"usermanagement-api/pkg/auth"

	"github.com/google/uuid"
)

// tokenManagerAdapter satisfies ports.TokenManager on top of the shared JWT
// library (pkg/auth). It is the only place that knows ports.TokenPair maps to
// auth.TokenPair / claims.
type tokenManagerAdapter struct {
	svc *auth.JWTService
}

func NewTokenManager(svc *auth.JWTService) ports.TokenManager {
	return &tokenManagerAdapter{svc: svc}
}

func (a *tokenManagerAdapter) GenerateTokenPair(userID uuid.UUID, email string) (*ports.TokenPair, error) {
	pair, err := a.svc.GenerateTokenPair(userID, email)
	if err != nil {
		return nil, err
	}
	return &ports.TokenPair{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		ExpiresIn:    pair.ExpiresIn,
	}, nil
}

func (a *tokenManagerAdapter) ValidateAccessToken(tokenString string) (uuid.UUID, error) {
	claims, err := a.svc.ValidateAccessToken(tokenString)
	if err != nil {
		return uuid.Nil, err
	}
	return claims.UserID, nil
}

func (a *tokenManagerAdapter) ValidateRefreshToken(tokenString string) (uuid.UUID, error) {
	claims, err := a.svc.ValidateRefreshToken(tokenString)
	if err != nil {
		return uuid.Nil, err
	}
	return claims.UserID, nil
}
