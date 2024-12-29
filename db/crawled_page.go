package db

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

func (cp *CrawledPage) Update(input CrawledPage) error {
	tx := DBConn.Select("url", "success", "crawl_duration", "status_code", "title", "description", "headings", "last_tested", "updated_at").Omit("created_at").Save(&input)
	if tx.Error != nil {
		log.Info(tx.Error)
		return tx.Error
	}

	return nil
}

// GetNextCrawlPages returns all pages which have NOT been previously been tested.
func (cp *CrawledPage) GetNextCrawlPages(limit int) ([]CrawledPage, error) {
	var urls []CrawledPage
	tx := DBConn.Where("last_tested IS NULL").Limit(limit).Find(&urls)
	if tx.Error != nil {
		log.Info(tx.Error)
		return nil, tx.Error
	}

	return urls, nil
}

func (cp *CrawledPage) Save() error {
	tx := DBConn.Save(&cp)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (cp *CrawledPage) GetNotIndexed() ([]CrawledPage, error) {
	var urls []CrawledPage
	tx := DBConn.Where("indexed = ? AND last_tested IS NOT NULL", false).Find(&urls)
	if tx.Error != nil {
		log.Info(tx.Error)
		return nil, tx.Error
	}

	return urls, nil
}

func (cp *CrawledPage) SetIndexedTrue(urls []CrawledPage) error {
	for _, u := range urls {
		u.Indexed = true

		tx := DBConn.Save(&u)
		if tx.Error != nil {
			log.Info(tx.Error)
			return tx.Error
		}
	}

	return nil
}
