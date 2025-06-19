package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Aftaza/sprintaza_backend/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetProfile mengambil profil pengguna yang sedang login.
func (h *UserHandler) GetProfile(c *gin.Context) {
	// Ambil userID dari context yang di-set oleh AuthMiddleware
	userID := c.MustGet("userID").(uint)

	user, err := h.userService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	// NOTE: Di aplikasi nyata, Anda mungkin ingin menggunakan DTO (Data Transfer Object)
	// untuk menyembunyikan field sensitif seperti PasswordHash sebelum mengirimkannya ke client.
	c.JSON(http.StatusOK, user)
}

// UpdateProfileRequest adalah struct untuk binding request update nama.
type UpdateProfileRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateProfile memperbarui nama pengguna.
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := h.userService.UpdateName(userID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui profil"})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// UpdatePasswordRequest adalah struct untuk binding request ganti password.
type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// UpdatePassword memperbarui password pengguna.
func (h *UserHandler) UpdatePassword(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userService.UpdatePassword(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		if err.Error() == "password lama tidak sesuai" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil diperbarui"})
}