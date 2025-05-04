package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()

	port := os.Getenv("PORT")
	if port == ""{
		port = "3002" // Default port if not specified
	}
	
	//START SERVER
	log.Printf("SERVER IS RUNNING ON THE PORT  %s\n", port)
	if err := app.Listen("0.0.0.0:"+ port); err !=nil {
		log.Fatalf("failed to start server on port %s: %v", port, err)
	}
}