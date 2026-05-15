# Load environment variables from .env file
include .env

# Database URL format for golang-migrate
DB_URL="postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=$(DB_SSLMODE)"

# Create a new migration file
# Usage: make migrate-create name=table_name
migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

# Run migrations (Up)
migrate-up:
	migrate -path migrations -database $(DB_URL) -verbose up

# Rollback migrations (Down 1 step)
migrate-down:
	migrate -path migrations -database $(DB_URL) -verbose down 1

# Run the server using hot-reload (Air)
run:
	air

# Run the server manually (Fallback)
run-manual:
	go run main.go