package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sebastian-nunez/golang-search-engine/handler"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	app.Get("/", handler.RenderHomePage)
	app.Get("/login", handler.RenderLoginPage)

	v1.Get("/", handler.GetHelloWorld)
}
