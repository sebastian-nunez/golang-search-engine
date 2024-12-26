package handler

import (
	"fmt"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

func render(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	return component.Render(c.Context(), c.Response().BodyWriter())
}

func htmlError(c *fiber.Ctx, message string) error {
	if len(message) == 0 {
		message = "an unspecified error occurred."
	}
	return c.SendString(fmt.Sprintf("<h2><strong>Error:</strong> %s</h2>", message))
}
