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
	datePlanRepo := repositories.NewDatePlanRepository(db)

	// Services
	cloudinarySvc := services.NewCloudinaryService()
	authService := services.NewAuthService(userRepo)
	coupleService := services.NewCoupleService(coupleRepo)
	memoryService := services.NewMemoryService(memoryRepo, coupleRepo)
	photoService := services.NewMemoryPhotoService(photoRepo, memoryRepo, coupleRepo, cloudinarySvc)
	datePlanService := services.NewDatePlanService(datePlanRepo, coupleRepo, memoryRepo)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	coupleHandler := handlers.NewCoupleHandler(coupleService)
	memoryHandler := handlers.NewMemoryHandler(memoryService)
	photoHandler := handlers.NewMemoryPhotoHandler(photoService)
	datePlanHandler := handlers.NewDatePlanHandler(datePlanService)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "success", "message": "Twoly API is running perfectly"})
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

				memories.POST("/:id/photos", photoHandler.UploadPhotos)
				memories.GET("/:id/photos", photoHandler.GetPhotos)
				memories.DELETE("/:id/photos/:photoId", photoHandler.DeletePhoto)
			}

			// Date Plan Routes
			datePlans := protected.Group("/date-plans")
			{
				datePlans.POST("/", datePlanHandler.CreateDatePlan)
				datePlans.GET("/", datePlanHandler.GetAllDatePlans)
				datePlans.GET("/:id", datePlanHandler.GetDatePlanDetail)
				datePlans.PUT("/:id", datePlanHandler.UpdateDatePlan)
				datePlans.DELETE("/:id", datePlanHandler.DeleteDatePlan)
				datePlans.PATCH("/:id/status", datePlanHandler.UpdateStatus)
				datePlans.PATCH("/:id/checklists/:checklistId", datePlanHandler.UpdateChecklistItem)
				datePlans.POST("/:id/checklists", datePlanHandler.AddChecklistItem)
				datePlans.DELETE("/:id/checklists/:checklistId", datePlanHandler.DeleteChecklistItem)
				datePlans.POST("/:id/convert-to-memory", datePlanHandler.ConvertToMemory)
			}
		}
	}
}