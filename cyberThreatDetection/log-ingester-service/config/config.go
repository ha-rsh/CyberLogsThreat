package config

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	Port 					string
	ElasticsearchURL		string
	APIKey					string
}

func Load() *Config {
	godotenv.Load()

	cfg := &Config{
		Port: getEnv("LOG_SERVICE_PORT", "8081"),
		ElasticsearchURL: getEnv("ELASTICSEARCH_URL", "http://localhost:9200"),
		APIKey: getEnv("API_KEY", "secret-api-key-12345"),
	}
	log.Printf("APIkey: %s", cfg.APIKey)
	log.Printf("Log Service configuration loaded")
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}