package authentication

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/sqlite3"
)

var StoredAuthenticationSession *session.Store

func SessionStart() {

	store := sqlite3.New(sqlite3.Config{
		Table: "fiber_storage",
	})

	authSession := session.New(session.Config{
		Storage: store,
	})

	StoredAuthenticationSession = authSession
}

func AuthStore(c *fiber.Ctx, userID uint) {
	session, err := StoredAuthenticationSession.Get(c)
	if err != nil {
		panic(err)
	}

	session.Set("authentication", userID)
	if err := session.Save(); err != nil {
		panic(err)
	}
}

func AuthGet(c *fiber.Ctx) (bool, uint) {
	session, err := StoredAuthenticationSession.Get(c)
	if err != nil {
		panic(err)
	}

	value := session.Get("authentication")
	if value == nil {
		return false, 0
	}

	return true, value.(uint)
}

func AuthDestroy(c *fiber.Ctx) {
	session, err := StoredAuthenticationSession.Get(c)
	if err != nil {
		panic(err)
	}

	session.Delete("authentication")
	if err := session.Save(); err != nil {
		panic(err)
	}
}
