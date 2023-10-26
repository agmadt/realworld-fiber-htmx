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

func ArticleFavoriteAction(c *fiber.Ctx) error {

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

	return c.Render("articles/partials/favorite-button", fiber.Map{
		"Article":            article,
		"Slug":               article.Slug,
		"IsArticleFavorited": isArticleFavorited,
		"IsOob":              true,
	}, "layouts/app-htmx")
}

func ArticleFollowAction(c *fiber.Ctx) error {

	var article model.Article
	var authenticatedUser model.User

	isFollowed := false

	isAuthenticated, userID := authentication.AuthGet(c)
	if !isAuthenticated {
		return helper.HTMXRedirectTo("/sign-in", "/htmx/sign-in", c)
	}

	db := database.Get()

	err := db.Model(&article).
		Where("slug = ?", c.Params("slug")).
		Preload("Favorites").
		Preload("User.Followers").
		Find(&article).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.HTMXRedirectTo("/sign-in", "/htmx/sign-in", c)
		}
	}

	authenticatedUser.ID = userID

	if article.User.FollowedBy(userID) {

		f := model.Follow{
			FollowerID:  article.UserID,
			FollowingID: userID,
		}

		db.Model(&article.User).Association("Followers").Find(&f)
		db.Delete(&f)

	} else {
		db.Model(&article.User).Association("Followers").Append(&model.Follow{FollowerID: article.UserID, FollowingID: userID})
		isFollowed = true
	}

	return c.Render("articles/partials/follow-button", fiber.Map{
		"Article":    article,
		"IsFollowed": isFollowed,
		"IsOob":      true,
	}, "layouts/app-htmx")
}
