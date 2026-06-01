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
	calendarRepo := repositories.NewCalendarRepository(db)
	loveNoteRepo := repositories.NewLoveNoteRepository(db)

	// Services
	cloudinarySvc := services.NewCloudinaryService()
	emailSvc := services.NewEmailService()
	authService := services.NewAuthService(userRepo, cloudinarySvc, emailSvc)
	coupleService := services.NewCoupleService(coupleRepo)
	memoryService := services.NewMemoryService(memoryRepo, coupleRepo)
	photoService := services.NewMemoryPhotoService(photoRepo, memoryRepo, coupleRepo, cloudinarySvc)
	datePlanService := services.NewDatePlanService(datePlanRepo, coupleRepo, memoryRepo)
	calendarService := services.NewCalendarService(calendarRepo, memoryRepo, datePlanRepo, coupleRepo)
	loveNoteService := services.NewLoveNoteService(loveNoteRepo, coupleRepo)
	userService := services.NewUserService(userRepo)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	coupleHandler := handlers.NewCoupleHandler(coupleService)
	memoryHandler := handlers.NewMemoryHandler(memoryService)
	photoHandler := handlers.NewMemoryPhotoHandler(photoService)
	datePlanHandler := handlers.NewDatePlanHandler(datePlanService)
	calendarHandler := handlers.NewCalendarHandler(calendarService)
	loveNoteHandler := handlers.NewLoveNoteHandler(loveNoteService)
	userHandler := handlers.NewUserHandler(userService)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "success", "message": "Twoly API is running perfectly 🚀"})
		})

		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/google", authHandler.GoogleLogin)
			auth.POST("/forgot-password", authHandler.ForgotPassword)
			auth.POST("/verify-otp", authHandler.VerifyOTP)
			auth.POST("/reset-password", authHandler.ResetPassword)
		}

		protected := v1.Group("")
		protected.Use(middleware.RequireAuth())
		{
			protected.GET("/auth/me", authHandler.GetMe)
			protected.PUT("/auth/me", authHandler.UpdateProfile)
			protected.PUT("/users/presence", userHandler.UpdatePresence)
			
			// Couple Routes
			couples := protected.Group("/couples")
			{
				couples.POST("/invite", coupleHandler.CreateInvite)
				couples.POST("/join", coupleHandler.JoinCouple)
				couples.GET("/me", coupleHandler.GetMyCouple)
				couples.PUT("/me", coupleHandler.UpdateCoupleSettings)
			}

			// Memory Routes
			memories := protected.Group("/memories")
			{
				memories.POST("", memoryHandler.CreateMemory)
				memories.GET("", memoryHandler.GetAllMemories)
				memories.GET("/photos", photoHandler.GetGalleryPhotos)
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
				datePlans.POST("", datePlanHandler.CreateDatePlan)
				datePlans.GET("", datePlanHandler.GetAllDatePlans)
				datePlans.GET("/:id", datePlanHandler.GetDatePlanDetail)
				datePlans.PUT("/:id", datePlanHandler.UpdateDatePlan)
				datePlans.DELETE("/:id", datePlanHandler.DeleteDatePlan)
				datePlans.PATCH("/:id/status", datePlanHandler.UpdateStatus)
				datePlans.PATCH("/:id/checklists/:checklistId", datePlanHandler.UpdateChecklistItem)
				datePlans.POST("/:id/checklists", datePlanHandler.AddChecklistItem)
				datePlans.DELETE("/:id/checklists/:checklistId", datePlanHandler.DeleteChecklistItem)
				datePlans.POST("/:id/convert-to-memory", datePlanHandler.ConvertToMemory)
			}

			// Calendar Routes
			calendar := protected.Group("/calendar")
			{
				calendar.GET("/events", calendarHandler.GetEvents)
				calendar.POST("/events", calendarHandler.CreateCustomEvent)
				calendar.PUT("/events/:id", calendarHandler.UpdateCustomEvent)
				calendar.DELETE("/events/:id", calendarHandler.DeleteCustomEvent)
			}

			// Love Notes Routes
			loveNotes := protected.Group("/love-notes")
			{
				loveNotes.POST("", loveNoteHandler.CreateLoveNote)
				loveNotes.GET("", loveNoteHandler.GetLoveNotes)
				loveNotes.POST("/:id/open", loveNoteHandler.OpenLoveNote)
				loveNotes.DELETE("/:id", loveNoteHandler.DeleteLoveNote)
			}
		}
	}
}