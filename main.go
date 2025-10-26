package main

import (
	"log"
	"os"

	db "deliveryAppBackend/config"
	routes "deliveryAppBackend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'")
		c.Next()
	}
}

func main() {
	// 📦 Load environment variables
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("❌ Error loading .env file")
	// }

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("❌ MONGO_URI is not set in environment variables")
	}

	// 🔌 Connect to MongoDB
	db.ConnectMongoDB(mongoURI)

	// 🚀 Setup Gin
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	router.Use(SecurityHeaders())

	// Setup routes
	routes.SetupRoutes(router)

	log.Printf("🚀 Delivery App Backend Server starting on port %s", port)
	router.Run(":" + port)
}

