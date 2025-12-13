package middleware

import (
	"time"
	"usermanagement-api/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware interface
type CORSMiddleware interface {
	SetupCORS() gin.HandlerFunc
}

type corsMiddleware struct {
	corsConfig config.CORSConfig
}

// NewCORSMiddleware creates a new CORS middleware
func NewCORSMiddleware(corsConfig config.CORSConfig) CORSMiddleware {
	return &corsMiddleware{
		corsConfig: corsConfig,
	}
}

// SetupCORS sets up CORS configuration from config
func (m *corsMiddleware) SetupCORS() gin.HandlerFunc {
	maxAge := time.Duration(m.corsConfig.MaxAge) * time.Second

	// Create CORS configuration
	return cors.New(cors.Config{
		AllowOrigins:     m.corsConfig.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: m.corsConfig.AllowCredentials,
		MaxAge:           maxAge,
	})
}