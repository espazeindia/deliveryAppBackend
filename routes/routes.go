package routes

import (
	"deliveryAppBackend/handlers"
	"deliveryAppBackend/infrastructure/mongodb"
	"deliveryAppBackend/middlewares"
	"deliveryAppBackend/usecase"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Initialize repositories
	partnerRepo := mongodb.NewDeliveryPartnerMongoRepository()
	deliveryRepo := mongodb.NewDeliveryMongoRepository()
	earningsRepo := mongodb.NewEarningsMongoRepository()

	// Initialize use cases
	authUseCase := usecase.NewAuthUseCase(partnerRepo)
	deliveryUseCase := usecase.NewDeliveryUseCase(deliveryRepo, partnerRepo, earningsRepo)
	profileUseCase := usecase.NewProfileUseCase(partnerRepo)
	earningsUseCase := usecase.NewEarningsUseCase(earningsRepo, deliveryRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authUseCase)
	deliveryHandler := handlers.NewDeliveryHandler(deliveryUseCase)
	profileHandler := handlers.NewProfileHandler(profileUseCase)
	earningsHandler := handlers.NewEarningsHandler(earningsUseCase)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "delivery-app-backend",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		delivery := v1.Group("/delivery")
		{
			// Authentication
			delivery.POST("/login", authHandler.Login)
			delivery.POST("/request-otp", authHandler.RequestOTP)
			delivery.POST("/verify-otp", authHandler.VerifyOTP)

			// Protected routes (authentication required)
			protected := delivery.Group("")
			protected.Use(middlewares.AuthMiddleware())
			{
				// Orders
				protected.GET("/orders/active", deliveryHandler.GetActiveOrders)
				protected.GET("/orders/history", deliveryHandler.GetOrderHistory)
				protected.GET("/orders/:id", deliveryHandler.GetOrderDetails)
				protected.POST("/orders/:id/accept", deliveryHandler.AcceptOrder)
				protected.POST("/orders/:id/status", deliveryHandler.UpdateOrderStatus)
				protected.POST("/orders/:id/complete", deliveryHandler.CompleteDelivery)

				// Profile
				protected.GET("/profile", profileHandler.GetProfile)
				protected.PUT("/profile", profileHandler.UpdateProfile)
				protected.POST("/location", profileHandler.UpdateLocation)
				protected.POST("/availability", profileHandler.ToggleAvailability)

				// Earnings
				protected.GET("/earnings", earningsHandler.GetEarnings)
				protected.GET("/earnings/history", earningsHandler.GetEarningsHistory)
			}
		}
	}
}

