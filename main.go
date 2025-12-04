package main

import (
	"log"
	"os"

	"warehouse-inventory-server/config"
	"warehouse-inventory-server/handlers"
	"warehouse-inventory-server/middleware"
	"warehouse-inventory-server/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	_ = godotenv.Load()

	// Initialize Database
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	log.Println("database connected")

	// Initialize Fiber app
	app := fiber.New()

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "healthy",
			"message": "Warehouse Inventory Server is running",
		})
	})

	// Auth routes
	userRepo := repositories.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepo)

	authRoute := app.Group("/api/auth")
	userHandler.RegisterRoute(authRoute)

	// Barang routes (admin only for write operations implied by handler, but we guard whole group as admin)
	barangRepo := repositories.NewBarangRepository(db)
	barangHandler := handlers.NewBarangHandler(barangRepo)

	barangRoute := app.Group("/api/barang", middleware.Authentication())
	barangHandler.RegisterRoute(barangRoute)

	// Stock routes (viewable by both admin and staff)
	stokRepo := repositories.NewStokRepository(db)
	stokHandler := handlers.NewStokHandler(stokRepo)

	stokRoute := app.Group("/api/stok", middleware.Authentication())
	stokHandler.RegisterStockRoute(stokRoute)

	historyRoute := app.Group("/api/history-stok", middleware.Authentication())
	stokHandler.RegisterHistoryRoute(historyRoute)

	// Pembelian routes (admin and staff allowed per requirement)
	pembelianRepo := repositories.NewPembelianRepository(db)
	pembelianHandler := handlers.NewPembelianHandler(pembelianRepo, stokRepo)

	pembelianRoute := app.Group("/api/pembelian", middleware.Authentication())
	pembelianHandler.RegisterRoute(pembelianRoute)

	// Penjualan routes (admin and staff allowed per requirement)
	penjualanRepo := repositories.NewPenjualanRepository(db)
	penjualanHandler := handlers.NewPenjualanHandler(penjualanRepo, stokRepo)

	penjualanRoute := app.Group("/api/penjualan", middleware.Authentication())
	penjualanHandler.RegisterRoute(penjualanRoute)

	port := os.Getenv("PORT")

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
