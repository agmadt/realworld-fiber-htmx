package authentication

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/mysql"
)

var StoredAuthenticationSession *session.Store

func SessionStart() {

	store := mysql.New(mysql.Config{
		ConnectionURI: "root@tcp(127.0.0.1:3306)/conduit?charset=utf8mb4&parseTime=True&loc=Local",
		Table:         "fiber_storage",
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
