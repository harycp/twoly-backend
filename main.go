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

	r.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

			requestHeaders := c.GetHeader("Access-Control-Request-Headers")
			if requestHeaders != "" {
				c.Writer.Header().Set("Access-Control-Allow-Headers", requestHeaders)
			} else {
				c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
			}

			c.Writer.Header().Set("Access-Control-Max-Age", "86400")
			c.Writer.Header().Set("Vary", "Origin")
		}

		log.Println("[CORS HIT]", c.Request.Method, c.Request.URL.Path, "Origin:", origin)

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

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