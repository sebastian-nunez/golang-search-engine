package core

import (
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/sebastian-nunez/golang-search-engine/db"
)

// RunEngine starts and manages the search engine crawling process.
// It retrieves crawl settings from the database, filters based on enabled
// crawling and URLs per hour limit, performs crawls on retrieved URLs,
// updates existing pages and potentially adds newly discovered external URLs
// to the database for future crawling.
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

	cp := &db.CrawledPage{}
	nextPages, err := cp.GetNextCrawlPages(int(settings.URLsPerHour))
	if err != nil {
		log.Infof("Unable to get next crawl pages from the database: %s", err)
		return
	}

	newURLs := make(map[string]struct{})
	lastTested := time.Now()
	for _, page := range nextPages {
		result := RunCrawl(page.URL)

		if !result.Success {
			err := page.Update(db.CrawledPage{
				ID:            page.ID,
				URL:           page.URL,
				Success:       false,
				CrawlDuration: result.ParsedPage.CrawlTime,
				StatusCode:    result.StatusCode,
				Title:         result.ParsedPage.Title,
				Description:   result.ParsedPage.Description,
				Headings:      result.ParsedPage.Headings,
				LastTested:    &lastTested,
			})
			if err != nil {
				log.Infof("Unable to save CrawledURL data for URL '%s':", page.URL, err)
			}

			continue
		}

		err := page.Update(db.CrawledPage{
			ID:            page.ID,
			URL:           page.URL,
			Success:       result.Success,
			CrawlDuration: result.ParsedPage.CrawlTime,
			StatusCode:    result.StatusCode,
			Title:         result.ParsedPage.Title,
			Description:   result.ParsedPage.Description,
			Headings:      result.ParsedPage.Headings,
			LastTested:    &lastTested,
		})
		if err != nil {
			log.Infof("Unable to save CrawledURL data for URL '%s':", page.URL, err)
		}

		// Only external URLs will be added to the database. However, we could also run the internal
		// links/URLs as well: this is out of scope for now.
		for _, url := range result.ParsedPage.Links.External {
			newURLs[url] = struct{}{}
		}
	} // End of range

	if !settings.AddNewURLs {
		log.Info("Adding new urls to database is disabled")
		return
	}

	added := 0
	for url := range newURLs {
		newPage := db.CrawledPage{URL: url}

		err := newPage.Save()
		if err != nil {
			log.Infof("Unable to save new URL '%s' to the database", newPage.URL)
		} else {
			added += 1
		}
	}

	// TODO: make this into a table.
	log.Infof("Crawled through %d pages. Found a total of %d new URLs. Added %d new pages to explore into the database.", len(nextPages), len(newURLs), added)
}

// RunIndex performs the process of building and saving the search engine index.
func RunIndex() {
	log.Info("Started search engine indexing...")
	defer log.Info("Search engine indexing has finished.")

	cp := &db.CrawledPage{}
	notIndexed, err := cp.GetNotIndexed()
	if err != nil {
		log.Info("Unable to get un-indexed pages from the database: %s", err)
		return
	}
	log.Infof("There are %d pages which are not indexed", len(notIndexed))

	// Add un-indexed pages to the current index in-memory.
	idx := make(InvertedIndex)
	idx.Add(notIndexed) // Large number of pages can cause memory issues here

	si := &db.SearchIndex{}
	err = si.Save(idx, notIndexed)
	if err != nil {
		log.Infof("Unable to save the new search index into the database: %s", err)
		return
	}

	err = cp.SetIndexedTrue(notIndexed)
	if err != nil {
		log.Infof("Unable to mark the newly indexed pages as indexed in the database: %s", err)
		return
	}
}
