package main

import (
	"log"
	"os"

	"warehouse-inventory-server/config"
	"warehouse-inventory-server/handlers"
	"warehouse-inventory-server/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load env (optional, don’t fail hard on missing)
	_ = godotenv.Load()

	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	log.Println("database connected")

	app := fiber.New()

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "Healthy"})
	})

	// Barang routes
	barangRepo := repositories.NewBarangRepository(db)
	barangHandler := handlers.NewBarangHandler(barangRepo)
	barangRoute := app.Group("/api/barang")
	barangHandler.RegisterRoute(barangRoute)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
