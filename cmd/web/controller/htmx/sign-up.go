package HTMXController

import (
	"realworld-fiber-htmx/cmd/web/model"
	"realworld-fiber-htmx/internal/authentication"
	"realworld-fiber-htmx/internal/database"
	"realworld-fiber-htmx/internal/helper"

	"github.com/gofiber/fiber/v2"
)

func SignUpPage(c *fiber.Ctx) error {

	return c.Render("sign-up/htmx-sign-up-page", fiber.Map{
		"PageTitle":    "Sign Up",
		"NavBarActive": "sign-up",
		"FiberCtx":     c,
	}, "layouts/app-htmx")
}

func SignUpAction(c *fiber.Ctx) error {

	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")

	if email == "" || username == "" || password == "" {

		return c.Render("sign-up/partials/sign-up-form", fiber.Map{
			"Errors": []string{
				"Username, email, and password cannot be null.",
			},
			"IsOob": true,
		}, "layouts/app-htmx")
	}

	user := model.User{Username: username, Email: email, Password: password, Name: username}
	user.HashPassword()

	db := database.Get()
	db.Create(&user)

	authentication.AuthStore(c, user.ID)

	return helper.HTMXRedirectTo("/", "/htmx/home", c)
}
