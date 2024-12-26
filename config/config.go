package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
}

// Envs is the global configuration for the application.
var Envs = initConfig()

func initConfig() Config {
	err := godotenv.Load()
	if err != nil {
		panic("Unable to locate environment variables")
	}

	return Config{
		Port: getEnv("PORT", "8080"),
	}
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
