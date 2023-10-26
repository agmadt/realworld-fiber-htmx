package renderer

import (
	"errors"
	"realworld-fiber-htmx/internal/authentication"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func ViewEngineStart() *html.Engine {

	viewEngine := html.New("./cmd/web/templates", ".tmpl")

	viewEngine.AddFunc("IsAuthenticated", func(c *fiber.Ctx) bool {
		isAuthenticated, _ := authentication.AuthGet(c)
		return isAuthenticated
	})

	viewEngine.AddFunc("Iterate", func(start int, end int) []int {
		n := end - start + 1
		result := make([]int, n)
		for i := 0; i < n; i++ {
			result[i] = start + i
		}
		return result
	})

	viewEngine.AddFunc("Dict", func(values ...interface{}) (map[string]interface{}, error) {
		if len(values)%2 != 0 {
			return nil, errors.New("invalid dict call")
		}
		dict := make(map[string]interface{}, len(values)/2)
		for i := 0; i < len(values); i += 2 {
			key, ok := values[i].(string)
			if !ok {
				return nil, errors.New("dict keys must be strings")
			}
			dict[key] = values[i+1]
		}
		return dict, nil
	})

	return viewEngine
}
