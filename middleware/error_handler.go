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

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type SpecificErrorResponse struct {
	Error   string            `json:"error"`
	Message map[string]string `json:"message"`
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
		// Return JSON response for Validation Error
		return c.Status(fiber.StatusBadRequest).JSON(SpecificErrorResponse{
			Error:   ve.Message,
			Message: ve.Errors,
		})
	} else if errors.As(err, &e) {
		code = e.Code
		message = e.Message
	} else {
		// Log internal errors that are not Fiber errors
		log.Printf("Internal Error: %v", err)
	}

	// Return JSON response for Internal Server Error
	return c.Status(code).JSON(ErrorResponse{
		Status:  "error",
		Message: message,
	})
}
