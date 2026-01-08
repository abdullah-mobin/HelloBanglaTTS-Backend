package routes

import (
	"github.com/abdullah-mobin/helloBanglaTTS-backend/handlers"
	"github.com/abdullah-mobin/helloBanglaTTS-backend/middlewares"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(route fiber.Router) {
	route.Post("/register", handlers.RegisterNewUser)
	route.Post("/login", handlers.Login)
	route.Post("/refresh", handlers.RefreshToken)
	route.Post("/logout", middlewares.IsAuthenticated, handlers.Logout)
}
