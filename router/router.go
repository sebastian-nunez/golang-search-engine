package router

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
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
	v1.Post("/search", handler.PostSearch)
	v1.Use("/search", cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("noCache") == "true"
		},
		Expiration:   30 * time.Minute,
		CacheControl: true, // Enable client side caching
	}))
}
