package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/harycp/twoly-backend/internal/config"
)

func main() {
	// 1. Load file .env
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  Warning: File .env tidak ditemukan, menggunakan environment system")
	}

	// 2. Koneksi ke Database
	config.ConnectDB()

	// 3. Setup Gin Router
	r := gin.Default()

	// Route dasar untuk cek health status backend
	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Twoly API is running perfectly 🚀",
		})
	})

	// 4. Jalankan Server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server berjalan di port: %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Gagal menjalankan server: %v", err)
	}
}