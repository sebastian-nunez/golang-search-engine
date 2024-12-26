package handler

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/sebastian-nunez/golang-search-engine/views"
)

func GetHelloWorld(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hello, world!",
	})
}

func PostLogin(c *fiber.Ctx) error {
	input := loginForm{}
	if err := c.BodyParser(&input); err != nil {
		log.Info(err)
		return c.SendString("<h2>Error: unable to parse login credentials</h2>")
	}

	return c.SendStatus(200)
}

func RenderHomePage(c *fiber.Ctx) error {
	return render(c, views.Home())
}

func RenderLoginPage(c *fiber.Ctx) error {
	return render(c, views.Login())
}

func render(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	return component.Render(c.Context(), c.Response().BodyWriter())
}

type loginForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}
