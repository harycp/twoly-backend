# Tahap 1: Builder
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk update && apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o twoly-server main.go


# Tahap 2: Runner
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Jakarta
ENV PORT=7860

WORKDIR /app

COPY --from=builder /app/twoly-server .

EXPOSE 7860

CMD ["./twoly-server"]