package handler

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/sebastian-nunez/golang-search-engine/model"
	"github.com/sebastian-nunez/golang-search-engine/utils"
	"github.com/sebastian-nunez/golang-search-engine/views"
	"gorm.io/gorm"
)

func GetPing(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "pong",
	})
}

func PostAdminLogin(c *fiber.Ctx, gdb *gorm.DB) error {
	payload := loginPayload{}
	if err := c.BodyParser(&payload); err != nil {
		log.Info(err)
		c.Status(fiber.StatusInternalServerError)
		return htmlError(c, "unable to parse login credentials")
	}

	if len(payload.Email) == 0 || len(payload.Password) == 0 {
		log.Info("invalid login credentials provided for admin login")
		c.Status(fiber.StatusUnauthorized)
		return htmlError(c, "invalid login credentials")
	}

	user := &model.User{}
	user, err := user.LoginAsAdmin(gdb, payload.Email, payload.Password)
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
	c.Append("HX-Redirect", "/dashboard")
	return c.SendStatus(fiber.StatusOK)
}

func PostLogout(c *fiber.Ctx) error {
	c.ClearCookie(utils.AdminCookie)
	c.Append("HX-Redirect", "/login")
	return c.SendStatus(fiber.StatusOK)
}

func PostSettings(c *fiber.Ctx, gdb *gorm.DB) error {
	payload := settingsPayload{}
	if err := c.BodyParser(&payload); err != nil {
		log.Info(err)
		return htmlError(c, "unable to parse search settings")
	}

	settings := &model.CrawlerSettings{
		URLsPerHour: uint(payload.URLsPerHour),
		SearchOn:    payload.SearchOn,
		AddNewURLs:  payload.AddNewURLs,
	}
	err := settings.Update(gdb)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return htmlError(c, "unable to update search settings")
	}

	// c.Append("HX-Refresh", "true")
	return c.SendString("Settings were saved.")
}

func PostSearch(c *fiber.Ctx, gdb *gorm.DB) error {
	// TODO: will probably be a good idea to add some sort of JSON validator library.
	if len(c.Body()) == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Missing required JSON field 'query' in the request body.",
		})
	}

	payload := searchPayload{}
	if err := c.BodyParser(&payload); err != nil {
		log.Info(err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": fmt.Sprintf("Unable to parse search query '%s'. Error: %s", payload.Query, err),
		})
	}

	if payload.Query == "" {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "JSON field 'query' can not be empty",
		})
	}

	si := &model.SearchIndex{}
	pages, err := si.FullTextSearch(gdb, payload.Query)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": fmt.Sprintf("Unable to run the full text search for search query '%s'. Error: %s", payload.Query, err),
		})
	}

	return c.JSON(fiber.Map{
		"results": pages,
		"total":   len(pages),
	})
}

func RenderDashboardPage(c *fiber.Ctx, gdb *gorm.DB) error {
	settings := &model.CrawlerSettings{}
	err := settings.Get(gdb)
	if err != nil {
		err := settings.CreateDefault(gdb)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return htmlError(c, "unable to create default search settings")
		}
	}

	urlsPerHour := strconv.FormatUint(uint64(settings.URLsPerHour), 10)
	return render(c, views.Dashboard(urlsPerHour, settings.SearchOn, settings.AddNewURLs))
}

func RenderLoginPage(c *fiber.Ctx) error {
	return render(c, views.Login())
}

func RenderSearchPage(c *fiber.Ctx) error {
	return render(c, views.Search())
}
