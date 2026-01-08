package main

import (
	"fmt"
	"log"
	"os"

	"github.com/abdullah-mobin/helloBanglaTTS-backend/config"
	"github.com/abdullah-mobin/helloBanglaTTS-backend/database"
	"github.com/abdullah-mobin/helloBanglaTTS-backend/middlewares"
	"github.com/abdullah-mobin/helloBanglaTTS-backend/routes"

	_ "github.com/abdullah-mobin/helloBanglaTTS-backend/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
)

//	@title						HelloBanglaTTS-Backend-APIs
//	@version					1.0
//	@description				Hello Bangla TTS is a Bangla language based Ai platform writen in go.
//	@host						localhost:5001
//	@BasePath					/api/v1
//	@schemes					http https
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization

func main() {
	err := config.LoadENV()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = database.InitCollections()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.CloseDB()

	app := fiber.New(fiber.Config{
		ErrorHandler: middlewares.ErrorHandler,
	})

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(helmet.New())
	app.Use(compress.New())
	app.Use(healthcheck.New())
	app.Get("/swagger/*", swagger.HandlerDefault)

	routes.SetupRoutes(app)
	port := os.Getenv("APP_PORT")
	fmt.Println("Server is running on port", port)
	log.Fatal(app.Listen(":" + port))
}
