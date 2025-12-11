package app

import (
	"log"
	"os"
	"strconv"
	"usermanagement-api/config"
	"usermanagement-api/pkg/cache"
	"usermanagement-api/pkg/database"
	"usermanagement-api/pkg/firebase"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// AppContainer holds all infrastructure dependencies
type AppContainer struct {
	Config    *config.Config
	DB        *gorm.DB
	Cache     cache.Cache
	FCMClient firebase.FCMClient
}

// NewAppContainer creates and initializes a new AppContainer
func NewAppContainer() (*AppContainer, error) {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	// Initialize Redis cache
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	cache := cache.NewRedisCache(redisAddr, redisPassword, redisDB)

	// Initialize Firebase Cloud Messaging client (optional)
	var fcmClient firebase.FCMClient
	fcmCredentialsFile := os.Getenv("FIREBASE_CREDENTIALS_FILE")
	if fcmCredentialsFile != "" {
		fcm, err := firebase.NewFCMClient(fcmCredentialsFile)
		if err != nil {
			log.Printf("Warning: Failed to initialize FCM client: %v", err)
		} else {
			fcmClient = fcm
		}
	}

	// Connect to database
	db := database.ConnectDB()

	// Migrate database
	database.MigrateDB(db)

	return &AppContainer{
		Config:    cfg,
		DB:        db,
		Cache:     cache,
		FCMClient: fcmClient,
	}, nil
}

// Close closes all connections in the container
func (c *AppContainer) Close() error {
	if c.DB != nil {
		sqlDB, err := c.DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
