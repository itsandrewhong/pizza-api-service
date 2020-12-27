package main

import (
	"log"
	"os"
)

func main() {
	a := App{}
	a.Initialize()

	// local test
	// port := "8000"

	// Get port from Heroku Environment
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	a.Run(":" + port)
}
