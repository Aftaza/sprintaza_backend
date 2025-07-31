package routes

import (
	handlerLogin "github.com/Aftaza/sprintaza_backend/handlers/auth-handlers/login"
	handlerRegister "github.com/Aftaza/sprintaza_backend/handlers/auth-handlers/register"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// setupAuthRoutes configures all authentication-related routes
func setupAuthRoutes(v1 *gin.RouterGroup, db *gorm.DB) {
	// Initialize handlers
	registerHandler := handlerRegister.NewHandler(db)
	loginHandler := handlerLogin.NewHandler(db)

	// Auth routes group
	authGroup := v1.Group("/auth")
	{
		// POST /api/v1/auth/register - Standard user registration
		authGroup.POST("/register", registerHandler.Register)

		// POST /api/v1/auth/register/google - Google OAuth registration
		authGroup.POST("/register/google", registerHandler.GoogleOAuthRegister)

		// POST /api/v1/auth/login - User login
		authGroup.POST("/login", loginHandler.Login)

		// POST /api/v1/auth/logout - User logout (not implemented)
		authGroup.POST("/logout", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "Logout endpoint not implemented yet"})
		})

		// GET /api/v1/auth/health - Auth service health check
		authGroup.GET("/health", registerHandler.HealthCheck)
	}
}