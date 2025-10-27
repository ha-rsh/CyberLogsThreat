package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	ElasticsearchURL string
	APIKey           string
}

func Load() *Config {
	godotenv.Load()
	
	cfg := &Config{
		Port:             getEnv("THREAT_SERVICE_PORT", "8082"),
		ElasticsearchURL: getEnv("ELASTICSEARCH_URL", "http://localhost:9200"),
		APIKey:           getEnv("API_KEY", "secret-api-key-12345"),
	}
	
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}