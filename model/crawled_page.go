package model

import (
	"time"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type CrawledPage struct {
	// Internal metadata
	ID      string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	URL     string `gorm:"unique;not null" json:"url"`
	Indexed bool   `gorm:"default:false" json:"indexed"` // This is NOT reference to a database index. Instead, it signals if the page has been indexed for full-text search.

	// Metadata from the crawling process
	Success       bool          `gorm:"default:false;not null" json:"success"`
	CrawlDuration time.Duration `json:"crawlDuration"`
	StatusCode    int           `gorm:"type:smallint" json:"statusCode"` // HTTP status code
	LastTested    *time.Time    `json:"lastTested"`                      // Use pointer so this value can be nil

	// Information extracted from the page
	Title       string `json:"title"`       // From the <title> meta tag
	Description string `json:"description"` // From the "description" <meta> tag
	Headings    string `json:"headings"`    // Combined string of all <h1> tags separated by ", "

	// Timestamps for tracking changes
	CreatedAt *time.Time     `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"` // Soft delete date
}

func (cp *CrawledPage) Update(gdb *gorm.DB, input CrawledPage) error {
	tx := gdb.Select("url", "success", "crawl_duration", "status_code", "title", "description", "headings", "last_tested", "updated_at").Omit("created_at").Save(&input)
	if tx.Error != nil {
		log.Info(tx.Error)
		return tx.Error
	}

	return nil
}

// GetNextCrawlPages returns all pages which have NOT been previously been tested.
func (cp *CrawledPage) GetNextCrawlPages(gdb *gorm.DB, limit int) ([]CrawledPage, error) {
	var pages []CrawledPage
	tx := gdb.Where("last_tested IS NULL").Limit(limit).Find(&pages)
	if tx.Error != nil {
		log.Info(tx.Error)
		return nil, tx.Error
	}

	return pages, nil
}

func (cp *CrawledPage) Save(gdb *gorm.DB) error {
	tx := gdb.Save(&cp)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (cp *CrawledPage) GetNotIndexed(gdb *gorm.DB) ([]CrawledPage, error) {
	var pages []CrawledPage
	tx := gdb.Where("indexed = false AND success = true AND last_tested IS NOT NULL").Find(&pages)
	if tx.Error != nil {
		log.Info(tx.Error)
		return nil, tx.Error
	}

	return pages, nil
}

func (cp *CrawledPage) SetIndexedTrue(gdb *gorm.DB, pages []CrawledPage) error {
	for _, p := range pages {
		p.Indexed = true

		tx := gdb.Save(&p)
		if tx.Error != nil {
			log.Info(tx.Error)
			return tx.Error
		}
	}

	return nil
}
