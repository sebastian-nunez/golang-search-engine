package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/sebastian-nunez/golang-search-engine/config"
	"github.com/sebastian-nunez/golang-search-engine/core"
	"github.com/sebastian-nunez/golang-search-engine/database"
	"github.com/sebastian-nunez/golang-search-engine/router"
)

func main() {
	app := fiber.New(fiber.Config{
		IdleTimeout: 10 * time.Second,
		ReadTimeout: 5 * time.Second,
	})

	app.Use(logger.New())
	app.Use(compress.New())
	app.Use(healthcheck.New(healthcheck.Config{
		LivenessEndpoint:  "/live",
		ReadinessEndpoint: "/ready",
	}))
	app.Get("/metrics", monitor.New())
	app.Use(limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Second,
	}))
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Hour,
	}))

	gdb, err := database.NewGormDB()
	if err != nil {
		panic(err)
	}
	router.SetupRoutes(gdb, app)
	core.StartCrawlerCronJobs(gdb)

	go func() {
		err := app.Listen(":" + config.Envs.Port)
		if err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
	<-sigch

	log.Info("Gracefully shutting down the server...")
	if err := app.Shutdown(); err != nil {
		log.Errorf("Error shutting down server: %v", err)
		os.Exit(1)
	}

	log.Info("Server shut down successfully.")
	os.Exit(0)
}
