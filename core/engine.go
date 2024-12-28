package core

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/sebastian-nunez/golang-search-engine/db"
)

func RunEngine() {
	log.Info("Started search engine crawl...")
	defer log.Info("Search engine crawl has finished.")

	settings := &db.SearchSettings{}
	err := settings.Get()
	if err != nil {
		log.Errorf("Unable to get search settings: %s", err)
		return
	}

	if !settings.SearchOn {
		log.Info("Search is disabled in the settings")
		return
	}
}

func RunIndex() {
	log.Info("Started search engine indexing...")
	defer log.Info("Search engine indexing has finished.")

	// TODO: implement
}
