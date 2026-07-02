# Twoly Backend

Twoly Backend is the REST API service for **Twoly**, a private digital space for couples to save memories, plan dates, manage shared moments, send love notes, and support realtime interaction.

This repository contains the backend service built with **Go**, **Gin**, **GORM**, and **PostgreSQL/Supabase**.

---

## Tech Stack

| Purpose          | Technology                       |
| ---------------- | -------------------------------- |
| Language         | Go                               |
| HTTP Framework   | Gin                              |
| ORM              | GORM                             |
| Database         | PostgreSQL / Supabase PostgreSQL |
| Migration        | golang-migrate                   |
| Authentication   | JWT                              |
| Password Hashing | bcrypt                           |
| Image Storage    | Cloudinary                       |
| Realtime         | Supabase Realtime                |
| Hot Reload       | Air                              |

---

## Main Features

- User authentication
- JWT protected routes
- Couple pairing system
- Memory timeline
- Shared album
- Date planner
- Shared calendar
- Love notes
- Realtime presence and touch
- Cloudinary media upload through backend with automatic compression

---

## Project Structure

```txt
twoly-backend/
├── internal/
│   ├── config/
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   ├── repositories/
│   ├── routes/
│   ├── services/
│   └── utils/
├── migrations/
├── .air.toml
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── main.go
└── README.md
```

---

## Getting Started

### 1. Clone Repository

```bash
git clone https://github.com/harycp/twoly-backend.git
cd twoly-backend
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Setup Environment Variables

Create a `.env` file in the root directory.

```env
APP_ENV=development
APP_PORT=8080

DB_HOST=your-db-host
DB_PORT=6543
DB_DATABASE=postgres
DB_USERNAME=your-db-username
DB_PASSWORD=your-db-password
DB_SSLMODE=require

DATABASE_URL=postgres://your-db-username:your-db-password@your-db-host:6543/postgres?sslmode=require&prefer_simple_protocol=true

JWT_SECRET=your-jwt-secret
JWT_EXPIRES_IN=24h

SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your-supabase-anon-key
SUPABASE_SERVICE_ROLE_KEY=your-supabase-service-role-key

CLOUDINARY_CLOUD_NAME=your-cloudinary-cloud-name
CLOUDINARY_API_KEY=your-cloudinary-api-key
CLOUDINARY_API_SECRET=your-cloudinary-api-secret
CLOUDINARY_FOLDER=twoly/memories
CLOUDINARY_AVATAR_FOLDER=twoly/avatar
MAX_MEDIA_UPLOAD_SIZE_MB=100

MAX_IMAGE_UPLOAD_SIZE_MB=5
```

> Do not commit `.env` to the repository.

---

## Run Development Server

### Run manually

```bash
go run main.go
```

### Run with hot reload

Install Air:

```bash
go install github.com/air-verse/air@latest
```

Run:

```bash
air
```

Server will run on:

```txt
http://localhost:8080
```

---

## Database Migration

This project uses `golang-migrate`.

### Install migrate

```bash
go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Create migration

```bash
migrate create -ext sql -dir migrations create_table_name
```

### Run migration

```bash
migrate -path migrations -database "postgres://your-db-username:your-db-password@your-db-host:6543/postgres?sslmode=require" up
```

### Rollback migration

```bash
migrate -path migrations -database "postgres://your-db-username:your-db-password@your-db-host:6543/postgres?sslmode=require" down 1
```

---

## Supabase Pooler Note

When using Supabase Pooler with GORM, add this parameter to the database URL:

```txt
prefer_simple_protocol=true
```

Example:

```env
DATABASE_URL=postgres://user:password@host:6543/postgres?sslmode=require&prefer_simple_protocol=true
```

This helps avoid prepared statement conflicts when using connection pooling.

---

## API Overview

Base URL:

```txt
http://localhost:8080/api/v1
```

### Public Routes

```txt
GET  /health
POST /auth/register
POST /auth/login
```

### Protected Routes

```txt
GET /auth/me
```

More endpoints will be added as the MVP modules are developed.

---

## Planned Modules

- Auth
- Couple System
- Memories
- Albums
- Date Plans
- Calendar Events
- Love Notes
- Realtime Touch
- Cloudinary Upload

---

## Security Notes

- Store secrets only in environment variables.
- Never commit `.env`.
- Hash passwords using bcrypt.
- Protect private routes with JWT middleware.
- Validate ownership before accessing couple resources.
- Upload images through backend only.
- Keep Cloudinary API Secret on the backend.

---

## Development Status

Twoly Backend is currently in early MVP development.

Current focus:

- Backend foundation
- Authentication
- Database migration
- Couple system
- Memory timeline API

---

## Author

Developed by **Hary Capri**

```txt
https://github.com/harycp
```
