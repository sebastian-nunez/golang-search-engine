package handler

import (
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
		return htmlError(c, "unable to parse login credentials")
	}

	return c.SendStatus(200)
}

func PostSettings(c *fiber.Ctx) error {
	settings := settingsForm{}
	if err := c.BodyParser(&settings); err != nil {
		log.Info(err)
		return htmlError(c, "unable to parse search settings")
	}

	return c.SendStatus(200)
}

func RenderHomePage(c *fiber.Ctx) error {
	return render(c, views.Home())
}

func RenderLoginPage(c *fiber.Ctx) error {
	return render(c, views.Login())
}
