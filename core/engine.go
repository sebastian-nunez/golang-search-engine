package core

import (
	"math"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/sebastian-nunez/golang-search-engine/model"
	"gorm.io/gorm"
)

// RunCrawler starts and manages the search engine crawling process.
// It begins crawling based on the search settings, performs crawls on retrieved URLs,
// updates existing pages and potentially adds newly discovered EXTERNAL URLs to the database for future crawling.
func RunCrawler(gdb *gorm.DB) {
	log.Info("Started search engine crawl...")
	defer log.Info("Search engine crawl has finished.")

	settings := &model.CrawlerSettings{}
	err := settings.Get(gdb)
	if err != nil {
		log.Errorf("Unable to get search settings: %s", err)
		return
	}

	if !settings.SearchOn {
		log.Info("Search is disabled in the settings")
		return
	}

	cp := &model.CrawledPage{}
	nextPages, err := cp.GetNextCrawlPages(gdb, int(settings.URLsPerHour))
	if err != nil {
		log.Infof("Unable to get next crawl pages from the database: %s", err)
		return
	}

	newURLs := make(map[string]struct{})
	lastTested := time.Now()
	numberWidth := int(math.Log10(float64(len(nextPages)))) + 1
	for i, page := range nextPages {
		log.Infof("Crawling: %3d%% (%0*d/%d) %s", int(float64(i+1)/float64(len(nextPages))*100), numberWidth, i+1, len(nextPages), page.URL)
		result := CrawlPage(page.URL)

		if !result.Success {
			err := page.Update(gdb, model.CrawledPage{
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

		err := page.Update(gdb, model.CrawledPage{
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
		newPage := model.CrawledPage{URL: url}

		err := newPage.Save(gdb)
		if err != nil {
			log.Infof("Unable to save new URL '%s' to the database", newPage.URL)
		} else {
			added += 1
		}
	}

	// TODO: make this into a table.
	log.Infof("Crawled through %d pages. Found a total of %d new URLs. Added %d new pages to explore into the database.", len(nextPages), len(newURLs), added)
}

// RunIndexer performs the process of building and saving the search engine index.
func RunIndexer(gdb *gorm.DB) {
	log.Info("Started search engine indexing...")
	defer log.Info("Search engine indexing has finished.")

	cp := &model.CrawledPage{}
	notIndexed, err := cp.GetNotIndexed(gdb)
	if err != nil {
		log.Info("Unable to get un-indexed pages from the database: %s", err)
		return
	}
	log.Infof("Found %d pages to be indexed.", len(notIndexed))

	// Add un-indexed pages to the current index in-memory.
	idx := make(InvertedIndex)
	idx.Add(notIndexed) // Large number of pages can cause memory issues here

	si := &model.SearchIndex{}
	err = si.Save(gdb, idx, notIndexed)
	if err != nil {
		log.Infof("Unable to save the new search index into the database: %s", err)
		return
	}

	err = cp.SetIndexedTrue(gdb, notIndexed)
	if err != nil {
		log.Infof("Unable to mark the newly indexed pages as indexed in the database: %s", err)
		return
	}
}
