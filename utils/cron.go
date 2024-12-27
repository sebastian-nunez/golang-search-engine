package utils

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/robfig/cron/v3"
)

func StartCronJobs() {
	c := cron.New()
	c.AddFunc("@hourly", runEngine)
	c.Start()

	cronCount := len(c.Entries())
	log.Infof("Set up '%d' cron jobs", cronCount)
}

func runEngine() {
	// TODO: implement engine start
	log.Info("Starting the engine...")
}
