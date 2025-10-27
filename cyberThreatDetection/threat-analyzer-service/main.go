package main

import (
	"cybersecuritySystem/shared/logger"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	"cybersecuritySystem/threat-analyzer-service/client"
	"cybersecuritySystem/threat-analyzer-service/config"
	_ "cybersecuritySystem/threat-analyzer-service/docs"
	"cybersecuritySystem/threat-analyzer-service/handlers"
	"cybersecuritySystem/threat-analyzer-service/service"
)

// @title Threat Analyzer Service API
// @version 1.0
// @description Threat detection and analysis service
// @host localhost:8082
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key

func main() {
	// Initialize logger
	logger.Init("threat-service", logger.INFO)
	
	cfg := config.Load()
	logger.Info("Configuration loaded")

	esClient, err := client.NewElasticsearchClient(cfg.ElasticsearchURL)
	if err != nil {
		logger.Fatal("Failed to create Elasticsearch client: %v", err)
	}
	logger.Success("Connected to Elasticsearch")

	svc := service.NewThreatService(esClient)
	handler := handlers.NewThreatHandler(svc)

	r := mux.NewRouter()
	
	r.Use(loggingMiddleware)
	// r.Use(corsMiddleware)
	// r.Use(authMiddleware(cfg.APIKey))

	r.HandleFunc("/api/threats/analyze", handler.AnalyzeThreats).Methods("POST")
	r.HandleFunc("/api/threats", handler.GetThreats).Methods("GET")
	r.HandleFunc("/api/threats/search", handler.SearchThreats).Methods("GET")
	r.HandleFunc("/api/threats/{threatId}", handler.GetThreatByID).Methods("GET")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	logger.Info("Threat Analyzer Service running on port %s", cfg.Port)
	logger.Info("Swagger docs: http://localhost:%s/swagger/index.html", cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		logger.Fatal("Server failed: %v", err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Skip logging for swagger
		if strings.HasPrefix(r.URL.Path, "/swagger") {
			next.ServeHTTP(w, r)
			return
		}
		
		next.ServeHTTP(w, r)
		
		duration := time.Since(start)
		logger.HTTP(r.Method, r.URL.Path, 200, duration)
	})
}
