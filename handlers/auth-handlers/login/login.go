package handlerLogin

import (
	"net/http"

	loginAuth "github.com/Aftaza/sprintaza_backend/controllers/auth-controllers/login"
	util "github.com/Aftaza/sprintaza_backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Handler struct {
	service *loginAuth.Service
}

func NewHandler(db *gorm.DB) *Handler {
	jwtSecret := util.GodotEnv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-key-change-in-production"
		logrus.Warn("JWT_SECRET not set in environment, using default key")
	}

	repository := loginAuth.NewRepository(db)
	service := loginAuth.NewService(repository, jwtSecret)

	return &Handler{
		service: service,
	}
}

func (h *Handler) Login(c *gin.Context) {
	var input loginAuth.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.Login(&input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
