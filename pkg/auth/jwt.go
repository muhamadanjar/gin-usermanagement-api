package auth

import (
	"errors"
	"fmt"
	"time"
	"usermanagement-api/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // Seconds until access token expires
}

// JWTService handles JWT token operations
type JWTService struct {
	jwtConfig config.JWTConfig
}

// NewJWTService creates a new JWT service
func NewJWTService(jwtConfig config.JWTConfig) *JWTService {
	return &JWTService{
		jwtConfig: jwtConfig,
	}
}

// GenerateTokenPair generates a new JWT token pair
func (s *JWTService) GenerateTokenPair(userID uuid.UUID, email string) (*TokenPair, error) {
	// Use access token expiration from config
	accessExpirationHours := s.jwtConfig.AccessTokenExpiration
	if accessExpirationHours == 0 {
		accessExpirationHours = 1 // Default to 1 hour
	}

	// Use refresh token expiration from config
	refreshExpirationHours := s.jwtConfig.RefreshTokenExpiration
	if refreshExpirationHours == 0 {
		refreshExpirationHours = 168 // Default to 7 days
	}

	// Generate access token
	accessExpiration := time.Now().Add(time.Duration(accessExpirationHours) * time.Hour)
	accessClaims := JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.jwtConfig.Secret))
	if err != nil {
		return nil, err
	}

	// Generate refresh token (with longer expiration)
	refreshExpiration := time.Now().Add(time.Duration(refreshExpirationHours) * time.Hour)
	refreshClaims := JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.jwtConfig.RefreshTokenSecret))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    accessExpirationHours * 3600, // Convert to seconds
	}, nil
}

// ValidateAccessToken validates an access token
func (s *JWTService) ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	return s.validateToken(tokenString, s.jwtConfig.Secret)
}

// ValidateRefreshToken validates a refresh token
func (s *JWTService) ValidateRefreshToken(tokenString string) (*JWTClaims, error) {
	return s.validateToken(tokenString, s.jwtConfig.RefreshTokenSecret)
}

func (s *JWTService) validateToken(tokenString, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// Legacy functions for backward compatibility (will be removed in future)
// These use a global JWTService instance that should be set via SetGlobalJWTService
var globalJWTService *JWTService

// SetGlobalJWTService sets the global JWT service instance
func SetGlobalJWTService(service *JWTService) {
	globalJWTService = service
}

// GenerateTokenPair is a legacy function that uses global JWT service
func GenerateTokenPair(userID uuid.UUID, email string) (*TokenPair, error) {
	if globalJWTService == nil {
		return nil, errors.New("JWT service not initialized")
	}
	return globalJWTService.GenerateTokenPair(userID, email)
}

// ValidateAccessToken is a legacy function that uses global JWT service
func ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	if globalJWTService == nil {
		return nil, errors.New("JWT service not initialized")
	}
	return globalJWTService.ValidateAccessToken(tokenString)
}

// ValidateRefreshToken is a legacy function that uses global JWT service
func ValidateRefreshToken(tokenString string) (*JWTClaims, error) {
	if globalJWTService == nil {
		return nil, errors.New("JWT service not initialized")
	}
	return globalJWTService.ValidateRefreshToken(tokenString)
}
