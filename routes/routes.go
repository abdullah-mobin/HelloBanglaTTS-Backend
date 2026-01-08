package routes

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	AuthRoutes(api.Group("/auth"))
	UserRoutes(api.Group("/user"))
	GenerateRoutes(api.Group("/generate"))
	ApiKeyRoutes(api.Group("/api-keys"))
	SupportRoutes(api.Group("/support"))

}
