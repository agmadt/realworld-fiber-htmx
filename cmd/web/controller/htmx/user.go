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

func UserDetailPage(c *fiber.Ctx) error {

	var authenticatedUser model.User
	var user model.User
	isSelf := false
	isFollowed := false
	navbarActive := "none"

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

	if isAuthenticated && user.ID == userID {
		isSelf = true
		navbarActive = "profile"
	}

	if isAuthenticated && !isSelf && user.FollowedBy(userID) {
		isFollowed = true
	}

	return c.Render("users/htmx-users-page", fiber.Map{
		"PageTitle":         user.Name,
		"IsSelf":            isSelf,
		"IsFollowed":        isFollowed,
		"AuthenticatedUser": authenticatedUser,
		"User":              user,
		"NavBarActive":      navbarActive,
		"FiberCtx":          c,
	}, "layouts/app-htmx")
}

func UserArticles(c *fiber.Ctx) error {

	var articles []model.Article
	var user model.User
	hasArticles := false

	_, userID := authentication.AuthGet(c)

	db := database.Get()

	err := db.Where(&user).
		Where("username = ?", c.Params("username")).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.HTMXRedirectTo("/", "/htmx/home", c)
		}
	}

	db.Where(&model.Article{UserID: user.ID}).
		Preload("Favorites").
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Order("tags.name asc")
		}).
		Preload("User").
		Order("created_at desc").
		Find(&articles)

	if len(articles) > 0 {
		hasArticles = true

		for i := 0; i < len(articles); i++ {
			articles[i].IsFavorited = articles[i].FavoritedBy(userID)
		}
	}

	feedNavbarItems := []fiber.Map{
		{
			"Title":     "Articles",
			"IsActive":  true,
			"HXPushURL": "/users/" + user.Username,
			"HXGetURL":  "/htmx/users/" + user.Username,
		},
		{
			"Title":     "Favorited Articles",
			"IsActive":  false,
			"HXPushURL": "/users/" + user.Username + "/favorites",
			"HXGetURL":  "/htmx/users/" + user.Username + "/favorites",
		},
	}

	return c.Render("users/htmx-users-articles", fiber.Map{
		"HasArticles":     hasArticles,
		"Articles":        articles,
		"User":            user,
		"FeedNavbarItems": feedNavbarItems,
	}, "layouts/app-htmx")
}

func UserArticlesFavorite(c *fiber.Ctx) error {

	var articles []model.Article
	var user model.User
	isSelf := false
	isFollowed := false
	hasArticles := false

	isAuthenticated, userID := authentication.AuthGet(c)

	db := database.Get()

	err := db.Model(&user).
		Where("username = ?", c.Params("username")).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.HTMXRedirectTo("/", "/htmx/home", c)
		}
	}

	db.Model(&user).
		Preload("Favorites").
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Order("tags.name asc")
		}).
		Preload("User").
		Order("created_at desc").
		Association("Favorites").
		Find(&articles)

	if len(articles) > 0 {
		hasArticles = true

		for i := 0; i < len(articles); i++ {
			articles[i].IsFavorited = articles[i].FavoritedBy(userID)
		}
	}

	if isAuthenticated && user.ID == userID {
		isSelf = true
	}

	if isAuthenticated && !isSelf && user.FollowedBy(userID) {
		isFollowed = true
	}

	feedNavbarItems := []fiber.Map{
		{
			"Title":     "Articles",
			"IsActive":  false,
			"HXPushURL": "/users/" + user.Username + "/articles",
			"HXGetURL":  "/htmx/users/" + user.Username + "/articles",
		},
		{
			"Title":     "Favorited Articles",
			"IsActive":  true,
			"HXPushURL": "/users/" + user.Username + "/favorites",
			"HXGetURL":  "/htmx/users/" + user.Username + "/favorites",
		},
	}

	return c.Render("users/htmx-users-articles", fiber.Map{
		"IsSelf":          isSelf,
		"IsFollowed":      isFollowed,
		"HasArticles":     hasArticles,
		"Articles":        articles,
		"User":            user,
		"FeedNavbarItems": feedNavbarItems,
		"IsLoadFavorites": true,
	}, "layouts/app-htmx")
}
