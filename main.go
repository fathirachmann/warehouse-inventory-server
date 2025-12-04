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

	// Auth routes
	userRepo := repositories.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepo)
	authRoute := app.Group("/api/auth")
	userHandler.RegisterRoute(authRoute)

	// Barang routes
	barangRepo := repositories.NewBarangRepository(db)
	barangHandler := handlers.NewBarangHandler(barangRepo)
	barangRoute := app.Group("/api/barang", middleware.JWTProtected())
	barangHandler.RegisterRoute(barangRoute)

	// 2. Implement: Stock routes (Guarded with JWT middleware)
	// GET all stock - /api/stok
	// GET stock by barang ID - /api/stok/:barang_id
	// GET History stock - /api/history-stok
	// GET History stock by barang ID - /api/history-stok/:barang_id

	// 3. Implement: Transaksi - Pembelian routes (Guarded with JWT middleware)
	// POST Create pembelian - /api/pembelian
	// GET all pembelian - /api/pembelian
	// GET pembelian by ID - /api/pembelian/:id

	// 4, Implement: Transaksi - Penjualan routes (Guarded with JWT middleware)
	// POST Create penjualan - /api/penjualan
	// GET all penjualan - /api/penjualan
	// GET penjualan by ID - /api/penjualan/:id

	port := os.Getenv("PORT")

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
