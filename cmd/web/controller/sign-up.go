package controller

import (
	"realworld-fiber-htmx/internal/authentication"

	"github.com/gofiber/fiber/v2"
)

func SignUpPage(c *fiber.Ctx) error {

	isAuthenticated, _ := authentication.AuthGet(c)
	if isAuthenticated {
		return c.Redirect("/")
	}

	return c.Render("sign-up/index", fiber.Map{
		"PageTitle":    "Sign Up — Conduit",
		"FiberCtx":     c,
		"NavBarActive": "sign-up",
	}, "layouts/app")
}
