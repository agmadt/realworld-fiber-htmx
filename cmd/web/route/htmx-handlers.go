package webroute

import (
	HTMXController "realworld-fiber-htmx/cmd/web/controller/htmx"

	"github.com/gofiber/fiber/v2"
)

func HTMXHandlers(app *fiber.App) {

	/* Sign In */
	app.Get("/htmx/sign-in", HTMXController.SignInPage)
	app.Post("/htmx/sign-in", HTMXController.SignInAction)
	app.Post("/htmx/sign-out", HTMXController.SignOut)

	/* Sign Up */
	app.Get("/htmx/sign-up", HTMXController.SignUpPage)
	app.Post("/htmx/sign-up", HTMXController.SignUpAction)

	/* Home */
	app.Get("/htmx/home", HTMXController.HomePage)
	app.Get("/htmx/home/your-feed", HTMXController.HomeYourFeed)
	app.Get("/htmx/home/global-feed", HTMXController.HomeGlobalFeed)
	app.Get("/htmx/home/tag-feed/:tag", HTMXController.HomeTagFeed)
	app.Get("/htmx/home/tag-list", HTMXController.HomeTagList)
	app.Post("/htmx/home/articles/:slug/favorite", HTMXController.HomeFavoriteAction)

	/* Article */
	app.Get("/htmx/articles/:slug", HTMXController.ArticleDetailPage)
	app.Get("/htmx/articles/:slug/comments", HTMXController.ArticleDetailCommentList)
	app.Post("/htmx/articles/:slug/comments", HTMXController.ArticleComment)
	app.Post("/htmx/articles/:slug/favorite", HTMXController.ArticleFavoriteAction)
	app.Post("/htmx/articles/follow-user/:slug", HTMXController.ArticleFollowAction)

	/* Editor */
	app.Get("/htmx/editor/:slug?", HTMXController.EditorPage)
	app.Post("/htmx/editor", HTMXController.StoreArticle)
	app.Patch("/htmx/editor/:slug?", HTMXController.UpdateArticle)

	/* User */
	app.Get("/htmx/users/:username", HTMXController.UserDetailPage)
	app.Get("/htmx/users/:username/articles", HTMXController.UserArticles)
	app.Get("/htmx/users/:username/favorites", HTMXController.UserArticlesFavorite)
	app.Post("/htmx/users/articles/:slug/favorite", HTMXController.UserArticleFavoriteAction)
	app.Post("/htmx/users/:username/follow", HTMXController.UserFollowAction)

	/* Setting */
	app.Get("/htmx/settings", HTMXController.SettingPage)
	app.Post("/htmx/settings", HTMXController.SettingAction)
}
