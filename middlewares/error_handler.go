package middlewares

import (
	"github.com/abdullah-mobin/helloBanglaTTS-backend/response"
	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	// Default to 500 Internal Server Error
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	// Check if it's a Fiber error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	} else if err != nil {
		message = err.Error()
	}

	// Return consistent JSON response
	return response.Custom(c, code, message, err)
}
