package db

import (
	"fmt"
	"time"
)

type SearchSettings struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	URLsPerHour uint      `json:"urlsPerHour"`
	SearchOn    bool      `json:"searchOn"`
	AddNewURLs  bool      `json:"addNewUrls"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (s *SearchSettings) Get() error {
	// We will only store 1 row for the settings
	if err := DBConn.Where("id = ?", 1).First(s).Error; err != nil {
		return err
	}

	return nil
}

func (s *SearchSettings) Update() error {
	tx := DBConn.Select("urls_per_hour", "search_on", "add_new_urls", "updated_at").Where("id = ?", 1).Updates(s)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (s *SearchSettings) CreateDefault() error {
	settings := SearchSettings{
		URLsPerHour: 10,
		SearchOn:    true,
		AddNewURLs:  true,
	}

	if err := DBConn.Create(&settings).Error; err != nil {
		return fmt.Errorf("unable to create the search settings in the database: %s", err)
	}

	*s = settings
	return nil
}
