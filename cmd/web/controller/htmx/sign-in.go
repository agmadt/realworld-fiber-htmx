package HTMXController

import (
	"errors"
	"realworld-fiber-htmx/cmd/web/model"
	"realworld-fiber-htmx/internal/authentication"
	"realworld-fiber-htmx/internal/database"
	"realworld-fiber-htmx/internal/helper"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SignInPage(c *fiber.Ctx) error {

	return c.Render("sign-in/htmx-sign-in-page", fiber.Map{
		"PageTitle":    "Sign In",
		"NavBarActive": "sign-in",
		"FiberCtx":     c,
	}, "layouts/app-htmx")

}

func SignInAction(c *fiber.Ctx) error {

	var user model.User
	email := c.FormValue("email")
	password := c.FormValue("password")

	if email == "" || password == "" {

		return c.Render("sign-in/partials/sign-in-form", fiber.Map{
			"Errors": []string{
				"Email or password cannot be null.",
			},
			"IsOob": true,
		}, "layouts/app-htmx")
	}

	db := database.Get()

	db.Model(&user)
	err := db.Where(&model.User{Email: email}).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Render("sign-in/partials/sign-in-form", fiber.Map{
				"Errors": []string{
					"Email and password did not match.",
				},
			}, "layouts/app-htmx")
		}
	}

	if !user.CheckPassword(password) {

		return c.Render("sign-in/partials/sign-in-form", fiber.Map{
			"Errors": []string{
				"Email and password did not match.",
			},
		}, "layouts/app-htmx")
	}

	authentication.AuthStore(c, user.ID)

	return helper.HTMXRedirectTo("/", "/htmx/home", c)
}

func SignOut(c *fiber.Ctx) error {

	isAuthenticated, _ := authentication.AuthGet(c)
	if !isAuthenticated {
		return c.Redirect("/")
	}

	authentication.AuthDestroy(c)

	return helper.HTMXRedirectTo("/", "/htmx/home", c)
}
