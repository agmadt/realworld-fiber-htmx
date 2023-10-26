package web

import (
	webroute "realworld-fiber-htmx/cmd/web/route"
	"realworld-fiber-htmx/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func Serve(app *fiber.App) {

	app.Static("/static", "./cmd/web/static")

	webroute.WebHandlers(app)
	webroute.HTMXHandlers(app)

	// Handle not founds
	app.Use(middleware.NotFound)
}
