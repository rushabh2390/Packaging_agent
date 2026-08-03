.PHONY: dev-backend dev-frontend docs docker-up docker-down

# Run Go Fiber backend locally
dev-backend:
	cd backend && swag init -g cmd/api/main.go && go run cmd/api/main.go

# Run Nuxt 3 frontend locally
dev-frontend:
	cd frontend && npm run dev

# Re-generate Swagger docs for Go backend
docs:
	cd backend && swag init -g cmd/api/main.go

# Run whole stack with Docker
docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down