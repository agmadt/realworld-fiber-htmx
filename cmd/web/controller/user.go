package controller

import (
	"errors"
	"realworld-fiber-htmx/cmd/web/model"
	"realworld-fiber-htmx/internal/authentication"
	"realworld-fiber-htmx/internal/database"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UserDetailPage(c *fiber.Ctx) error {

	var user model.User
	var authenticatedUser model.User
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
			return c.Redirect("/")
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

	return c.Render("users/show", fiber.Map{
		"PageTitle":         user.Name + " — Conduit",
		"FiberCtx":          c,
		"IsSelf":            isSelf,
		"IsFollowed":        isFollowed,
		"User":              user,
		"AuthenticatedUser": authenticatedUser,
		"NavBarActive":      navbarActive,
	}, "layouts/app")
}

func UserDetailFavoritePage(c *fiber.Ctx) error {

	var user model.User
	var authenticatedUser model.User
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
			return c.Redirect("/")
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

	return c.Render("users/show", fiber.Map{
		"PageTitle":         user.Name + " — Conduit",
		"FiberCtx":          c,
		"IsSelf":            isSelf,
		"IsFollowed":        isFollowed,
		"User":              user,
		"AuthenticatedUser": authenticatedUser,
		"NavBarActive":      navbarActive,
		"IsLoadFavorites":   true,
	}, "layouts/app")
}
