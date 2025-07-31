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

// Register handles standard user registration
func (h *Handler) Register(c *gin.Context) {
	var input registerAuth.RegisterUserInput

	// Bind JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("Invalid JSON input for registration")

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
	}).Info("Registration attempt")

	// Process registration
	response, err := h.service.RegisterUser(&input)
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
		}).Error("Failed to process registration")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Registration failed",
			"message": "Unable to process registration at this time",
		})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    response,
	})

	// Log successful registration
	logrus.WithFields(logrus.Fields{
		"user_id": response.User.ID,
		"email":   response.User.Email,
	}).Info("Registration completed successfully")
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

// @Summary Register a new user
// @Description Registers a new user with their name and email.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body registerAuth.RegisterUserInput true "User registration details"
// @Success 201 {object} registerAuth.RegisterResponse "User registered successfully"
// @Failure 400 {object} gin.H "Invalid input or validation error"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /auth/register [post]
func Register(db *gorm.DB) gin.HandlerFunc {
	return NewHandler(db).Register
}

// @Summary Register a new user with Google OAuth
// @Description Registers a new user or logs in an existing user using Google OAuth details.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body registerAuth.GoogleOAuthRegisterInput true "User registration details"
// @Success 201 {object} registerAuth.RegisterResponse "User registered successfully"
// @Success 200 {object} registerAuth.RegisterResponse "User already exists, logged in successfully"
// @Failure 400 {object} gin.H "Invalid input or validation error"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /auth/register/google [post]
func GoogleOAuthRegister(db *gorm.DB) gin.HandlerFunc {
	return NewHandler(db).GoogleOAuthRegister
}
