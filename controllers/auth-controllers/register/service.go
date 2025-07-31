package registerAuth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	model "github.com/Aftaza/sprintaza_backend/models"
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

// RegisterUserInput defines the input for standard user registration
type RegisterUserInput struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// Validate checks if the input fields are valid
func (i *RegisterUserInput) Validate() error {
	if i.Name == "" {
		return &ValidationError{Field: "name", Message: "Name is required"}
	}
	if i.Email == "" {
		return &ValidationError{Field: "email", Message: "Email is required"}
	}
	return nil
}

// RegisterUser handles standard user registration
func (s *Service) RegisterUser(input *RegisterUserInput) (*RegisterResponse, error) {
	// Validate input
	if err := input.Validate(); err != nil {
		logrus.WithFields(logrus.Fields{
			"email": input.Email,
			"error": err.Error(),
		}).Warn("Invalid registration input")
		return nil, err
	}

	// Check if user already exists
	existingUser, err := s.repository.CheckUserExists(input.Email)
	if err != nil {
		return nil, errors.New("failed to check user existence")
	}

	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Generate a random password
	password, err := util.GenerateRandomPassword(12)
	if err != nil {
		return nil, errors.New("failed to generate secure password")
	}

	// Create new user entity
	newUser := &model.EntityUsers{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: password, // This will be hashed by the BeforeCreate hook
	}

	// Create user in database
	if err := s.repository.CreateUser(newUser); err != nil {
		return nil, errors.New("failed to create user account")
	}

	// Send welcome email
	emailData := map[string]string{
		"Name":     newUser.Name,
		"Email":    newUser.Email,
		"Password": password,
	}
	go func() {
		if err := util.SendEmail(newUser.Name, newUser.Email, "Welcome to Sprintaza!", "templates/email/registration_email.html", emailData); err != nil {
			logrus.WithFields(logrus.Fields{
				"user_id": newUser.ID,
				"email":   newUser.Email,
				"error":   err.Error(),
			}).Error("Failed to send welcome email")
		}
	}()

	// Award welcome achievement
	var welcomeAchievement *AchievementResponse
	if achievementData, err := s.awardWelcomeAchievement(newUser.ID); err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id": newUser.ID,
			"error":   err.Error(),
		}).Warn("Failed to award welcome achievement, but user registration succeeded")
	} else {
		welcomeAchievement = achievementData
	}

	// Generate JWT token
	token, err := s.generateJWTToken(newUser.ID, newUser.Email)
	if err != nil {
		return nil, errors.New("failed to generate authentication token")
	}

	logrus.WithFields(logrus.Fields{
		"user_id": newUser.ID,
		"email":   newUser.Email,
		"name":    newUser.Name,
	}).Info("New user registered successfully")

	return &RegisterResponse{
		Message: "User registered successfully",
		User: UserResponse{
			ID:    newUser.ID,
			Name:  newUser.Name,
			Email: newUser.Email,
		},
		Token:       token,
		IsNewUser:   true,
		Achievement: welcomeAchievement,
	}, nil
}

// RegisterUserWithGoogleOAuth handles Google OAuth user registration
func (s *Service) RegisterUserWithGoogleOAuth(input *GoogleOAuthRegisterInput) (*RegisterResponse, error) {
	// Validate input
	if err := input.Validate(); err != nil {
		logrus.WithFields(logrus.Fields{
			"email": input.Email,
			"error": err.Error(),
		}).Warn("Invalid registration input")
		return nil, err
	}

	// Check if user already exists
	existingUser, err := s.repository.CheckUserExists(input.Email)
	if err != nil {
		return nil, errors.New("failed to check user existence")
	}

	// If user exists, return login response instead of creating new user
	if existingUser != nil {
		logrus.WithFields(logrus.Fields{
			"email":   input.Email,
			"user_id": existingUser.ID,
		}).Info("User already exists, returning login response")

		// Generate JWT token for existing user
		token, err := s.generateJWTToken(existingUser.ID, existingUser.Email)
		if err != nil {
			return nil, errors.New("failed to generate authentication token")
		}

		return &RegisterResponse{
			Message: "User already exists, logged in successfully",
			User: UserResponse{
				ID:        existingUser.ID,
				Name:      existingUser.Name,
				Email:     existingUser.Email,
				AvatarURL: existingUser.AvatarURL,
			},
			Token:     token,
			IsNewUser: false,
		}, nil
	}

	// Create new user entity
	// For Google OAuth, we don't need a password, so we generate a random one
	randomPassword, err := s.generateRandomPassword()
	if err != nil {
		return nil, errors.New("failed to generate secure password")
	}

	newUser := &model.EntityUsers{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: randomPassword, // This will be hashed by the BeforeCreate hook
		AvatarURL:    "",             // Can be set later from Google profile
	}

	// Create user in database
	if err := s.repository.CreateUser(newUser); err != nil {
		return nil, errors.New("failed to create user account")
	}

	// Award welcome achievement (ID: 1) for new users
	var welcomeAchievement *AchievementResponse
	if achievementData, err := s.awardWelcomeAchievement(newUser.ID); err != nil {
		// Log the error but don't fail the registration
		logrus.WithFields(logrus.Fields{
			"user_id": newUser.ID,
			"error":   err.Error(),
		}).Warn("Failed to award welcome achievement, but user registration succeeded")
	} else {
		welcomeAchievement = achievementData
	}

	// Generate JWT token for new user
	token, err := s.generateJWTToken(newUser.ID, newUser.Email)
	if err != nil {
		return nil, errors.New("failed to generate authentication token")
	}

	logrus.WithFields(logrus.Fields{
		"user_id": newUser.ID,
		"email":   newUser.Email,
		"name":    newUser.Name,
	}).Info("New user registered successfully via Google OAuth")

	return &RegisterResponse{
		Message: "User registered successfully",
		User: UserResponse{
			ID:        newUser.ID,
			Name:      newUser.Name,
			Email:     newUser.Email,
			AvatarURL: newUser.AvatarURL,
		},
		Token:       token,
		IsNewUser:   true,
		Achievement: welcomeAchievement,
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
		logrus.WithFields(logrus.Fields{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("Failed to generate JWT token")
		return "", err
	}

	return tokenString, nil
}

// generateRandomPassword generates a secure random password for OAuth users
func (s *Service) generateRandomPassword() (string, error) {
	bytes := make([]byte, 32) // 32 bytes = 256 bits
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// awardWelcomeAchievement awards the welcome achievement (ID: 1) to a new user and returns achievement details
func (s *Service) awardWelcomeAchievement(userID uint) (*AchievementResponse, error) {
	const welcomeAchievementID = 1

	// Award the achievement
	if err := s.repository.AwardAchievement(userID, welcomeAchievementID); err != nil {
		return nil, err
	}

	// Get achievement details to award XP
	achievement, err := s.repository.GetAchievementByID(welcomeAchievementID)
	if err != nil {
		return nil, err
	}

	if achievement == nil {
		return nil, errors.New("welcome achievement not found")
	}

	// Update user's XP with achievement reward
	if err := s.repository.UpdateUserXP(userID, achievement.XPReward); err != nil {
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"user_id":        userID,
		"achievement_id": welcomeAchievementID,
		"xp_reward":      achievement.XPReward,
	}).Info("Welcome achievement awarded successfully")

	// Return achievement details for response
	return &AchievementResponse{
		ID:          achievement.ID,
		Name:        achievement.Name,
		Description: achievement.Description,
		XPReward:    achievement.XPReward,
		IconURL:     achievement.IconURL,
	}, nil
}

// RegisterResponse represents the response for user registration
type RegisterResponse struct {
	Message     string                  `json:"message"`
	User        UserResponse            `json:"user"`
	Token       string                  `json:"token"`
	IsNewUser   bool                    `json:"is_new_user"`
	Achievement *AchievementResponse    `json:"achievement,omitempty"` // Only for new users
}

// UserResponse represents user data in API responses
type UserResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// AchievementResponse represents achievement data in API responses
type AchievementResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	XPReward    int    `json:"xp_reward"`
	IconURL     string `json:"icon_url"`
}
