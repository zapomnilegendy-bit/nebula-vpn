#!/bin/bash
set -euo pipefail

echo "=== NebulaVPN Development Setup ==="

# Check prerequisites
echo "Checking prerequisites..."

command -v go >/dev/null 2>&1 || { echo "Go is required. Install: brew install go"; exit 1; }
command -v docker >/dev/null 2>&1 || { echo "Docker is required. Install: brew install --cask docker"; exit 1; }
command -v rustc >/dev/null 2>&1 || { echo "Rust is recommended. Install: curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh"; }

echo "All prerequisites found!"

# Copy environment file
if [ ! -f .env ]; then
    cp .env.example .env
    echo "Created .env from .env.example"
fi

# Install Go tools
echo "Installing Go tools..."
cd backend
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
cd ..

# Start infrastructure
echo "Starting Docker infrastructure..."
docker compose -f infrastructure/docker/docker-compose.yml up -d postgres redis

echo "Waiting for services to become healthy..."
sleep 5

# Run migrations
echo "Running database migrations..."
cd backend
make migrate-up
cd ..

echo ""
echo "=== Setup Complete ==="
echo ""
echo "Start the API server:"
echo "  cd backend && go run ./cmd/api-gateway"
echo ""
echo "Server will be available at: http://localhost:8080"
echo "Health check: http://localhost:8080/health"
echo ""
echo "Stop infrastructure:"
echo "  docker compose -f infrastructure/docker/docker-compose.yml down"
