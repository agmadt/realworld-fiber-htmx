package helper

import "github.com/gofiber/fiber/v2"

func HTMXRedirectTo(HXURL string, HXGETURL string, c *fiber.Ctx) error {

	c.Append("HX-Replace-Url", HXURL)
	c.Append("HX-Reswap", "none")

	return c.Render("components/redirect", fiber.Map{
		"HXGet":     HXGETURL,
		"HXTarget":  "#app-body",
		"HXTrigger": "load",
	}, "layouts/app-htmx")
}
