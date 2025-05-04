package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber"
)

func main() {
	
	app := fiber.New()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3002" // Default port if not specified
	}

	//START SERVER
	log.Printf("Server is running on port %s\n", port)
	if err := app.Listen("0.0.0.0:"+ port); err !=nil {
		log.Fatalf("failed to start server on port %s: %v", port, err)
	}
}