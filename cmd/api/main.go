package main

import (
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"razum-backend/internal/config"
	"razum-backend/internal/database"
	"razum-backend/internal/handlers"
	"razum-backend/internal/middleware"
	"razum-backend/internal/repository"
	"razum-backend/internal/services"
)

func main() {
	cfg := config.Load()
	gin.SetMode(cfg.GinMode)

	if err := database.Init(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	db := database.GetDB()

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(db)
	eventRepo := repository.NewEventRepository(db)
	participationRepo := repository.NewParticipationRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	ratingRepo := repository.NewRatingRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	dashboardRepo := repository.NewDashboardRepository(db)
	cadreRepo := repository.NewCadreRepository(db)
	adminRepo := repository.NewAdminRepository(db)
	filterRepo := repository.NewFilterRepository(db)

	// Инициализация сервисов
	authService := services.NewAuthService(userRepo, cfg)
	filterService := services.NewFilterService(filterRepo)
	profileService := services.NewProfileService(userRepo)
	eventService := services.NewEventService(eventRepo, userRepo, participationRepo)
	participationService := services.NewParticipationService(
		participationRepo,
		eventRepo,
		userRepo,
		auditRepo,
		ratingRepo,
	)
	ratingService := services.NewRatingService(ratingRepo)
	reviewService := services.NewReviewService(reviewRepo, eventRepo, userRepo)
	dashboardService := services.NewDashboardService(dashboardRepo)
	cadreService := services.NewCadreService(cadreRepo)
	pdfService := services.NewPDFService(cadreService)
	adminService := services.NewAdminService(adminRepo, userRepo, db)
	searchService := services.NewSearchService(db)

	// Инициализация handlers
	authHandler := handlers.NewAuthHandler(authService)
	profileHandler := handlers.NewProfileHandler(profileService)
	eventHandler := handlers.NewEventHandler(eventService)
	participationHandler := handlers.NewParticipationHandler(participationService, eventService)
	ratingHandler := handlers.NewRatingHandler(ratingService)
	reviewHandler := handlers.NewReviewHandler(reviewService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)
	cadreHandler := handlers.NewCadreHandler(cadreService, pdfService, filterService)
	adminHandler := handlers.NewAdminHandler(adminService)
	searchHandler := handlers.NewSearchHandler(searchService)
	filterHandler := handlers.NewFilterHandler(filterService)

	// Настройка роутера
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	// Добавляем middleware для принудительной UTF-8
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		c.Next()
	})

	// Публичные эндпоинты
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Razum Backend is running",
		})
	})

	api := router.Group("/api")
	{
		// Auth endpoints (публичные)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Profile endpoints (требуют авторизации)
		profile := api.Group("/profile")
		profile.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			profile.GET("", profileHandler.GetMyProfile)
			profile.PUT("", profileHandler.UpdateMyProfile)
			profile.GET("/:id", profileHandler.GetPublicProfile)
		}

		// Events endpoints
		events := api.Group("/events")
		{
			// Публичные эндпоинты с опциональной авторизацией
			events.GET("", middleware.OptionalAuthMiddleware(cfg.JWTSecret), eventHandler.ListEvents)
			events.GET("/:id", middleware.OptionalAuthMiddleware(cfg.JWTSecret), eventHandler.GetEvent)

			// Защищенные эндпоинты (требуют авторизацию)
			protected := events.Group("")
			protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
			{
				// CRUD мероприятий (только организаторы)
				protected.POST("", middleware.RequireOrganizer(), eventHandler.CreateEvent)
				protected.PUT("/:id", middleware.RequireOrganizer(), eventHandler.UpdateEvent)
				protected.DELETE("/:id", middleware.RequireOrganizer(), eventHandler.DeleteEvent)

				// Участие в мероприятиях
				protected.POST("/:id/register", participationHandler.RegisterForEvent)
				protected.GET("/:id/my-qr", participationHandler.GetMyQRCode)
				protected.DELETE("/:id/cancel", participationHandler.CancelRegistration)

				// Подтверждение участия (только организатор мероприятия)
				protected.POST("/:id/confirm", middleware.RequireOrganizer(), participationHandler.ConfirmParticipation)

				// Список участников мероприятия (только организатор этого мероприятия)
				protected.GET("/:id/participants", middleware.RequireOrganizer(), participationHandler.GetEventParticipants)
			}
		}

		// Мои участия
		my := api.Group("/my")
		my.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			my.GET("/participations", participationHandler.GetMyParticipations)
		}

		// Rating endpoints
		rating := api.Group("/rating")
		{
			rating.GET("/global", ratingHandler.GetGlobalRating)
			rating.GET("/direction/:direction", ratingHandler.GetRatingByDirection)
			rating.GET("/user/:id", ratingHandler.GetUserRating)

			// Защищенные эндпоинты
			protectedRating := rating.Group("")
			protectedRating.Use(middleware.AuthMiddleware(cfg.JWTSecret))
			{
				protectedRating.GET("/me", ratingHandler.GetMyRating)
				protectedRating.POST("/refresh", middleware.RequireRole("admin"), ratingHandler.RefreshRatingCache)
			}
		}

		// Review endpoints (отзывы об организаторах)
		reviews := api.Group("/organizers")
		{
			// Публичные эндпоинты
			reviews.GET("/:id/reviews", reviewHandler.GetOrganizerReviews)

			// Защищенные эндпоинты
			protectedReviews := reviews.Group("")
			protectedReviews.Use(middleware.AuthMiddleware(cfg.JWTSecret))
			{
				protectedReviews.POST("/:id/reviews", middleware.RequireParticipant(), reviewHandler.CreateReview)
			}
		}

		// Dashboard endpoints
		dashboard := api.Group("/dashboard")
		dashboard.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			dashboard.GET("", dashboardHandler.GetDashboard)
			dashboard.GET("/events", dashboardHandler.GetRecentEvents)
			dashboard.GET("/rating-history", dashboardHandler.GetRatingHistory)
			dashboard.GET("/tags", dashboardHandler.GetTrendingTags)
			dashboard.GET("/stats", dashboardHandler.GetActivityStats)
		}

		// Cadre endpoints (кадровый резерв) - только для наблюдателей
		cadre := api.Group("/cadre")
		cadre.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			cadre.GET("/candidates", middleware.RequireObserver(), cadreHandler.GetCandidates)
			cadre.GET("/candidates/:id/export", middleware.RequireObserver(), cadreHandler.ExportCandidatePDF)
		}

		// Admin endpoints
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(cfg.JWTSecret)) // только аутентификация, без роли
		{
			// Управление организаторами (только admin)
			admin.GET("/organizers/pending", middleware.RequireRole("admin"), adminHandler.GetPendingOrganizers)
			admin.POST("/organizers/:id/approve", middleware.RequireRole("admin"), adminHandler.ApproveOrganizer)
			admin.POST("/organizers/:id/reject", middleware.RequireRole("admin"), adminHandler.RejectOrganizer)

			// Настройки (только admin)
			admin.GET("/settings/difficulty", middleware.RequireRole("admin"), adminHandler.GetDifficultySettings)
			admin.PUT("/settings/difficulty", middleware.RequireRole("admin"), adminHandler.UpdateDifficultySettings)

			// Статистика (admin и observer)
			admin.GET("/stats", middleware.RequireRole("admin", "observer"), adminHandler.GetStats)

			// Управление пользователями (только admin)
			admin.GET("/users", middleware.RequireRole("admin"), adminHandler.GetAllUsers)
			admin.DELETE("/users/:id", middleware.RequireRole("admin"), adminHandler.DeleteUser)
			admin.PUT("/users/:id/role", middleware.RequireRole("admin"), adminHandler.ChangeUserRole)

			// Управление мероприятиями (только admin)
			admin.GET("/events", middleware.RequireRole("admin"), adminHandler.GetAllEvents)
			admin.DELETE("/events/:id", middleware.RequireRole("admin"), adminHandler.DeleteEventByAdmin)
		}

		// Эндпоинты фильтров (только для наблюдателей)
		filters := api.Group("/filters")
		filters.Use(middleware.AuthMiddleware(cfg.JWTSecret), middleware.RequireObserver())
		{
			filters.POST("", filterHandler.CreateFilter)       // создать
			filters.GET("", filterHandler.GetMyFilters)        // список моих
			filters.GET("/:id", filterHandler.GetFilter)       // получить
			filters.PUT("/:id", filterHandler.UpdateFilter)    // обновить
			filters.DELETE("/:id", filterHandler.DeleteFilter) // удалить
		}

		// Search
		api.GET("/search", searchHandler.Search)
	}

	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
