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
		log.Infof("Unable to get next crawl URLs from the database: %s", err)
		return
	}

	newURLs := make(map[string]struct{})
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
			newURLs[newURL] = struct{}{}
		}
	} // End of range

	if !settings.AddNewURLs {
		log.Info("Adding new urls to database is disabled")
		return
	}

	added := 0
	for u := range newURLs {
		newURL := db.CrawledURL{URL: u}

		err := newURL.Save()
		if err != nil {
			log.Infof("Unable to save new URL '%s' to the database", newURL.URL)
		} else {
			added += 1
		}
	}

	// TODO: make this into a table.
	log.Infof("Crawled through %d URLs. Found a total of %d new URLs. Added %d new URLs to the database", len(nextURLs), len(newURLs), added)
}

func RunIndex() {
	log.Info("Started search engine indexing...")
	defer log.Info("Search engine indexing has finished.")

	// TODO: implement
}
