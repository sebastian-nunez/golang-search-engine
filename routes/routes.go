package routes

import (
	"github.com/gofiber/fiber/v2"
)

func Register(router fiber.Router) {
	v1 := router.Group("/v1")

	v1.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello world!",
		})
	})
}
