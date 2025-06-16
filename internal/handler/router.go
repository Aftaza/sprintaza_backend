
package handler

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/Aftaza/sprintaza_backend/internal/handler/http" // Handler spesifik
	// "github.com/Aftaza/sprintaza_backend/internal/middleware"
	"github.com/Aftaza/sprintaza_backend/internal/repository"
	"github.com/Aftaza/sprintaza_backend/internal/service"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

// SetupRoutes menginisialisasi semua dependency dan mendaftarkan semua route.
// Fungsi ini mengembalikan gin.Engine yang sudah siap dijalankan.
func SetupRoutes(db *gorm.DB) *gin.Engine {
	// --- Inisialisasi Engine Gin ---
	r := gin.Default()

	// TODO: Tambahkan middleware global jika perlu, misalnya CORS
	// r.Use(cors.Default())

	// --- Inisialisasi Dependencies (Dependency Injection) ---

	// 1. Repository
	userRepo := repository.NewUserRepository(db)
	// projectRepo := repository.NewProjectRepository(db) // (Untuk nanti)

	// 2. Konfigurasi Google OAuth
	googleOAuthConfig := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			// Tambahkan scope Google Calendar di sini nanti
			// "https://www.googleapis.com/auth/calendar.events.readonly", 
		},
		Endpoint: google.Endpoint,
	}

	// 3. Service
	authService := service.NewAuthService(userRepo, googleOAuthConfig)
	// projectService := service.NewProjectService(projectRepo) // (Untuk nanti)

	// 4. Handler
	authHandler := http.NewAuthHandler(authService, googleOAuthConfig)
	// projectHandler := http.NewProjectHandler(projectService) // (Untuk nanti)

	// --- Pendaftaran Route ---
	api := r.Group("/v1")
	{
		// Grup untuk route autentikasi (Publik)
		authRoutes := api.Group("/auth")
		{
			authRoutes.GET("/google/login", authHandler.GoogleLogin)
			authRoutes.GET("/google/callback", authHandler.GoogleCallback)
			authRoutes.POST("/login", authHandler.Login)
		}

		// Grup untuk route yang memerlukan autentikasi (Dilindungi Middleware)
		// protectedRoutes := api.Group("").Use(middleware.AuthMiddleware())
		// {
		// 	// Contoh route yang dilindungi
		// 	// projects := protectedRoutes.Group("/projects")
		// 	// {
		// 	// 	projects.POST("/", projectHandler.CreateProject)
		// 	// 	projects.GET("/", projectHandler.GetProjects)
		// 	// }

		// 	// Route untuk mendapatkan profil user yang sedang login
		// 	// protectedRoutes.GET("/me", userHandler.GetMyProfile)
		// }
	}

	return r
}