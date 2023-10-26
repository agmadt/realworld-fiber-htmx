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

func HomeFavoriteAction(c *fiber.Ctx) error {

	var article model.Article
	var authenticatedUser model.User

	isArticleFavorited := false

	isAuthenticated, userID := authentication.AuthGet(c)
	if !isAuthenticated {
		return helper.HTMXRedirectTo("/sign-in", "/htmx/sign-in", c)
	}

	db := database.Get()

	err := db.Model(&article).
		Where("slug = ?", c.Params("slug")).
		Preload("Favorites").
		Find(&article).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.HTMXRedirectTo("/sign-in", "/htmx/sign-in", c)
		}
	}

	authenticatedUser.ID = userID

	if article.FavoritedBy(userID) {
		db.Model(&article).Association("Favorites").Delete(&authenticatedUser)
	} else {
		db.Model(&article).Association("Favorites").Append(&authenticatedUser)
		isArticleFavorited = true
	}

	return c.Render("home/partials/article-favorite-button", fiber.Map{
		"GetFavoriteCount": article.GetFavoriteCount(),
		"Slug":             article.Slug,
		"IsFavorited":      isArticleFavorited,
	}, "layouts/app-htmx")
}
