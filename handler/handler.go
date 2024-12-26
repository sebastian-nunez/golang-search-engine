package handler

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/sebastian-nunez/golang-search-engine/views"
)

func GetHelloWorld(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hello, world!",
	})
}

func RenderHomePage(c *fiber.Ctx) error {
	return render(c, views.Home())
}

func render(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	return component.Render(c.Context(), c.Response().BodyWriter())
}
