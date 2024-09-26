package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv(filePath string) {
	err := godotenv.Load(filePath)
	if err != nil {
		log.Fatalf("Error loading .env file from %s: %v", filePath, err)
	}
}
