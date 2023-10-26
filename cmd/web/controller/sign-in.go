package controller

import (
	"realworld-fiber-htmx/internal/authentication"

	"github.com/gofiber/fiber/v2"
)

func SignInPage(c *fiber.Ctx) error {

	isAuthenticated, _ := authentication.AuthGet(c)
	if isAuthenticated {
		return c.Redirect("/")
	}

	return c.Render("sign-in/index", fiber.Map{
		"PageTitle":    "Sign In — Conduit",
		"FiberCtx":     c,
		"NavBarActive": "sign-in",
	}, "layouts/app")
}
