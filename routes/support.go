package routes

import (
	"github.com/abdullah-mobin/helloBanglaTTS-backend/handlers"
	"github.com/abdullah-mobin/helloBanglaTTS-backend/middlewares"
	"github.com/gofiber/fiber/v2"
)

func SupportRoutes(route fiber.Router) {

	ticketRoutes(route.Group("/ticket"))

}

func ticketRoutes(route fiber.Router) {
	route.Get("/", middlewares.IsAuthenticated, handlers.FindTickets)
	route.Post("/new", middlewares.IsAuthenticated, handlers.CreateTicket)
	route.Put("/:id", middlewares.IsAuthenticated, handlers.UpdateTicket)
	route.Delete("/:id", middlewares.IsAuthenticated, handlers.DeleteTicket)
}
