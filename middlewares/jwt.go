package middlewares

import (
	"os"
	"strings"

	"github.com/abdullah-mobin/helloBanglaTTS-backend/response"
	"github.com/abdullah-mobin/helloBanglaTTS-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func IsAuthenticated(c *fiber.Ctx) error {
	var tokenString string

	// 1️⃣ Try to get token from Authorization header
	authHeader := c.Get("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			tokenString = parts[1]
		} else {
			// Wrong format but header exists → return error
			return response.UnauthorizedException(c, "Invalid Authorization header format", nil)
		}
	}

	// 2️⃣ If no token in header → try to get from query param ?token=
	if tokenString == "" {
		tokenString = c.Query("token")
	}

	// 3️⃣ If still empty → unauthorized
	if tokenString == "" {
		return response.UnauthorizedException(c, "Token not provided", nil)
	}

	// 4️⃣ Validate JWT
	userId, err := utils.ValidateToken(tokenString, os.Getenv("JWT_SECRET"))
	if err != nil {
		return response.UnauthorizedException(c, "Invalid or expired token", nil)
	}

	// 5️⃣ store userId in context locals
	c.Locals("userId", userId)

	return c.Next()
}
