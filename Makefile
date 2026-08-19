.PHONY: dev-backend dev-frontend docs docker-up docker-down migrate

# Run Go Fiber backend locally from backend/main.go
dev-backend:
	cd backend && swag init -g main.go && go run main.go

# Run database migrations
migrate:
	cd backend && go run cmd/migrate/migrate.go

# Run Nuxt 3 frontend locally
dev-frontend:
	cd frontend && npm run dev

# Re-generate Swagger docs for Go backend
docs:
	cd backend && swag init -g main.go

# Run whole stack with Docker
docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down