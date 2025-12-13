package app

import (
	"usermanagement-api/config"
	"usermanagement-api/pkg/auth"
	"usermanagement-api/pkg/cache"
	"usermanagement-api/pkg/database"
	"usermanagement-api/pkg/firebase"
	"usermanagement-api/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AppContainer holds all infrastructure dependencies
type AppContainer struct {
	Config    *config.Config
	Logger    *zap.Logger
	DB        *gorm.DB
	Cache     cache.Cache
	FCMClient firebase.FCMClient
}

// NewAppContainer creates and initializes a new AppContainer
func NewAppContainer() (*AppContainer, error) {
	// Load configuration first (needed for logger config)
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	// Initialize logger
	if err := logger.Initialize(cfg.Logger); err != nil {
		return nil, err
	}
	zapLogger := logger.GetLogger()
	zapLogger.Info("Logger initialized", zap.String("level", cfg.Logger.Level), zap.String("mode", cfg.Logger.Mode))

	// Initialize Redis cache
	cacheInstance := cache.NewRedisCache(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	zapLogger.Info("Redis cache initialized", zap.String("addr", cfg.Redis.Addr))

	// Initialize Firebase Cloud Messaging client (optional)
	var fcmClient firebase.FCMClient
	if cfg.Firebase.CredentialsFile != "" {
		fcm, err := firebase.NewFCMClient(cfg.Firebase.CredentialsFile, zapLogger)
		if err != nil {
			zapLogger.Warn("Failed to initialize FCM client", zap.Error(err))
		} else {
			fcmClient = fcm
			zapLogger.Info("FCM client initialized")
		}
	} else {
		zapLogger.Info("FCM client skipped (no credentials file configured)")
	}

	// Connect to database
	db, err := database.ConnectDB(cfg, zapLogger)
	if err != nil {
		zapLogger.Fatal("Failed to connect to database", zap.Error(err))
		return nil, err
	}

	// Migrate database
	if err := database.MigrateDB(db, zapLogger); err != nil {
		zapLogger.Fatal("Failed to migrate database", zap.Error(err))
		return nil, err
	}

	// Initialize JWT service and set as global for backward compatibility
	jwtService := auth.NewJWTService(cfg.JWT)
	auth.SetGlobalJWTService(jwtService)

	return &AppContainer{
		Config:    cfg,
		Logger:    zapLogger,
		DB:        db,
		Cache:     cacheInstance,
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
		if err := sqlDB.Close(); err != nil {
			c.Logger.Error("Failed to close database connection", zap.Error(err))
			return err
		}
		c.Logger.Info("Database connection closed")
	}

	// Sync logger before exit
	if err := logger.Sync(); err != nil {
		return err
	}

	return nil
}
