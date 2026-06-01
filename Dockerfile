# Tahap 1: Builder
FROM golang:1.22-alpine AS builder

# Set working directory di dalam container
WORKDIR /app

# Install git dan sertifikat SSL
RUN apk update && apk add --no-cache git ca-certificates tzdata

# Copy file module terlebih dahulu agar cache layer Docker optimal
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy SELURUH sisa source code dari root direktori Anda ke dalam /app
COPY . .

# Build aplikasi Go menjadi file binary statis bernama 'twoly-backend'
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o twoly-backend main.go

# Tahap 2: Runner (Image final yang sangat ringan)
FROM alpine:latest

# Set zona waktu dan sertifikat SSL
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Jakarta

WORKDIR /root/

# Copy binary 'twoly-backend' dari Tahap 1
COPY --from=builder /app/twoly-backend .

# Beritahu Docker bahwa kita menggunakan port 8080
EXPOSE 8080

# Jalankan aplikasi
CMD ["./twoly-backend"]