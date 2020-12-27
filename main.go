package main

import (
	"log"
	"os"
)

func main() {
	// Create an app struct and initialize the DB connection
	a := App{}
	a.Initialize()

	// local test
	// port := "8000"

	// Get port from Heroku Environment
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// Run the API Server
	a.Run(":" + port)
}
