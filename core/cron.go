package core

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/robfig/cron/v3"
)

func StartCrawlerCronJobs() {
	c := cron.New()

	c.AddFunc("0 * * * *", RunEngine) // Run every hour
	c.AddFunc("15 * * * *", RunIndex) // Run every hour at 15 minutes past
	c.Start()

	cronCount := len(c.Entries())
	log.Infof("Finished setting up %d cron jobs", cronCount)
}
