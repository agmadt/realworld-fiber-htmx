package HTMXController

import (
	"encoding/json"
	"errors"
	"realworld-fiber-htmx/cmd/web/model"
	"realworld-fiber-htmx/internal"
	"realworld-fiber-htmx/internal/authentication"
	"realworld-fiber-htmx/internal/database"
	"realworld-fiber-htmx/internal/helper"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

func EditorPage(c *fiber.Ctx) error {

	var authenticatedUser model.User
	var article model.Article
	hasArticle := false
	NavBarActive := "editor"

	isAuthenticated, userID := authentication.AuthGet(c)

	db := database.Get()

	if isAuthenticated {
		db.Model(&authenticatedUser).
			Where("id = ?", userID).
			First(&authenticatedUser)
	}

	if c.Params("slug") != "" {

		err := db.Model(&article).
			Where("slug = ?", c.Params("slug")).
			Preload("Tags", func(db *gorm.DB) *gorm.DB {
				return db.Order("tags.name asc")
			}).
			Find(&article).Error

		if err == nil {
			hasArticle = true
			NavBarActive = "none"
		}
	}

	return c.Render("editor/htmx-editor-page", fiber.Map{
		"PageTitle":         "Editor",
		"FiberCtx":          c,
		"NavBarActive":      NavBarActive,
		"AuthenticatedUser": authenticatedUser,
		"HasArticle":        hasArticle,
		"Article":           article,
	}, "layouts/app-htmx")
}

func StoreArticle(c *fiber.Ctx) error {

	type TagItem struct {
		Value string
	}

	var (
		errorBag          []string
		tagItems          []TagItem
		authenticatedUser model.User
	)
	validate := internal.NewValidator()

	isAuthenticated, userID := authentication.AuthGet(c)
	if !isAuthenticated {
		return c.Redirect("/")
	}

	db := database.Get()

	db.Model(&authenticatedUser).
		Where("id = ?", userID).
		First(&authenticatedUser)

	article := &model.Article{
		Title:       c.FormValue("title"),
		Slug:        slug.Make(c.FormValue("title")),
		Description: c.FormValue("description"),
		Body:        c.FormValue("content"),
		UserID:      userID,
	}

	err := validate.Validate(article)
	if err != nil {

		for _, err := range err.(validator.ValidationErrors) {
			errorBag = append(errorBag, internal.ErrorMessage(err.Field(), err.Tag()))
		}

		return c.Render("editor/htmx-editor-page", fiber.Map{
			"IsOob":             true,
			"FiberCtx":          c,
			"NavBarActive":      "editor",
			"Errors":            errorBag,
			"AuthenticatedUser": authenticatedUser,
		}, "layouts/app-htmx")
	}

	db.Create(article)

	if c.FormValue("tags") != "" {
		json.Unmarshal([]byte(c.FormValue("tags")), &tagItems)

		for i := 0; i < len(tagItems); i++ {
			tagItem := tagItems[i]
			tag := model.Tag{Name: tagItem.Value}

			err := db.Model(&tag).Where("name = ?", tagItem.Value).First(&tag).Error
			if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
				db.Create(&tag)
			}

			if err := db.Model(&article).Association("Tags").Append(&tag); err != nil {
				return err
			}
		}
	}

	return helper.HTMXRedirectTo("/articles/"+article.Slug, "/htmx/articles/"+article.Slug, c)
}

func UpdateArticle(c *fiber.Ctx) error {

	type TagItem struct {
		Value string
	}

	var (
		errorBag          []string
		tagItems          []TagItem
		authenticatedUser model.User
		article           model.Article
		tags              []model.Tag
	)

	validate := internal.NewValidator()

	isAuthenticated, userID := authentication.AuthGet(c)
	if !isAuthenticated {
		return c.Redirect("/")
	}

	db := database.Get()

	db.Model(&authenticatedUser).
		Where("id = ?", userID).
		First(&authenticatedUser)

	err := db.Model(&article).
		Where("slug = ?", c.Params("slug")).
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Order("tags.name asc")
		}).
		Find(&article).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Redirect("/")
		}
	}

	article.Title = c.FormValue("title")
	article.Description = c.FormValue("description")
	article.Body = c.FormValue("content")

	err = validate.Validate(article)
	if err != nil {

		for _, err := range err.(validator.ValidationErrors) {
			errorBag = append(errorBag, internal.ErrorMessage(err.Field(), err.Tag()))
		}

		return c.Render("editor/htmx-editor-page", fiber.Map{
			"IsOob":             true,
			"FiberCtx":          c,
			"NavBarActive":      "editor",
			"Errors":            errorBag,
			"AuthenticatedUser": authenticatedUser,
			"HasArticle":        true,
			"Article":           article,
		}, "layouts/app-htmx")
	}

	article.Slug = slug.Make(c.FormValue("title"))

	db.Updates(article)

	if c.FormValue("tags") != "" {
		json.Unmarshal([]byte(c.FormValue("tags")), &tagItems)

		for i := 0; i < len(tagItems); i++ {
			tagItem := tagItems[i]
			tag := model.Tag{Name: tagItem.Value}

			err := db.Model(&tag).Where("name = ?", tagItem.Value).First(&tag).Error
			if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
				db.Create(&tag)
			}

			tags = append(tags, tag)
		}

		if err := db.Model(&article).Association("Tags").Replace(&tags); err != nil {
			return err
		}
	}

	return helper.HTMXRedirectTo("/articles/"+article.Slug, "/htmx/articles/"+article.Slug, c)
}
