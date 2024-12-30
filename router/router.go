package router

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/sebastian-nunez/golang-search-engine/handler"
	"github.com/sebastian-nunez/golang-search-engine/middleware"
	"gorm.io/gorm"
)

func SetupRoutes(gdb *gorm.DB, app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	app.Get("/", middleware.WithAuth, func(c *fiber.Ctx) error {
		return handler.RenderHomePage(c, gdb)
	})
	app.Get("/login", handler.RenderLoginPage)
	app.Post("/logout", handler.PostLogout)

	v1.Get("/ping", handler.GetPing)
	v1.Post("/login", func(c *fiber.Ctx) error {
		return handler.PostAdminLogin(c, gdb)
	})
	v1.Post("/settings", middleware.WithAuth, func(c *fiber.Ctx) error {
		return handler.PostSettings(c, gdb)
	})
	v1.Post("/search", func(c *fiber.Ctx) error {
		return handler.PostSearch(c, gdb)
	})
	v1.Use("/search", cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("noCache") == "true"
		},
		Expiration:   30 * time.Minute,
		CacheControl: true, // Enable client side caching
	}))
}
