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

	// Repositories
	userRepo := repositories.NewUserRepository(db)
	coupleRepo := repositories.NewCoupleRepository(db)
	memoryRepo := repositories.NewMemoryRepository(db)
	photoRepo := repositories.NewMemoryPhotoRepository(db)

	// Services
	cloudinarySvc := services.NewCloudinaryService()
	authService := services.NewAuthService(userRepo)
	coupleService := services.NewCoupleService(coupleRepo)
	memoryService := services.NewMemoryService(memoryRepo, coupleRepo)
	photoService := services.NewMemoryPhotoService(photoRepo, memoryRepo, coupleRepo, cloudinarySvc)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	coupleHandler := handlers.NewCoupleHandler(coupleService)
	memoryHandler := handlers.NewMemoryHandler(memoryService)
	photoHandler := handlers.NewMemoryPhotoHandler(photoService)

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
			
			// Couple Routes
			couples := protected.Group("/couples")
			{
				couples.POST("/invite", coupleHandler.CreateInvite)
				couples.POST("/join", coupleHandler.JoinCouple)
				couples.GET("/me", coupleHandler.GetMyCouple)
			}

			// Memory Routes
			memories := protected.Group("/memories")
			{
				memories.POST("/", memoryHandler.CreateMemory)
				memories.GET("/", memoryHandler.GetAllMemories)
				memories.GET("/:id", memoryHandler.GetMemoryDetail)
				memories.PUT("/:id", memoryHandler.UpdateMemory)
				memories.DELETE("/:id", memoryHandler.DeleteMemory)

				// Memory Photos Routes
				memories.POST("/:id/photos", photoHandler.UploadPhotos)
				memories.GET("/:id/photos", photoHandler.GetPhotos)
				memories.DELETE("/:id/photos/:photoId", photoHandler.DeletePhoto)
			}
		}
	}
}