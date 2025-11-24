package main

import (
	"log"
	"os"

	"github.com/hieu9721/media-store-backend/config"
	"github.com/hieu9721/media-store-backend/routes"
	"github.com/joho/godotenv"
)

func main() {
    // Load .env file
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    // Connect to MongoDB
    config.ConnectDB()

    // Setup routes
    router := routes.SetupRoutes()

    // Get port from env
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    // Start server
    log.Printf("🚀 Server starting on port %s...", port)
    if err := router.Run(":" + port); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
