package middleware

import (
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Authentication() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Header Authorization error",
				"message": "header Authorization tidak ditemukan"})
		}

		// Expect format: Bearer <token>
		var tokenString string
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Header Authorization error",
				"message": "Token tidak valid atau hilang",
			})
		}

		secret := os.Getenv("JWT_SECRET")

		if secret == "" {
			log.Println("Environment variable error: JWT_SECRET is empty", "middleware.go:Authentication", "Error at line 32")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server error"})
		}

		parsed, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(secret), nil
		})
		if err != nil || !parsed.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "	"})
		}

		// Simpan claims ke fiber context
		if claims, ok := parsed.Claims.(jwt.MapClaims); ok {
			c.Locals("user", claims)
		}

		return c.Next()
	}
}

// GuardAdmin memastikan bahwa hanya user "admin" yang dapat mengakses route tertentu
func GuardAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil claims dari context
		userClaims, ok := c.Locals("user").(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "JWT claims error",
				"message": "claims user tidak ditemukan dalam context",
			})
		}

		// Ambil role dari claims
		roleVal, ok := userClaims["role"]
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   "JWT claims error",
				"message": "role tidak ditemukan dalam token",
			})
		}

		// Validasi tipe data role
		roleStr, ok := roleVal.(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   "JWT claims error",
				"message": "format role dalam token tidak valid",
			})
		}

		// Case validation: Cek apakah role adalah "admin"
		if strings.ToLower(roleStr) != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   "Akses ditolak",
				"message": "Tidak memiliki izin - admin only",
			})
		}

		// Lolos
		return c.Next()
	}
}

// Note: Hanya ada satu authorization GuardAdmin karena hanya ada dua peran (admin dan staff). Admin juga bisa mengakses semua route staff.
