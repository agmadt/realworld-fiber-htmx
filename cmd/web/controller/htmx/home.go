package HTMXController

import (
	"math"
	"realworld-fiber-htmx/cmd/web/model"
	"realworld-fiber-htmx/internal/authentication"
	"realworld-fiber-htmx/internal/database"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func HomePage(c *fiber.Ctx) error {

	var authenticatedUser model.User

	isAuthenticated, userID := authentication.AuthGet(c)

	if isAuthenticated {
		db := database.Get()
		db.Model(&authenticatedUser).
			Where("id = ?", userID).
			First(&authenticatedUser)
	}

	return c.Render("home/htmx-home-page", fiber.Map{
		"PageTitle":         "Home",
		"NavBarActive":      "home",
		"FiberCtx":          c,
		"AuthenticatedUser": authenticatedUser,
	}, "layouts/app-htmx")
}

func HomeYourFeed(c *fiber.Ctx) error {
	var (
		articles        []model.Article
		hasArticles     bool
		user            model.User
		followings      []model.Follow
		hasPagination   bool
		totalPagination int
		count           int64
	)

	page := 0
	if c.QueryInt("page") > 1 {
		page = c.QueryInt("page") - 1
	}

	isAuthenticated, userID := authentication.AuthGet(c)
	if !isAuthenticated {
		return c.Redirect("/")
	}

	db := database.Get()
	db.Model(&user).Where("id = ?", userID).First(&user)

	db.Model(&user).Preload("Followings").Association("Followings").Find(&followings)
	if len(followings) == 0 {
		hasArticles = false
	}

	ids := make([]uint, len(followings))
	for i, f := range followings {
		ids[i] = f.FollowerID
	}

	db.Where("user_id in (?)", ids).
		Preload("Favorites").
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Order("tags.name asc")
		}).
		Preload("User").
		Limit(5).
		Offset(page * 5).
		Order("created_at desc").
		Find(&articles)

	db.Model(&articles).Where("user_id in (?)", ids).Count(&count)

	if count > 0 && (count/5 > 0) {
		pageDivision := float64(count) / float64(5)
		totalPagination = int(math.Ceil(pageDivision))
		hasPagination = true
	}

	feedNavbarItems := []fiber.Map{
		{
			"Title":     "Your Feed",
			"IsActive":  true,
			"HXPushURL": "/your-feed",
			"HXGetURL":  "/htmx/home/your-feed",
		},
		{
			"Title":     "Global Feed",
			"IsActive":  false,
			"HXPushURL": "/",
			"HXGetURL":  "/htmx/home/global-feed",
		},
	}

	if len(articles) > 0 {
		hasArticles = true
	}

	c.Render("home/htmx-home-feed", fiber.Map{
		"HasArticles":        hasArticles,
		"Articles":           articles,
		"FeedNavbarItems":    feedNavbarItems,
		"Personal":           isAuthenticated,
		"TotalPagination":    totalPagination,
		"HasPagination":      hasPagination,
		"CurrentPagination":  page + 1,
		"PushPathPagination": "your-feed",
		"PathPagination":     "your-feed",
	}, "layouts/app-htmx")

	return nil
}

func HomeGlobalFeed(c *fiber.Ctx) error {

	var (
		articles        []model.Article
		hasArticles     bool
		hasPagination   bool
		totalPagination int
		count           int64
	)

	page := 0
	if c.QueryInt("page") > 1 {
		page = c.QueryInt("page") - 1
	}

	isAuthenticated, userID := authentication.AuthGet(c)

	db := database.Get()
	db.Model(&articles).
		Preload("Favorites").
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Order("tags.name asc")
		}).
		Preload("User").
		Limit(5).
		Offset(page * 5).
		Order("created_at desc").
		Find(&articles)

	db.Model(&articles).Count(&count)

	feedNavbarItems := []fiber.Map{
		{
			"Title":     "Global Feed",
			"IsActive":  true,
			"HXPushURL": "/",
			"HXGetURL":  "/htmx/home/global-feed",
		},
	}

	if count > 0 && (count/5 > 0) {
		pageDivision := float64(count) / float64(5)
		totalPagination = int(math.Ceil(pageDivision))
		hasPagination = true
	}

	if isAuthenticated {

		feedNavbarItems = append([]fiber.Map{
			{
				"Title":     "Your Feed",
				"IsActive":  false,
				"HXPushURL": "/your-feed",
				"HXGetURL":  "/htmx/home/your-feed",
			},
		}, feedNavbarItems...)
	}

	if len(articles) > 0 {
		hasArticles = true

		for i := 0; i < len(articles); i++ {
			articles[i].IsFavorited = articles[i].FavoritedBy(userID)
		}
	}

	c.Render("home/htmx-home-feed", fiber.Map{
		"HasArticles":         hasArticles,
		"Articles":            articles,
		"FeedNavbarItems":     feedNavbarItems,
		"AuthenticatedUserID": userID,
		"TotalPagination":     totalPagination,
		"HasPagination":       hasPagination,
		"CurrentPagination":   page + 1,
		"PathPagination":      "global-feed",
	}, "layouts/app-htmx")

	return nil
}

func HomeTagFeed(c *fiber.Ctx) error {

	var (
		tag             model.Tag
		articles        []model.Article
		hasArticles     bool
		hasPagination   bool
		totalPagination int
		count           int64
	)

	page := 0
	if c.QueryInt("page") > 1 {
		page = c.QueryInt("page") - 1
	}

	isAuthenticated, _ := authentication.AuthGet(c)

	tagText := c.Params("tag")

	db := database.Get()

	db.Where(&model.Tag{Name: tagText}).First(&tag)

	db.Model(&tag).
		Preload("Favorites").
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Order("tags.name asc")
		}).
		Preload("User").
		Limit(5).
		Offset(page * 5).
		Order("created_at desc").
		Association("Articles").
		Find(&articles)

	count = db.Model(&tag).
		Association("Articles").
		Count()

	if len(articles) > 0 {
		hasArticles = true
	}

	if count > 0 && (count/5 > 0) {
		pageDivision := float64(count) / float64(5)
		totalPagination = int(math.Ceil(pageDivision))
		hasPagination = true
	}

	feedNavbarItems := []fiber.Map{
		{
			"Title":     "Global Feed",
			"IsActive":  false,
			"HXPushURL": "/",
			"HXGetURL":  "/htmx/home/global-feed",
		},
	}

	if isAuthenticated {

		feedNavbarItems = append([]fiber.Map{
			{
				"Title":     "Your Feed",
				"IsActive":  false,
				"HXPushURL": "/your-feed",
				"HXGetURL":  "/htmx/home/your-feed",
			},
		}, feedNavbarItems...)
	}

	feedNavbarItems = append(feedNavbarItems,
		fiber.Map{
			"Title":     tagText,
			"IsActive":  true,
			"HXPushURL": "/",
			"HXGetURL":  "/htmx/home/global-feed",
		},
	)

	c.Render("home/htmx-home-feed", fiber.Map{
		"HasArticles":        hasArticles,
		"Articles":           articles,
		"FeedNavbarItems":    feedNavbarItems,
		"TotalPagination":    totalPagination,
		"HasPagination":      hasPagination,
		"CurrentPagination":  page + 1,
		"PushPathPagination": "tag-feed/" + tag.Name,
		"PathPagination":     "tag-feed/" + tag.Name,
	}, "layouts/app-htmx")

	return nil
}

func HomeTagList(c *fiber.Ctx) error {

	var (
		tag     model.Tag
		tags    []model.Tag
		hasTags bool
	)

	db := database.Get()
	db.Model(&tag).
		Select("*, COUNT(id) as favorite_count").
		Preload("Articles").
		Limit(5).
		Order("favorite_count DESC").
		Group("id").
		Find(&tags)

	if len(tags) > 0 {
		hasTags = true
	}

	c.Render("home/partials/tag-item-list", fiber.Map{
		"Tags":    tags,
		"HasTags": hasTags,
	}, "layouts/app-htmx")

	return nil
}
