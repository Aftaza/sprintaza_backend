package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Aftaza/sprintaza_backend/internal/service"
	"golang.org/x/oauth2"
)

// AuthHandler menangani request HTTP yang berhubungan dengan autentikasi.
type AuthHandler struct {
	authService       service.AuthService
	googleOAuthConfig *oauth2.Config
}

// NewAuthHandler membuat instance baru dari AuthHandler.
func NewAuthHandler(authService service.AuthService, googleOAuthConfig *oauth2.Config) *AuthHandler {
	return &AuthHandler{
		authService:       authService,
		googleOAuthConfig: googleOAuthConfig,
	}
}

// GoogleLogin mengarahkan pengguna ke halaman login Google.
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	// State bisa digunakan untuk mencegah serangan CSRF, untuk sekarang kita gunakan string statis.
	url := h.googleOAuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback menangani redirect dari Google setelah login.
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code not provided"})
		return
	}

	// Proses kode dan dapatkan JWT
	token, err := h.authService.ProcessGoogleCallback(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Kirim token sebagai respons
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// LoginRequest adalah struktur untuk binding request body login email/password.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login menangani permintaan login dengan email dan password.
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.LoginWithPassword(req.Email, req.Password)
	if err != nil {
		// Menggunakan StatusUnauthorized untuk kredensial yang salah.
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}