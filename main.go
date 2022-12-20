// Main package is used to start the api-gateway
package main

import (
	"api-gateway/database"
	"api-gateway/logger"
	"api-gateway/router"
	"api-gateway/server"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Starts the api gateway and all the services that work in it
func main() {
	godotenv.Load()

	errorLogger := logger.CreateNewLogger()
	if errorLogger != nil {
		log.Fatalf("Failed to Create Logger: %v", errorLogger)
	}

	client, errorDatabase := database.NewDatabase()
	if errorDatabase != nil {
		logger.Log.Errorf("Failed to Create DB : %v", errorDatabase)
	}

	defer client.Def.Close()
	myServer := server.NewServer()

	router.Apply(myServer.Def)
	myServer.Def.Run(":" + os.Getenv("API_GATEWAY_PORT"))
}
