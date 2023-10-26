package webroute

import (
	"realworld-fiber-htmx/cmd/web/controller"

	"github.com/gofiber/fiber/v2"
)

type PageData struct {
	PageTitle string
}

func WebHandlers(app *fiber.App) {

	/* Sign In */
	app.Get("/sign-in", controller.SignInPage)

	/* Sign Up */
	app.Get("/sign-up", controller.SignUpPage)

	/* Home */
	app.Get("/", controller.HomePage)
	app.Get("/your-feed", controller.YourFeedPage)
	app.Get("/tag-feed/:slug", controller.TagFeedPage)

	/* Article */
	app.Get("/articles/:slug", controller.ArticleDetailPage)

	/* Editor */
	app.Get("/editor/:slug?", controller.EditorPage)

	/* User */
	app.Get("/users/:username", controller.UserDetailPage)
	app.Get("/users/:username/articles", controller.UserDetailPage)
	app.Get("/users/:username/favorites", controller.UserDetailFavoritePage)

	/* Setting */
	app.Get("/settings", controller.SettingPage)
}
