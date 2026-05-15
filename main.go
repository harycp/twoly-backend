package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/harycp/twoly-backend/internal/config"
	"github.com/harycp/twoly-backend/internal/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("[WARN] .env file not found. Using system environment variables.")
	}

	config.ConnectDB()
	config.ConnectCloudinary()

	r := gin.Default()

	routes.SetupRoutes(r)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	printServerBanner(port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("[ERROR] Failed to start server: %v", err)
	}
}

func printServerBanner(port string) {
	log.Printf(`
+--------------------------------------------------+
|              TWOLY BACKEND SYSTEM                |
+--------------------------------------------------+
| STATUS       : ONLINE                            |
| ENVIRONMENT  : DEVELOPMENT                       |
| PROTOCOL     : HTTP                              |
| HOST         : localhost                         |
| PORT         : %-32s |
| BASE URL     : http://localhost:%-15s |
+--------------------------------------------------+
| SYSTEM LOG   : SERVER BOOT SEQUENCE COMPLETED    |
+--------------------------------------------------+
`, port, port)
}