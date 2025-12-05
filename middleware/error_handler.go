package middleware

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
)

// ValidationError represents a validation error with a map of field errors
type ValidationError struct {
	Message string
	Errors  map[string]string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// ErrorHandler untuk menangani error secara global dan mengembalikan response dengan standardized format
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Default status code and message
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	// Check for specific error types
	var e *fiber.Error
	var ve *ValidationError

	if errors.As(err, &ve) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   ve.Message,
			"message": ve.Errors,
		})
	} else if errors.As(err, &e) {
		code = e.Code
		message = e.Message
	} else {
		// Log internal errors that are not Fiber errors
		log.Printf("Internal Error: %v", err)
	}

	// Return JSON response
	return c.Status(code).JSON(fiber.Map{
		"status":  "error",
		"message": message,
	})
}
