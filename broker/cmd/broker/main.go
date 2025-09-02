package main

import (
	"log"
)

func main() {
	app, err := NewApp()
	if err != nil {
		log.Fatalf("Error initializing app: %v", err)
	}

	app.Run()

	// Keep the program running
	select {}
}
