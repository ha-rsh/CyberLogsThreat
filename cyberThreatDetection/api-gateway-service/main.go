package main

import (
	"cybersecuritySystem/api-gateway-service/config"
	"cybersecuritySystem/api-gateway-service/handlers"
	"cybersecuritySystem/api-gateway-service/middleware"
	"cybersecuritySystem/api-gateway-service/proxy"
	"cybersecuritySystem/shared/auth"
	"cybersecuritySystem/shared/logger"
	"log"
	"net/http"
	"time"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/gorilla/mux"
	_ "cybersecuritySystem/api-gateway-service/docs"
)


func main() {
	logger.Init("api-gateway", logger.INFO)
	cfg := config.Load()

	jwtSecret := cfg.JWTSecret
	auth.SetJWTSecret(jwtSecret)

	reverseProxy := proxy.NewReverseProxy(cfg.LogServiceURL, cfg.ThreatServiceURL)
	healthHandler := handlers.NewHealthHandler()
	authHandler := handlers.NewAuthHandler()

	r := mux.NewRouter()

	r.Use(middleware.CORS)
	r.Use(middleware.Logging)
	r.Use(middleware.RateLimit(200, time.Minute))

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	r.HandleFunc("/health", healthHandler.Health).Methods("GET")

	r.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST")

	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.JWTMiddleware)

	api.HandleFunc("/auth/refresh", authHandler.RefreshToken).Methods("POST")

	api.PathPrefix("/").HandlerFunc(reverseProxy.Route)

	logger.Info("API Gateway running on port %s", cfg.Port)
	logger.Info("Swagger UI available at: http://localhost:%s/swagger/index.html", cfg.Port)
	logger.Info("Proxying /api/logs/* to %s", cfg.LogServiceURL)
	logger.Info("Proxying /api/threats/* to %s", cfg.ThreatServiceURL)

	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal(err)
	}
}

