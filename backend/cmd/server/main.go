package main

import (
	"log"

	"github.com/gsbcamargo/freenstagram/backend/internal/db"
	"github.com/gsbcamargo/freenstagram/backend/internal/routes"
)

func main() {
	// Initialize the database connection.
	if err := db.Init(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Set up the routes.
	router := routes.SetupRouter()

	// Start the server on port 8080.
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
