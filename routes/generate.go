package routes

import (
	"github.com/abdullah-mobin/helloBanglaTTS-backend/handlers"
	"github.com/abdullah-mobin/helloBanglaTTS-backend/middlewares"
	"github.com/gofiber/fiber/v2"
)

func GenerateRoutes(route fiber.Router) {

	route.Post("/tts", middlewares.IsAuthenticated, handlers.GenerateTTS)
	// route.Post("/video")
	// route.Post("/image")
	// route.Post("/story")
	// route.Post("/smartify")
	// route.Post("/humanizer")
}
