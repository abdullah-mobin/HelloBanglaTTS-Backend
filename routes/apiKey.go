package routes

import (
	"github.com/abdullah-mobin/helloBanglaTTS-backend/handlers"
	"github.com/abdullah-mobin/helloBanglaTTS-backend/middlewares"
	"github.com/gofiber/fiber/v2"
)

func ApiKeyRoutes(route fiber.Router) {
	route.Get("/", middlewares.IsAuthenticated, handlers.FindApiKeys)
	route.Post("/new", middlewares.IsAuthenticated, handlers.CreateApiKey)
	route.Put("/:key", middlewares.IsAuthenticated, handlers.UpdateApiKey)
}
