package controller

import (
	"realworld-fiber-htmx/cmd/web/model"
	"realworld-fiber-htmx/internal/authentication"
	"realworld-fiber-htmx/internal/database"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func EditorPage(c *fiber.Ctx) error {

	var authenticatedUser model.User
	var article model.Article
	hasArticle := false

	isAuthenticated, userID := authentication.AuthGet(c)
	if !isAuthenticated {
		return c.Redirect("/")
	}

	db := database.Get()

	db.Model(&authenticatedUser).
		Where("id = ?", userID).
		First(&authenticatedUser)

	if c.Params("slug") != "" {

		err := db.Model(&article).
			Where("slug = ?", c.Params("slug")).
			Preload("Tags", func(db *gorm.DB) *gorm.DB {
				return db.Order("tags.name asc")
			}).
			Find(&article).Error

		if err == nil {
			hasArticle = true
		}
	}

	return c.Render("editor/form", fiber.Map{
		"PageTitle":         "Editor — Conduit",
		"FiberCtx":          c,
		"NavBarActive":      "editor",
		"AuthenticatedUser": authenticatedUser,
		"HasArticle":        hasArticle,
		"Article":           article,
	}, "layouts/app")
}
