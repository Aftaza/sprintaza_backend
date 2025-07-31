package loginAuth

import (
	"errors"
	"time"

	util "github.com/Aftaza/sprintaza_backend/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Service struct {
	repository *Repository
	jwtSecret  string
}

func NewService(repository *Repository, jwtSecret string) *Service {
	return &Service{
		repository: repository,
		jwtSecret:  jwtSecret,
	}
}

// Login authenticates a user and returns a JWT token
func (s *Service) Login(input *LoginInput) (*LoginResponse, error) {
	// Find user by email
	user, err := s.repository.GetUserByEmail(input.Email)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"email": input.Email,
			"error": err.Error(),
		}).Warn("User not found during login attempt")
		return nil, errors.New("invalid email or password")
	}

	// Compare password
	if err := util.ComparePassword(user.PasswordHash, input.Password); err != nil {
		logrus.WithFields(logrus.Fields{
			"email": input.Email,
		}).Warn("Invalid password during login attempt")
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := s.generateJWTToken(user.ID, user.Email)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id": user.ID,
			"error":   err.Error(),
		}).Error("Failed to generate JWT token during login")
		return nil, errors.New("failed to generate authentication token")
	}

	logrus.WithFields(logrus.Fields{
		"user_id": user.ID,
		"email":   user.Email,
	}).Info("User logged in successfully")

	return &LoginResponse{
		Message: "Login successful",
		User: UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			AvatarURL: user.AvatarURL,
		},
		Token: token,
	}, nil
}

// generateJWTToken creates a JWT token for the user
func (s *Service) generateJWTToken(userID uint, email string) (string, error) {
	// Create the claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days expiration
		"iat":     time.Now().Unix(),
		"jti":     uuid.New().String(), // JWT ID for uniqueness
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// LoginResponse represents the response for a successful login
type LoginResponse struct {
	Message string       `json:"message"`
	User    UserResponse `json:"user"`
	Token   string       `json:"token"`
}

// UserResponse represents user data in API responses
type UserResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}
