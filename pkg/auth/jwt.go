package auth

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

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

// GenerateToken generates a new JWT token
func GenerateTokenPair(userID uuid.UUID, email string) (*TokenPair, error) {
	// Get expiration times from env
	accessExpirationHours, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRATION"))
	if err != nil {
		accessExpirationHours = 1 // Default to 1 hour
	}

	refreshExpirationHours, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRATION"))
	if err != nil {
		refreshExpirationHours = 24 * 7 // Default to 7 days
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
	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
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
	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("REFRESH_TOKEN_SECRET")))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    accessExpirationHours * 3600, // Convert to seconds
	}, nil
}

// ValidateToken validates a JWT token
func ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	return validateToken(tokenString, os.Getenv("JWT_SECRET"))
}

// ValidateRefreshToken validates a refresh token
func ValidateRefreshToken(tokenString string) (*JWTClaims, error) {
	return validateToken(tokenString, os.Getenv("REFRESH_TOKEN_SECRET"))
}

func validateToken(tokenString, secret string) (*JWTClaims, error) {
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
