package main

import (
	"log"
	"os"

	"warehouse-inventory-server/config"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	port := os.Getenv(`PORT`)

	if err != nil {
		log.Fatalf("Error loading .env file: %v", err.Error())
	}

	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err.Error())
	}

	log.Println("database connected")
	_ = db

	app := fiber.New()

	// Server health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "Healthy",
		})
	})

	app.Listen(":" + port)
}
