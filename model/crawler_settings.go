package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type CrawlerSettings struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	URLsPerHour uint      `json:"urlsPerHour"`
	SearchOn    bool      `json:"searchOn"`
	AddNewURLs  bool      `json:"addNewUrls"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (cs *CrawlerSettings) Get(gdb *gorm.DB) error {
	// We will only store 1 row for the settings
	if err := gdb.Where("id = ?", 1).First(cs).Error; err != nil {
		return err
	}

	return nil
}

func (cs *CrawlerSettings) Update(gdb *gorm.DB) error {
	tx := gdb.Select("urls_per_hour", "search_on", "add_new_urls", "updated_at").Where("id = ?", 1).Updates(cs)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (cs *CrawlerSettings) CreateDefault(gdb *gorm.DB) error {
	settings := CrawlerSettings{
		URLsPerHour: 10,
		SearchOn:    true,
		AddNewURLs:  true,
	}

	if err := gdb.Create(&settings).Error; err != nil {
		return fmt.Errorf("unable to create the search settings in the database: %s", err)
	}

	return nil
}
