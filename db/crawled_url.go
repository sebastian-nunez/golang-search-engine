package db

import (
	"time"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type CrawledURL struct {
	ID              string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	URL             string         `gorm:"unique;not null" json:"url"`
	Success         bool           `gorm:"not null" json:"success"`
	CrawlDuration   time.Duration  `json:"crawlDuration"`
	ResponseCode    int            `gorm:"type:smallint" json:"responseCode"`
	PageTitle       string         `json:"pageTitle"`
	PageDescription string         `json:"pageDescription"`
	Headings        string         `json:"headings"`
	Indexed         bool           `gorm:"default:false" json:"indexed"`
	LastTested      *time.Time     `json:"lastTested"` // Use pointer so this value can be nil
	CreatedAt       *time.Time     `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

func (url *CrawledURL) UpdateURL(input CrawledURL) error {
	tx := DBConn.Select("url", "success", "crawl_duration", "response_code", "page_title", "page_description", "headings", "last_tested", "updated_at").Omit("created_at").Save(&input)
	if tx.Error != nil {
		log.Info(tx.Error)
		return tx.Error
	}

	return nil
}

// GetNextCrawlURLs returns all URLs which have NOT been previously been tested.
func (url *CrawledURL) GetNextCrawlURLs(limit int) ([]CrawledURL, error) {
	var urls []CrawledURL
	tx := DBConn.Where("last_tested IS NULL").Limit(limit).Find(&urls)
	if tx.Error != nil {
		log.Info(tx.Error)
		return nil, tx.Error
	}

	return urls, nil
}

func (url *CrawledURL) Save() error {
	tx := DBConn.Save(&url)
	if tx.Error != nil {
		log.Info(tx.Error)
		return tx.Error
	}

	return nil
}

func (url *CrawledURL) GetNotIndexed() ([]CrawledURL, error) {
	var urls []CrawledURL
	tx := DBConn.Where("indexed = ? AND last_tested IS NOT NULL", false).Find(&urls)
	if tx.Error != nil {
		log.Info(tx.Error)
		return nil, tx.Error
	}

	return urls, nil
}

func (url *CrawledURL) SetIndexedTrue(urls []CrawledURL) error {
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
