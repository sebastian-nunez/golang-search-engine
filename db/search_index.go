package db

import (
	"time"

	"gorm.io/gorm"
)

// Index is an in-memory inverted index. It maps tokens to Page IDs (stored in the database).
type Index map[string][]string

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

func (si *SearchIndex) Save(index Index, pages []CrawledPage) error {
	for token, pageIDs := range index {
		newIndex := &SearchIndex{Value: token}

		// Indexes must not be overwritten in the database. We must either fetch the existing index, or create a new one.
		if err := DBConn.Where(SearchIndex{Value: token}).FirstOrCreate(newIndex).Error; err != nil {
			return err
		}

		var pagesToAppend []CrawledPage
		for _, pageID := range pageIDs {
			for _, page := range pages {
				if page.ID == pageID {
					pagesToAppend = append(pagesToAppend, page)
					break
				}
			}
		}

		if err := DBConn.Model(&newIndex).Association("pages").Append(&pagesToAppend); err != nil {
			return err
		}
	}

	return nil
}
