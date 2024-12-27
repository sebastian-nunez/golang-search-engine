package db

import (
	"fmt"

	"github.com/gofiber/fiber/v2/log"
	"github.com/sebastian-nunez/golang-search-engine/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBConn *gorm.DB

func InitDB() {
	dbURL := config.Envs.DatabaseURL
	var err error

	DBConn, err = gorm.Open(postgres.Open(dbURL))
	if err != nil {
		msg := fmt.Sprintf("Failed to init database connection: %s", err)
		log.Info(msg)
		panic(msg)
	}

	err = DBConn.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	if err != nil {
		msg := fmt.Sprintf("Unable to install UUID extension in the database: %s", err)
		log.Info(msg)
		panic(err)
	}

	err = DBConn.AutoMigrate(&User{}, &SearchSettings{}, &CrawledURL{})
	if err != nil {
		msg := fmt.Sprintf("Unable to create database migrations: %s", err)
		log.Info(msg)
		panic(err)
	}
}

func GetDB() *gorm.DB {
	return DBConn
}
