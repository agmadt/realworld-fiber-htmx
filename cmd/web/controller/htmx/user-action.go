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

func UserArticleFavoriteAction(c *fiber.Ctx) error {

	var article model.Article
	var authenticatedUser model.User

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
		article.IsFavorited = false
	} else {
		db.Model(&article).Association("Favorites").Append(&authenticatedUser)
		article.IsFavorited = true
	}

	return c.Render("users/partials/article-favorite-button", fiber.Map{
		"Article": article,
	}, "layouts/app-htmx")
}

func UserFollowAction(c *fiber.Ctx) error {

	var authenticatedUser model.User
	var user model.User
	isFollowed := false

	isAuthenticated, userID := authentication.AuthGet(c)

	db := database.Get()

	err := db.Model(&user).
		Where("username = ?", c.Params("username")).
		Preload("Followers").
		Find(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.HTMXRedirectTo("/", "/htmx/home", c)
		}
	}

	if isAuthenticated {
		db.Model(&authenticatedUser).
			Where("id = ?", userID).
			First(&authenticatedUser)
	}

	f := model.Follow{
		FollowerID:  user.ID,
		FollowingID: userID,
	}

	if user.FollowedBy(userID) {
		db.Model(&user).Association("Followers").Find(&f)
		db.Delete(&f)
	} else {
		db.Model(&user).Association("Followers").Append(&f)
		isFollowed = true
	}

	return c.Render("users/partials/follow-button", fiber.Map{
		"User":       user,
		"IsFollowed": isFollowed,
	}, "layouts/app-htmx")
}
