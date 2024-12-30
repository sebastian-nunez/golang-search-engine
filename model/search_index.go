package model

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type SearchIndex struct {
	ID        string         `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Value     string         `json:"value"`
	Pages     []CrawledPage  `gorm:"many2many:token_pages" json:"pages"` // Create virtual join-table named "token_pages"
	CreatedAt *time.Time     `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

func (si *SearchIndex) TableName() string {
	return "search_index"
}

// Save persists the inverted index data to the database, establishing
// associations between tokens and crawled pages.
//
// `index` is an invented index mapping tokens (words) to a set of pageIDs.
func (si *SearchIndex) Save(gdb *gorm.DB, index map[string]map[string]struct{}, pages []CrawledPage) error {
	for token, pageIDs := range index {
		newIndex := &SearchIndex{Value: token}

		// Indexes must not be overwritten in the database. We must either fetch the existing index, or create a new one.
		if err := gdb.Where(SearchIndex{Value: token}).FirstOrCreate(newIndex).Error; err != nil {
			return err
		}

		// Since the inverted index only holds the pageIDs, we have to manually find which page
		// corresponds with that ID to get the full page object for Gorm to create the association.
		var pagesToAppend []CrawledPage
		for pageID := range pageIDs {
			for _, page := range pages {
				if page.ID == pageID {
					pagesToAppend = append(pagesToAppend, page)
					break
				}
			}
		}

		if err := gdb.Model(&newIndex).Association("Pages").Append(&pagesToAppend); err != nil {
			return err
		}
	}

	return nil
}

// FullTextSearch retrieves SearchIndex records from the database that contain the given term(s),
// along with their associated Pages.
func (si *SearchIndex) FullTextSearch(gdb *gorm.DB, text string) ([]CrawledPage, error) {
	searchTerms := strings.Fields(text)
	if len(searchTerms) == 0 {
		return nil, fmt.Errorf("no search terms were found")
	}

	pages := []CrawledPage{}
	for _, term := range searchTerms {
		var searchIndexes []SearchIndex
		// Run a "contains" text query across all index tokens for the search term
		// and eagerly fetch the corresponding CrawledPages at the same time.
		if err := gdb.Preload("Pages").Where("value LIKE ?", "%"+term+"%").Find(&searchIndexes).Error; err != nil {
			return nil, err
		}

		for _, idx := range searchIndexes {
			pages = append(pages, idx.Pages...)
		}
	}

	return pages, nil
}
