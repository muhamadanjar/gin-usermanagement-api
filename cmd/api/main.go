package main

import (
	"usermanagement-api/internal/app"
	"usermanagement-api/internal/container"
	"usermanagement-api/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	// Initialize App Container (Infrastructure layer)
	// Logger will be initialized inside NewAppContainer based on config
	appContainer, err := app.NewAppContainer()
	if err != nil {
		// If logger is not initialized yet, use a fallback
		log := logger.GetLogger()
		if log == nil {
			// Last resort: initialize default logger
			logger.Initialize(logger.Config{Level: "info", Mode: "development"})
			log = logger.GetLogger()
		}
		log.Fatal("Failed to initialize app container", zap.Error(err))
	}
	defer func() {
		if err := appContainer.Close(); err != nil {
			appContainer.Logger.Error("Error closing app container", zap.Error(err))
		}
	}()

	// Initialize Business Container (Domain layer)
	businessContainer := container.NewBusinessContainer(
		appContainer.DB,
		appContainer.Cache,
		appContainer.FCMClient,
		appContainer.Config.CORS,
	)

	// Initialize Server with containers
	server := app.NewServer(appContainer, businessContainer)

	if err := server.Initialize(); err != nil {
		appContainer.Logger.Fatal("Failed to initialize server", zap.Error(err))
	}

	if err := server.Run(); err != nil {
		appContainer.Logger.Fatal("Failed to run server", zap.Error(err))
	}
}
