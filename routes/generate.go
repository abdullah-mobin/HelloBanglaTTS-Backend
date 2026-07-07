package routes

import (
	"github.com/abdullah-mobin/helloBanglaTTS-backend/handlers"
	"github.com/gofiber/fiber/v2"
)

func GenerateRoutes(route fiber.Router) {

	route.Post("/tts", handlers.GenerateTTS)
	// route.Post("/video")
	route.Post("/image", handlers.GenerateImage)
	route.Post("/translate", handlers.TranslateBanglaPrompt)
	route.Post("/story", handlers.GenerateStory)
	route.Post("/humanizer", handlers.HumanizeBanglaTextHandler)
	route.Post("/smartify", handlers.SmartifyBanglaTextHandler)
}
