package HTMXController

import (
	"realworld-fiber-htmx/cmd/web/model"
	"realworld-fiber-htmx/internal"
	"realworld-fiber-htmx/internal/authentication"
	"realworld-fiber-htmx/internal/database"
	"realworld-fiber-htmx/internal/helper"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func SettingPage(c *fiber.Ctx) error {

	var authenticatedUser model.User

	isAuthenticated, userID := authentication.AuthGet(c)

	if isAuthenticated {
		db := database.Get()
		db.Model(&authenticatedUser).
			Where("id = ?", userID).
			First(&authenticatedUser)
	}

	return c.Render("settings/htmx-setting-page", fiber.Map{
		"PageTitle":         "Settings",
		"NavBarActive":      "settings",
		"FiberCtx":          c,
		"AuthenticatedUser": authenticatedUser,
	}, "layouts/app-htmx")

}

func SettingAction(c *fiber.Ctx) error {

	var errorBag []string
	validate := internal.NewValidator()

	isAuthenticated, userID := authentication.AuthGet(c)
	if !isAuthenticated {
		return helper.HTMXRedirectTo("/sign-in", "/htmx/sign-in", c)
	}

	user := &model.User{
		ID:       userID,
		Image:    c.FormValue("image"),
		Name:     c.FormValue("name"),
		Bio:      c.FormValue("bio"),
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
	}

	err := validate.Validate(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errorBag = append(errorBag, internal.ErrorMessage(err.Field(), err.Tag()))
		}

		return c.Render("settings/partials/form-message", fiber.Map{
			"IsOob":  true,
			"Errors": errorBag,
		}, "layouts/app-htmx")
	}

	if user.Password != "" {
		user.HashPassword()
	}

	db := database.Get()
	db.Model(user).Updates(user)

	return c.Render("settings/partials/htmx-form-message", fiber.Map{
		"IsOob":             true,
		"SuccessMessages":   []string{"Data successfully saved."},
		"NavBarActive":      "settings",
		"FiberCtx":          c,
		"AuthenticatedUser": user,
	}, "layouts/app-htmx")
}
