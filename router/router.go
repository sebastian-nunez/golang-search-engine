package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sebastian-nunez/golang-search-engine/database"
	"github.com/sebastian-nunez/golang-search-engine/handler"
	"github.com/sebastian-nunez/golang-search-engine/middleware"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	app.Get("/", middleware.WithAuth, handler.RenderHomePage)
	app.Get("/login", handler.RenderLoginPage)
	app.Post("/logout", handler.PostLogout)

	v1.Get("/ping", handler.GetPing)
	v1.Post("/login", handler.PostAdminLogin)
	v1.Post("/settings", middleware.WithAuth, handler.PostSettings)

	// DEBUG ONLY
	// TODO: For testing, always create a set of basic credentials. Remove for a production app.
	v1.Get("/create-admin", func(c *fiber.Ctx) error {
		u := &database.User{}
		err := u.CreateAdmin("jdoe@google.com", "password")
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "User created!",
		})
	})
}
