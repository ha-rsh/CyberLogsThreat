.PHONY: help up down restart logs build clean install status

help:
	@echo "Available commands:"
	@echo "  make up        - Start all services (Elasticsearch, API Gateway, Log Service, Threat Service, Frontend)"
	@echo "  make down      - Stop all services"
	@echo "  make restart   - Restart all services"
	@echo "  make logs      - View all logs"
	@echo "  make build     - Rebuild and start all services"
	@echo "  make clean     - Remove all containers and volumes"
	@echo "  make install   - Install all dependencies"
	@echo "  make status    - Check status of all services"

up:
	docker compose up -d
	@echo "✓ All services started:"
	@echo "  Frontend:       http://localhost:3000"
	@echo "  API Gateway:    http://localhost:8080"
	@echo "  Log Service:    http://localhost:8081"
	@echo "  Threat Service: http://localhost:8082"
	@echo "  Elasticsearch:  http://localhost:9200"

down:
	docker compose down
	@echo "✓ All services stopped"

restart:
	docker compose restart
	@echo "✓ All services restarted"

logs:
	docker compose logs -f

build:
	docker compose up -d --build
	@echo "✓ All services rebuilt and started"

clean:
	docker compose down -v
	@echo "✓ All containers, volumes, and data removed"

install:
	@echo "Installing Go dependencies..."
	cd cyberThreatDetection && go mod tidy
	@echo "Installing Node.js dependencies..."
	cd cyber_threat_frontend && npm install
	@echo "✓ All dependencies installed"

status:
	@echo "Service Status:"
	@docker compose ps