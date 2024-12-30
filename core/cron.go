package core

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

func StartCrawlerCronJobs(gdb *gorm.DB) {
	c := cron.New()

	c.AddFunc("0 * * * *", func() { RunCrawler(gdb) })  // Run every hour
	c.AddFunc("15 * * * *", func() { RunIndexer(gdb) }) // Run every hour at 15 minutes past
	c.Start()

	cronCount := len(c.Entries())
	log.Infof("Successfully set up %d cron jobs", cronCount)
}
