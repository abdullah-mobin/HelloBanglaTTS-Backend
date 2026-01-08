package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadENV() error {
	// Ignore error if .env doesn't exist (production scenario)
	err := godotenv.Load()
	if err != nil {
		log.Println("ℹ️  No .env file found (using environment variables)")
	} else {
		log.Println("✅ Loaded .env file")
	}
	return nil
}
