package main

import (
	"log"

	"zeneye-gateway/internal/adapter/http"
	"zeneye-gateway/internal/adapter/repository/postgres"
	"zeneye-gateway/pkg/logger"
	"zeneye-gateway/pkg/utils"
)

func main() {

	// Load Env
	utils.LoadConfig()

	// Initialize and ensure logs are flushed
	logger.InitLogger()
	defer logger.SyncLogger()

	db := postgres.InitDB() // Initialize the database

	// Setup and run the HTTP router
	router := http.SetupRouter(db)
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
