.PHONY: all build clean test lint proto migrate-up migrate-down dev docker-up docker-down

# Generate go.sum
go-sum:
	cd backend && go mod tidy

# Build all backend services
build: go-sum
	cd backend && go build ./cmd/...

# Run all tests
test: go-sum
	cd backend && go test -v -race -coverprofile=coverage.out ./...

# Lint
lint:
	cd backend && golangci-lint run ./...

# Clean build artifacts
clean:
	rm -rf backend/build dist

# Database migrations
migrate-up:
	cd backend && migrate -path migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" up

migrate-down:
	cd backend && migrate -path migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" down

# Run dev server
dev:
	cd backend && go run ./cmd/api-gateway

# Docker
docker-up:
	docker compose -f infrastructure/docker/docker-compose.yml up -d

docker-down:
	docker compose -f infrastructure/docker/docker-compose.yml down

# Proto generation
proto:
	cd backend && buf generate proto

# Tools installation
tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/bufbuild/buf/cmd/buf@latest
