package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// RateLimiter membatasi jumlah request dari satu IP
func RateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		// Konfigurasi rate limiter - 20 request per 30 detik
		Max:        20,
		Expiration: 30 * time.Second,

		// Menggunakan IP address sebagai key
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},

		// Respon ketika batas tercapai
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"status":  "Too Many Requests",
				"message": "Terlalu banyak request. Silakan coba lagi nanti.",
			})
		},
	})
}
