package controller

import (
	"realworld-fiber-htmx/cmd/web/model"
	"realworld-fiber-htmx/internal/authentication"
	"realworld-fiber-htmx/internal/database"

	"github.com/gofiber/fiber/v2"
)

func HomePage(c *fiber.Ctx) error {

	var authenticatedUser model.User

	isAuthenticated, userID := authentication.AuthGet(c)

	if isAuthenticated {
		db := database.Get()
		db.Model(&authenticatedUser).
			Where("id = ?", userID).
			First(&authenticatedUser)
	}

	return c.Render("home/index", fiber.Map{
		"PageTitle":         "Home — Conduit",
		"FiberCtx":          c,
		"NavBarActive":      "home",
		"AuthenticatedUser": authenticatedUser,
		"CurrentPage":       c.QueryInt("page"),
	}, "layouts/app")
}

func YourFeedPage(c *fiber.Ctx) error {

	var authenticatedUser model.User

	isAuthenticated, userID := authentication.AuthGet(c)
	if isAuthenticated {
		db := database.Get()

		db.Model(&authenticatedUser).
			Where("id = ?", userID).
			First(&authenticatedUser)
	} else {
		return c.Redirect("/")
	}

	return c.Render("home/index", fiber.Map{
		"PageTitle":         "Home — Conduit",
		"Personal":          true,
		"FiberCtx":          c,
		"NavBarActive":      "home",
		"AuthenticatedUser": authenticatedUser,
		"CurrentPage":       c.QueryInt("page"),
	}, "layouts/app")
}

func TagFeedPage(c *fiber.Ctx) error {

	var user model.User

	isAuthenticated, userID := authentication.AuthGet(c)

	if isAuthenticated {
		db := database.Get()
		db.Model(&model.User{ID: userID}).
			First(&user)
	}

	return c.Render("home/index", fiber.Map{
		"PageTitle":    "Home — Conduit",
		"Tag":          true,
		"TagSlug":      c.Params("slug"),
		"FiberCtx":     c,
		"NavBarActive": "home",
		"User":         user,
		"CurrentPage":  c.QueryInt("page"),
	}, "layouts/app")
}
