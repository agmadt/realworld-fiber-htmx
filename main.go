package main

import (
	"log"
	"realworld-fiber-htmx/cmd/web"
	"realworld-fiber-htmx/internal/authentication"
	"realworld-fiber-htmx/internal/database"
	"realworld-fiber-htmx/internal/renderer"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var validate *validator.Validate

func main() {

	viewEngine := renderer.ViewEngineStart()
	app := fiber.New(fiber.Config{
		Views: viewEngine,
	})

	database.Open()
	authentication.SessionStart()

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())

	web.Serve(app)

	log.Fatal(app.Listen("localhost:8181"))
}
