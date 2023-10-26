package HTMXController

import (
	"errors"
	"realworld-fiber-htmx/cmd/web/model"
	"realworld-fiber-htmx/internal"
	"realworld-fiber-htmx/internal/authentication"
	"realworld-fiber-htmx/internal/database"
	"realworld-fiber-htmx/internal/helper"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ArticleDetailCommentList(c *fiber.Ctx) error {

	var article model.Article

	isAuthenticated, _ := authentication.AuthGet(c)

	db := database.Get()

	db.Model(&article).
		Where("slug = ?", c.Params("slug")).
		Preload("User").
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Order("id DESC").Preload("User")
		}).
		Find(&article)

	return c.Render("articles/partials/comments-wrapper", fiber.Map{
		"Article":         article,
		"IsAuthenticated": isAuthenticated,
	}, "layouts/app-htmx")
}

func ArticleComment(c *fiber.Ctx) error {

	var (
		errorBag          []string
		article           model.Article
		comment           model.Comment
		authenticatedUser model.User
	)
	validate := internal.NewValidator()

	isAuthenticated, userID := authentication.AuthGet(c)
	if !isAuthenticated {
		return helper.HTMXRedirectTo("/sign-in", "/htmx/sign-in", c)
	}

	db := database.Get()

	db.Model(&authenticatedUser).
		Where("id = ?", userID).
		First(&authenticatedUser)

	err := db.Model(&article).
		Where("slug = ?", c.Params("slug")).
		Preload("User").
		Find(&article).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.HTMXRedirectTo("/", "/htmx/home", c)
		}
	}

	comment.UserID = userID
	comment.ArticleID = article.ID
	comment.Body = c.FormValue("comment")

	err = validate.Validate(comment)
	if err != nil {

		for _, err := range err.(validator.ValidationErrors) {
			errorBag = append(errorBag, internal.ErrorMessage(err.Field(), err.Tag()))
		}

		return c.Render("components/error-message", fiber.Map{
			"Errors": errorBag,
		}, "layouts/app-htmx")
	}

	comment.User = authenticatedUser

	db.Create(&comment)

	return c.Render("articles/htmx-post-comments", fiber.Map{
		"IsOob":   true,
		"Article": article,
		"Comment": comment,
		"User":    authenticatedUser,
	}, "layouts/app-htmx")
}
