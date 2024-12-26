package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/sebastian-nunez/golang-search-engine/config"
)

func main() {
	app := fiber.New(fiber.Config{
		IdleTimeout: 10 * time.Second,
	})

	app.Use(logger.New())
	app.Use(healthcheck.New(healthcheck.Config{
		LivenessEndpoint:  "/live",
		ReadinessEndpoint: "/ready",
	}))
	app.Get("/metrics", monitor.New())

	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/", func(c *fiber.Ctx) error {
		log.Info("Getting data!")
		return c.SendString("Hello world!")
	})

	app.Listen(":" + config.Envs.Port)
}
