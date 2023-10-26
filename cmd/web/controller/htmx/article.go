package HTMXController

import (
	"errors"
	"realworld-fiber-htmx/cmd/web/model"
	"realworld-fiber-htmx/internal/authentication"
	"realworld-fiber-htmx/internal/database"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ArticleDetailPage(c *fiber.Ctx) error {

	var article model.Article
	isSelf := false
	isFollowed := false
	var authenticatedUser model.User

	isAuthenticated, userID := authentication.AuthGet(c)

	db := database.Get()

	if isAuthenticated {
		db.Model(&authenticatedUser).
			Where("id = ?", userID).
			First(&authenticatedUser)
	}

	err := db.Model(&article).
		Where("slug = ?", c.Params("slug")).
		Preload("Favorites").
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Order("tags.name asc")
		}).
		Preload("User").
		Find(&article).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Redirect("/")
		}
	}

	if isAuthenticated && article.User.FollowedBy(userID) {
		isFollowed = true
	}

	if isAuthenticated && article.User.ID == userID {
		isSelf = true
	}

	return c.Render("articles/htmx-article-page", fiber.Map{
		"PageTitle":          article.Title,
		"NavBarActive":       "none",
		"Article":            article,
		"IsOob":              false,
		"IsSelf":             isSelf,
		"IsFollowed":         isFollowed,
		"IsArticleFavorited": article.FavoritedBy(userID),
		"AuthenticatedUser":  authenticatedUser,
		"FiberCtx":           c,
	}, "layouts/app-htmx")
}
