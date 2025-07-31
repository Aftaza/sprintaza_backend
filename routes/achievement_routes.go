package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// setupAchievementRoutes configures all achievement-related routes
func setupAchievementRoutes(v1 *gin.RouterGroup, db *gorm.DB) {
	// TODO: Initialize achievement handlers when implemented
	// achievementHandler := achievementHandler.NewHandler(db)

	// Achievement routes group
	achievementGroup := v1.Group("/achievements")
	{
		// GET /api/v1/achievements - Get all achievements (not implemented)
		achievementGroup.GET("/", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "Get achievements endpoint not implemented yet"})
		})
		
		// GET /api/v1/achievements/:id - Get achievement by ID (not implemented)
		achievementGroup.GET("/:id", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "Get achievement by ID endpoint not implemented yet"})
		})
	}
}