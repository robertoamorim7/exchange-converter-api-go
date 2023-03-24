package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	API_KEY = os.Getenv("API_KEY")

	app := fiber.New()

	InitRoutes(app)

	app.Listen(":8080")
}
