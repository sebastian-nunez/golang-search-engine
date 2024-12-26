package routes

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/sebastian-nunez/golang-search-engine/views"
)

func Register(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	app.Get("/", func(c *fiber.Ctx) error {
		return render(c, views.Home())
	})

	v1.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, world!",
		})
	})
}

func render(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	return component.Render(c.Context(), c.Response().BodyWriter())
}
