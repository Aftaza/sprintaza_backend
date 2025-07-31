package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes initializes all routes for the application
func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "OK",
			"message":   "Sprintaza Backend is running",
			"timestamp": gin.H{},
		})
	})

	// API version 1 group
	v1 := router.Group("/api/v1")
	{
		// Setup route groups
		setupAuthRoutes(v1, db)
		setupUserRoutes(v1, db)
		setupProjectRoutes(v1, db)
		setupTaskRoutes(v1, db)
		setupAchievementRoutes(v1, db)
	}

	// 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"error":   "Not Found",
			"message": "The requested endpoint does not exist",
			"path":    c.Request.URL.Path,
		})
	})

	// 405 handler
	router.NoMethod(func(c *gin.Context) {
		c.JSON(405, gin.H{
			"error":   "Method Not Allowed",
			"message": "The requested method is not allowed for this endpoint",
			"method":  c.Request.Method,
			"path":    c.Request.URL.Path,
		})
	})
}