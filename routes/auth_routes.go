package routes

import (
	handlerRegister "github.com/Aftaza/sprintaza_backend/handlers/auth-handlers/register"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// setupAuthRoutes configures all authentication-related routes
func setupAuthRoutes(v1 *gin.RouterGroup, db *gorm.DB) {
	// Initialize handlers
	registerHandler := handlerRegister.NewHandler(db)

	// Auth routes group
	authGroup := v1.Group("/auth")
	{
		// POST /api/v1/auth/register - Google OAuth registration
		authGroup.POST("/register", registerHandler.GoogleOAuthRegister)
		
		// POST /api/v1/auth/login - User login (not implemented)
		authGroup.POST("/login", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "Login endpoint not implemented yet"})
		})
		
		// POST /api/v1/auth/logout - User logout (not implemented)
		authGroup.POST("/logout", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "Logout endpoint not implemented yet"})
		})
		
		// GET /api/v1/auth/health - Auth service health check
		authGroup.GET("/health", registerHandler.HealthCheck)
	}
}