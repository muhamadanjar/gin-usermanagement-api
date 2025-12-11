package main

import (
	"log"
	"usermanagement-api/internal/app"
	"usermanagement-api/internal/container"
)

func main() {
	// Initialize App Container (Infrastructure layer)
	appContainer, err := app.NewAppContainer()
	if err != nil {
		log.Fatalf("Failed to initialize app container: %v", err)
	}
	defer appContainer.Close()

	// Initialize Business Container (Domain layer)
	businessContainer := container.NewBusinessContainer(
		appContainer.DB,
		appContainer.Cache,
		appContainer.FCMClient,
	)

	// Initialize Server with containers
	server := app.NewServer(appContainer, businessContainer)

	if err := server.Initialize(); err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	if err := server.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
