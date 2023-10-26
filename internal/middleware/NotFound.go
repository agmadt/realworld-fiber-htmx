package middleware

import "github.com/gofiber/fiber/v2"

func NotFound(c *fiber.Ctx) error {

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"message": "Page not found.",
	})
}
