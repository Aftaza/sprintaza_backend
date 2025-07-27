package main

import (
	"fmt"
	"os"

	config "github.com/Aftaza/sprintaza_backend/configs"
	util "github.com/Aftaza/sprintaza_backend/utils"
	
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func init() {
	// Set up logging
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
	
	// Load environment variables
	if os.Getenv("GO_ENV") != "production" {
		logrus.SetLevel(logrus.DebugLevel)
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}

func main() {
	// Initialize database connection
	db = config.Connection()
	if db == nil {
		logrus.Fatal("Failed to connect to database")
	}

	// Initialize Gin router
	router := gin.New()

	// Add middleware
	setupMiddleware(router)

	// Setup routes
	setupRoutes(router)

	// Get port from environment or use default
	port := util.GodotEnv("GO_PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	logrus.Infof("Starting server on port %s", port)
	logrus.Infof("Server mode: %s", gin.Mode())
	
	if err := router.Run(":" + port); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}

func setupMiddleware(router *gin.Engine) {
	// Recovery middleware
	router.Use(gin.Recovery())

	// Logger middleware
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format("02/Jan/2006:15:04:05 -0700"),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	// CORS middleware
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With"}
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

	// Gzip compression middleware
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	// Security headers
	router.Use(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	})
}

func setupRoutes(router *gin.Engine) {
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
		// Auth routes
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", func(c *gin.Context) {
				c.JSON(501, gin.H{"message": "Register endpoint not implemented yet"})
			})
			authGroup.POST("/login", func(c *gin.Context) {
				c.JSON(501, gin.H{"message": "Login endpoint not implemented yet"})
			})
			authGroup.POST("/logout", func(c *gin.Context) {
				c.JSON(501, gin.H{"message": "Logout endpoint not implemented yet"})
			})
		}

		// User routes
		userGroup := v1.Group("/users")
		{
			userGroup.GET("/", func(c *gin.Context) {
				c.JSON(501, gin.H{"message": "Get users endpoint not implemented yet"})
			})
			userGroup.GET("/:id", func(c *gin.Context) {
				c.JSON(501, gin.H{"message": "Get user by ID endpoint not implemented yet"})
			})
			userGroup.POST("/", func(c *gin.Context) {
				c.JSON(501, gin.H{"message": "Create user endpoint not implemented yet"})
			})
			userGroup.PUT("/:id", func(c *gin.Context) {
				c.JSON(501, gin.H{"message": "Update user endpoint not implemented yet"})
			})
			userGroup.DELETE("/:id", func(c *gin.Context) {
				c.JSON(501, gin.H{"message": "Delete user endpoint not implemented yet"})
			})
		}

		// Project routes
		projectGroup := v1.Group("/projects")
		{
			projectGroup.GET("/", func(c *gin.Context) {
				c.JSON(501, gin.H{"message": "Get projects endpoint not implemented yet"})
			})
			projectGroup.GET("/:id", func(c *gin.Context) {
				c.JSON(501, gin.H{"message": "Get project by ID endpoint not implemented yet"})
			})
			projectGroup.POST("/", func(c *gin.Context) {
				c.JSON(501, gin.H{"message": "Create project endpoint not implemented yet"})
			})
			projectGroup.PUT("/:id", func(c *gin.Context) {
				c.JSON(501, gin.H{"message": "Update project endpoint not implemented yet"})
			})
			projectGroup.DELETE("/:id", func(c *gin.Context) {
				c.JSON(501, gin.H{"message": "Delete project endpoint not implemented yet"})
			})

			// Project member routes
			projectGroup.POST("/:id/members", func(c *gin.Context) {
				c.JSON(501, gin.H{"message": "Add project member endpoint not implemented yet"})
			})
			projectGroup.DELETE("/:id/members/:userId", func(c *gin.Context) {
				c.JSON(501, gin.H{"message": "Remove project member endpoint not implemented yet"})
			})
		}

		// Task routes
		taskGroup := v1.Group("/tasks")
		{
			taskGroup.GET("/", func(c *gin.Context) {
				c.JSON(501, gin.H{"message": "Get tasks endpoint not implemented yet"})
			})
			taskGroup.GET("/:id", func(c *gin.Context) {
				c.JSON(501, gin.H{"message": "Get task by ID endpoint not implemented yet"})
			})
			taskGroup.POST("/", func(c *gin.Context) {
				c.JSON(501, gin.H{"message": "Create task endpoint not implemented yet"})
			})
			taskGroup.PUT("/:id", func(c *gin.Context) {
				c.JSON(501, gin.H{"message": "Update task endpoint not implemented yet"})
			})
			taskGroup.DELETE("/:id", func(c *gin.Context) {
				c.JSON(501, gin.H{"message": "Delete task endpoint not implemented yet"})
			})

			// Subtask routes
			taskGroup.POST("/:id/subtasks", func(c *gin.Context) {
				c.JSON(501, gin.H{"message": "Create subtask endpoint not implemented yet"})
			})
			taskGroup.PUT("/:id/subtasks/:subtaskId", func(c *gin.Context) {
				c.JSON(501, gin.H{"message": "Update subtask endpoint not implemented yet"})
			})
			taskGroup.DELETE("/:id/subtasks/:subtaskId", func(c *gin.Context) {
				c.JSON(501, gin.H{"message": "Delete subtask endpoint not implemented yet"})
			})
		}

		// Achievement routes
		achievementGroup := v1.Group("/achievements")
		{
			achievementGroup.GET("/", func(c *gin.Context) {
				c.JSON(501, gin.H{"message": "Get achievements endpoint not implemented yet"})
			})
			achievementGroup.GET("/:id", func(c *gin.Context) {
				c.JSON(501, gin.H{"message": "Get achievement by ID endpoint not implemented yet"})
			})
		}
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

// GetDB returns the database instance for use in other packages
func GetDB() *gorm.DB {
	return db
}