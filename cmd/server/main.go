package main

import (
	"httpapp/internal/server/serverapp"
	"log"
)

const (
	addr = ":8080"
)

func main() {
	app := serverapp.NewServerApp(addr)
	if err := app.Start(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	log.Println("Server stopped successfully")
}
