package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/Aftaza/sprintaza_backend/internal/utils"
)

// AuthMiddleware adalah middleware untuk memvalidasi JWT dari Authorization header.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// 2. Cek format "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			return
		}

		tokenString := parts[1]

		// 3. Validasi token menggunakan fungsi dari util
		claims, err := util.ValidateJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// 4. Jika valid, simpan informasi user ke dalam context Gin
		// Ini memungkinkan handler berikutnya untuk mengetahui siapa pengguna yang sedang login.
		c.Set("userID", claims.UserID)

		// 5. Lanjutkan ke handler berikutnya
		c.Next()
	}
}