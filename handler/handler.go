package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/sebastian-nunez/golang-search-engine/database"
	"github.com/sebastian-nunez/golang-search-engine/utils"
	"github.com/sebastian-nunez/golang-search-engine/views"
)

func GetHelloWorld(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hello, world!",
	})
}

func PostAdminLogin(c *fiber.Ctx) error {
	input := loginForm{}
	if err := c.BodyParser(&input); err != nil {
		log.Info(err)
		c.Status(fiber.StatusInternalServerError)
		return htmlError(c, "unable to parse login credentials")
	}

	if len(input.Email) == 0 || len(input.Password) == 0 {
		log.Info("invalid login credentials provided for admin login")
		c.Status(fiber.StatusUnauthorized)
		return htmlError(c, "invalid login credentials")
	}

	user := &database.User{}
	user, err := user.LoginAsAdmin(input.Email, input.Password)
	if err != nil {
		log.Info(err)
		c.Status(fiber.StatusUnauthorized)
		return htmlError(c, "invalid login credentials")
	}

	signedToken, err := utils.CreateNewAuthToken(user.ID, user.Email, user.IsAdmin)
	if err != nil {
		log.Info(err)
		c.Status(fiber.StatusInternalServerError)
		return htmlError(c, "unable to sign JWT token")
	}

	cookie := fiber.Cookie{
		Name:     utils.AdminCookie,
		Value:    signedToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)
	c.Append("HX-Redirect", "/")
	return c.SendStatus(fiber.StatusOK)
}

func GetLogout(c *fiber.Ctx) error {
	c.ClearCookie(utils.AdminCookie)
	c.Append("HX-Redirect", "/login")
	return c.SendStatus(fiber.StatusOK)
}

func PostSettings(c *fiber.Ctx) error {
	settings := settingsForm{}
	if err := c.BodyParser(&settings); err != nil {
		log.Info(err)
		return htmlError(c, "unable to parse search settings")
	}

	return c.SendStatus(fiber.StatusOK)
}

func RenderHomePage(c *fiber.Ctx) error {
	return render(c, views.Home())
}

func RenderLoginPage(c *fiber.Ctx) error {
	return render(c, views.Login())
}
