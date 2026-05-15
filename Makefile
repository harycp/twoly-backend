# Load environment variables dari file .env
include .env

# Membuat format URL Database untuk golang-migrate
DB_URL="postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=$(DB_SSLMODE)"

# Command untuk membuat file migration baru
# Penggunaan: make migrate-create name=nama_tabel
migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

# Command untuk menjalankan migration (Apply ke database)
migrate-up:
	migrate -path migrations -database $(DB_URL) -verbose up

# Command untuk me-rollback migration 1 langkah ke belakang
migrate-down:
	migrate -path migrations -database $(DB_URL) -verbose down 1

# Command untuk menjalankan server backend
run:
	go run main.go