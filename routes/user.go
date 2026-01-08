package routes

import (
	"github.com/abdullah-mobin/helloBanglaTTS-backend/handlers"
	"github.com/abdullah-mobin/helloBanglaTTS-backend/middlewares"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(route fiber.Router) {
	route.Get("/profile", middlewares.IsAuthenticated, handlers.UserProfile)
	route.Get("/", middlewares.IsAuthenticated, handlers.Users)
}
