package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/harycp/twoly-backend/internal/config"
	"github.com/harycp/twoly-backend/internal/handlers"
	"github.com/harycp/twoly-backend/internal/middleware"
	"github.com/harycp/twoly-backend/internal/repositories"
	"github.com/harycp/twoly-backend/internal/services"
)

func SetupRoutes(r *gin.Engine) {
	db := config.GetDB()

	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "success", "message": "Twoly API is running perfectly 🚀"})
		})

		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		protected := v1.Group("/")
		protected.Use(middleware.RequireAuth())
		{
			protected.GET("/auth/me", authHandler.GetMe)
		}
	}
}