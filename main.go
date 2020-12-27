package main

import (
	"log"
	"os"
)

func main() {
	// Create an app struct and initialize the DB connection
	a := App{}
	a.Initialize()

	// Get port from Heroku Environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
		log.Println("Local test")
	}

	// Run the API Server
	a.Run(":" + port)
}
