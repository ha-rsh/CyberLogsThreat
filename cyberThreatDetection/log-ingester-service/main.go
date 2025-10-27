package main

import (
	"cybersecuritySystem/log-ingester-service/config"
	"cybersecuritySystem/log-ingester-service/handlers"
	"cybersecuritySystem/log-ingester-service/repository"
	"cybersecuritySystem/log-ingester-service/service"
	"cybersecuritySystem/shared/logger"
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "cybersecuritySystem/log-ingester-service/docs"
)

// @title Log Ingestor Service API
// @version 1.0
// @description Log collection and storage service for threat detection system
// @host localhost:8081
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
func main() {
	logger.Init("log-service", logger.INFO)
	cfg := config.Load()
	logger.Info("Configuration loaded")

	repo, err := repository.NewLogRepository(cfg.ElasticsearchURL)
	if err != nil {
		logger.Fatal("Failed to create repository: %v", err)
	}
	logger.Success("Connected to Elasticsearch")

	svc := service.NewLogService(repo)
	handler := handlers.NewLogHandler(svc)

	r := mux.NewRouter()

	r.HandleFunc("/api/logs/upload", handler.UploadCSV).Methods("POST")
	r.HandleFunc("/api/logs", handler.CreateLog).Methods("POST")
	r.HandleFunc("/api/logs", handler.GetLogs).Methods("GET")
	r.HandleFunc("/api/logs/search", handler.SearchLogs).Methods("GET")
	r.HandleFunc("/api/logs/{logId}", handler.GetLogByID).Methods("GET")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	logger.Info("Log Ingestor Service running on port %s", cfg.Port)
	logger.Info("Swagger docs: http://localhost:%s/swagger/index.html", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		logger.Fatal("Server failed: %v", err)
	}
}