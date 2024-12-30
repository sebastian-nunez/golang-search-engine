package handler

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/sebastian-nunez/golang-search-engine/db"
	"github.com/sebastian-nunez/golang-search-engine/utils"
	"github.com/sebastian-nunez/golang-search-engine/views"
)

func GetPing(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "pong",
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

	user := &db.User{}
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

func PostLogout(c *fiber.Ctx) error {
	c.ClearCookie(utils.AdminCookie)
	c.Append("HX-Redirect", "/login")
	return c.SendStatus(fiber.StatusOK)
}

func PostSettings(c *fiber.Ctx) error {
	input := settingsForm{}
	if err := c.BodyParser(&input); err != nil {
		log.Info(err)
		return htmlError(c, "unable to parse search settings")
	}

	settings := &db.SearchSettings{
		URLsPerHour: uint(input.URLsPerHour),
		SearchOn:    input.SearchOn,
		AddNewURLs:  input.AddNewURLs,
	}
	err := settings.Update()
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return htmlError(c, "unable to update search settings")
	}

	// c.Append("HX-Refresh", "true")
	return c.SendString("Settings were saved.")
}

func PostSearch(c *fiber.Ctx) error {
	// TODO: will probably be a good idea to add some sort of JSON validator library.
	if len(c.Body()) == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"success": false,
			"message": "Missing required JSON field 'term' in the request body.",
			"data":    []db.CrawledPage{},
		})
	}

	input := &searchInput{}
	if err := c.BodyParser(&input); err != nil {
		log.Info(err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": fmt.Sprintf("Unable to parse search input '%s'. Error: %s", input.Term, err),
		})
	}

	if input.Term == "" {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "JSON field 'term' can not be empty",
		})
	}

	si := &db.SearchIndex{}
	pages, err := si.FullTextSearch(input.Term)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": fmt.Sprintf("Unable to run the full text search for search term '%s'. Error: %s", input.Term, err),
		})
	}

	return c.JSON(fiber.Map{
		"results": pages,
		"total":   len(pages),
	})
}

func RenderHomePage(c *fiber.Ctx) error {
	settings := &db.SearchSettings{}
	err := settings.Get()
	if err != nil {
		err := settings.CreateDefault()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return htmlError(c, "unable to create default search settings")
		}
	}

	urlsPerHour := strconv.FormatUint(uint64(settings.URLsPerHour), 10)
	return render(c, views.Home(urlsPerHour, settings.SearchOn, settings.AddNewURLs))
}

func RenderLoginPage(c *fiber.Ctx) error {
	return render(c, views.Login())
}
