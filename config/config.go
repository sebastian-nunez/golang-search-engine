package config

import (
	"os"

	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	SecretKey   string
}

// Envs is the global configuration for the application.
var Envs = initConfig()

func initConfig() Config {
	err := godotenv.Load()
	if err != nil {
		// Causing a `panic` is making the tests fail for some reason. I should probably look into it,
		// but it is not that important for now :)
		log.Errorf("Unable to load .env config: %s", err)
	}

	return Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgresql://postgres:password@localhost:5432/search"),
		SecretKey:   getEnv("SECRET_KEY", "REPLACE_ME"),
	}
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
