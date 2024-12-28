package core

import (
	"time"

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

	crawl := &db.CrawledURL{}
	nextURLs, err := crawl.GetNextCrawlURLs(int(settings.URLsPerHour))
	if err != nil {
		log.Info("Unable to get next crawl URLs with '%d' URLs per hour: %s", settings.URLsPerHour, err)
		return
	}

	newURLs := []db.CrawledURL{}
	lastTested := time.Now()
	for _, curURL := range nextURLs {
		result := RunCrawl(curURL.URL)

		if !result.Success {
			err := curURL.UpdateURL(db.CrawledURL{
				ID:              curURL.ID,
				URL:             curURL.URL,
				Success:         false,
				CrawlDuration:   result.ParsedBody.CrawlTime,
				ResponseCode:    result.ResponseCode,
				PageTitle:       result.ParsedBody.PageTitle,
				PageDescription: result.ParsedBody.PageDescription,
				Headings:        result.ParsedBody.Headings,
				LastTested:      &lastTested,
			})
			if err != nil {
				log.Infof("Unable to save CrawledURL data for URL '%s':", curURL.URL, err)
			}

			continue
		}

		err := curURL.UpdateURL(db.CrawledURL{
			ID:              curURL.ID,
			URL:             curURL.URL,
			Success:         result.Success,
			CrawlDuration:   result.ParsedBody.CrawlTime,
			ResponseCode:    result.ResponseCode,
			PageTitle:       result.ParsedBody.PageTitle,
			PageDescription: result.ParsedBody.PageDescription,
			Headings:        result.ParsedBody.Headings,
			LastTested:      &lastTested,
		})
		if err != nil {
			log.Infof("Unable to save CrawledURL data for URL '%s':", curURL.URL, err)
		}

		// Only external URLs will be added to the database. However, we could also run the internal
		// links/URLs as well: this is out of scope for now.
		for _, newURL := range result.ParsedBody.Links.External {
			newURLs = append(newURLs, db.CrawledURL{URL: newURL})
		}
	} // End of range

	if !settings.AddNewURLs {
		log.Info("Adding new urls to database is disabled")
		return
	}

	added := 0
	for _, newURL := range newURLs {
		err := newURL.Save()
		if err != nil {
			log.Info("Unable to save new URL '%s' to the database: %s", newURL, err)
		} else {
			added += 1
		}
	}

	log.Infof("Added %d new URLs to the database", added)
}

func RunIndex() {
	log.Info("Started search engine indexing...")
	defer log.Info("Search engine indexing has finished.")

	// TODO: implement
}
