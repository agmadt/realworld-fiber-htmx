package controller

import (
	"realworld-fiber-htmx/cmd/web/model"
	"realworld-fiber-htmx/internal/authentication"
	"realworld-fiber-htmx/internal/database"

	"github.com/gofiber/fiber/v2"
)

func SettingPage(c *fiber.Ctx) error {

	var authenticatedUser model.User

	isAuthenticated, userID := authentication.AuthGet(c)
	if !isAuthenticated {
		return c.Redirect("/")
	}

	db := database.Get()

	db.Model(&authenticatedUser).
		Where("id = ?", userID).
		First(&authenticatedUser)

	return c.Render("settings/index", fiber.Map{
		"PageTitle":         "Settings — Conduit",
		"FiberCtx":          c,
		"NavBarActive":      "settings",
		"AuthenticatedUser": authenticatedUser,
	}, "layouts/app")
}
