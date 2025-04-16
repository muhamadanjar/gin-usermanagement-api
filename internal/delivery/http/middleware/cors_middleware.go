package middleware

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware interface
type CORSMiddleware interface {
	SetupCORS() gin.HandlerFunc
}

type corsMiddleware struct{}

// NewCORSMiddleware creates a new CORS middleware
func NewCORSMiddleware() CORSMiddleware {
	return &corsMiddleware{}
}

// SetupCORS sets up CORS configuration from environment variables
func (m *corsMiddleware) SetupCORS() gin.HandlerFunc {
	// Get CORS configuration from environment
	allowedOriginsStr := os.Getenv("CORS_ALLOWED_ORIGINS")
	allowedOrigins := []string{"http://localhost:3000"} // Default value
	if allowedOriginsStr != "" {
		allowedOrigins = strings.Split(allowedOriginsStr, ",")
	}

	allowCredentials := true // Default value
	if os.Getenv("CORS_ALLOW_CREDENTIALS") == "false" {
		allowCredentials = false
	}

	maxAgeStr := os.Getenv("CORS_MAX_AGE")
	maxAge := 12 * time.Hour // Default value
	if maxAgeStr != "" {
		if maxAgeSeconds, err := strconv.Atoi(maxAgeStr); err == nil {
			maxAge = time.Duration(maxAgeSeconds) * time.Second
		}
	}

	// Create CORS configuration
	return cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: allowCredentials,
		MaxAge:           maxAge,
	})
}
