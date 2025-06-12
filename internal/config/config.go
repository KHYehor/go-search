package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr              string
	MongoURI          string
	MongoDBName       string
	MongoDBCollection string
	OutputDir         string
}

func Load() *Config {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, falling back to system environment")
	}

	return &Config{
		Addr:              getEnv("ADDR", ":3000"),
		MongoURI:          getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDBName:       getEnv("MONGO_DB_NAME", "DB"),
		MongoDBCollection: getEnv("MONGO_DB_COLLECTION", "Collection"),
		OutputDir:         getEnv("OUTPUT_DIRECTORY", "/output/"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
