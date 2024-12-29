package db

import (
	"time"

	"github.com/sebastian-nunez/golang-search-engine/types"
	"gorm.io/gorm"
)

type SearchIndex struct {
	ID        string         `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Value     string         `json:"value"`
	URLs      []CrawledURL   `gorm:"many2many:token_urls" json:"urls"` // Create virtual join-table named "token_urls"
	CreatedAt *time.Time     `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

func (s *SearchIndex) TableName() string {
	return "search_index"
}

func (s *SearchIndex) Save(index types.Index, crawledURLs []CrawledURL) error {
	for token, urlIDs := range index {
		newIndex := &SearchIndex{Value: token}

		// Indexes must not be overwritten in the database. We must either fetch the existing index, or create a new one.
		if err := DBConn.Where(SearchIndex{Value: token}).FirstOrCreate(newIndex).Error; err != nil {
			return err
		}

		var urlsToAppend []CrawledURL
		for _, urlID := range urlIDs {
			for _, url := range crawledURLs {
				if url.ID == urlID {
					urlsToAppend = append(urlsToAppend, url)
					break
				}
			}
		}

		if err := DBConn.Model(&newIndex).Association("urls").Append(&urlsToAppend); err != nil {
			return err
		}
	}

	return nil
}
