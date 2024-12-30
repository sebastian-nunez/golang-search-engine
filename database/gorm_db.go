package database

import (
	"fmt"

	"github.com/sebastian-nunez/golang-search-engine/config"
	"github.com/sebastian-nunez/golang-search-engine/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGormDB() (*gorm.DB, error) {
	dbURL := config.Envs.DatabaseURL
	var err error

	gdb, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init database connection: %s", err)
	}

	err = gdb.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	if err != nil {
		return nil, fmt.Errorf("unable to install UUID extension in the database: %s", err)

	}

	err = gdb.AutoMigrate(&model.User{}, &model.CrawlerSettings{}, &model.CrawledPage{}, &model.SearchIndex{})
	if err != nil {
		return nil, fmt.Errorf("unable to create database migrations: %s", err)
	}

	return gdb, nil
}
