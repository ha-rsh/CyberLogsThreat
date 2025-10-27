package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	LogServiceURL    string
	ThreatServiceURL string
	APIKey           string
	JWTSecret        string
}

func Load() *Config {
	godotenv.Load()
	
	cfg := &Config{
		Port:             getEnv("API_GATEWAY_PORT", "8080"),
		LogServiceURL:    getEnv("LOG_SERVICE_URL", "http://localhost:8081"),
		ThreatServiceURL: getEnv("THREAT_SERVICE_URL", "http://localhost:8082"),
		APIKey:           getEnv("API_KEY", "secret"),
		JWTSecret:        getEnv("JWT_SECRET", "secret"),
	}
	
	log.Printf("API Gateway Configuration loaded")
	log.Printf("LogServiceURL: %s", cfg.LogServiceURL)
	log.Printf("ThreatServiceURL: %s", cfg.ThreatServiceURL)
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}