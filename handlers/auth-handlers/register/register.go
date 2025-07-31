package handlerRegister

import (
	"net/http"
	"os"

	registerAuth "github.com/Aftaza/sprintaza_backend/controllers/auth-controllers/register"
	util "github.com/Aftaza/sprintaza_backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Handler struct {
	service *registerAuth.Service
}

func NewHandler(db *gorm.DB) *Handler {
	// Get JWT secret from environment
	jwtSecret := util.GodotEnv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-key-change-in-production" // Fallback for development
		logrus.Warn("JWT_SECRET not set in environment, using default key")
	}

	repository := registerAuth.NewRepository(db)
	service := registerAuth.NewService(repository, jwtSecret)
	
	return &Handler{
		service: service,
	}
}

// GoogleOAuthRegister handles Google OAuth user registration
func (h *Handler) GoogleOAuthRegister(c *gin.Context) {
	var input registerAuth.GoogleOAuthRegisterInput
	
	// Bind JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("Invalid JSON input for Google OAuth registration")
		
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"message": "Please provide valid email and name",
			"details": err.Error(),
		})
		return
	}

	// Log the registration attempt
	logrus.WithFields(logrus.Fields{
		"email": input.Email,
		"name":  input.Name,
	}).Info("Google OAuth registration attempt")

	// Process registration
	response, err := h.service.RegisterUserWithGoogleOAuth(&input)
	if err != nil {
		// Check if it's a validation error
		if validationErr, ok := err.(*registerAuth.ValidationError); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"message": validationErr.Message,
				"field":   validationErr.Field,
			})
			return
		}

		// Handle other errors
		logrus.WithFields(logrus.Fields{
			"email": input.Email,
			"error": err.Error(),
		}).Error("Failed to process Google OAuth registration")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Registration failed",
			"message": "Unable to process registration at this time",
		})
		return
	}

	// Determine HTTP status code based on whether user is new or existing
	statusCode := http.StatusCreated
	if !response.IsNewUser {
		statusCode = http.StatusOK
	}

	// Return success response
	c.JSON(statusCode, gin.H{
		"success": true,
		"data":    response,
	})

	// Log successful registration/login
	logrus.WithFields(logrus.Fields{
		"user_id":     response.User.ID,
		"email":       response.User.Email,
		"is_new_user": response.IsNewUser,
	}).Info("Google OAuth registration/login completed successfully")
}

// HealthCheck provides a health check endpoint for the auth service
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "OK",
		"service":   "Authentication Service",
		"timestamp": gin.H{},
		"version":   os.Getenv("APP_VERSION"),
	})
}