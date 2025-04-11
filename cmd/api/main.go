package main

import (
	"log"
	"usermanagement-api/internal/app"
)

func main() {
	server := app.NewServer()

	if err := server.Initialize(); err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	if err := server.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
