package registerAuth

import (
	model "github.com/Aftaza/sprintaza_backend/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// CheckUserExists checks if a user with the given email already exists
func (r *Repository) CheckUserExists(email string) (*model.EntityUsers, error) {
	var user model.EntityUsers
	err := r.db.Where("email = ?", email).First(&user).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // User doesn't exist, which is expected for registration
		}
		logrus.WithFields(logrus.Fields{
			"email": email,
			"error": err.Error(),
		}).Error("Failed to check if user exists")
		return nil, err
	}
	
	return &user, nil // User already exists
}

// CreateUser creates a new user in the database
func (r *Repository) CreateUser(user *model.EntityUsers) error {
	// Start a transaction to ensure both user and user_xp are created together
	tx := r.db.Begin()
	
	// Create UserXP first
	userXP := &model.EntityUserXP{
		TotalXP: 0,
	}
	
	if err := tx.Create(userXP).Error; err != nil {
		tx.Rollback()
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to create user XP")
		return err
	}
	
	// Set the UserXPID in the user entity
	user.UserXPID = userXP.ID
	
	// Create the user
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		logrus.WithFields(logrus.Fields{
			"email": user.Email,
			"error": err.Error(),
		}).Error("Failed to create user")
		return err
	}
	
	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"email": user.Email,
			"error": err.Error(),
		}).Error("Failed to commit user creation transaction")
		return err
	}
	
	logrus.WithFields(logrus.Fields{
		"user_id": user.ID,
		"email":   user.Email,
		"name":    user.Name,
	}).Info("User created successfully")
	
	return nil
}

// GetUserByID retrieves a user by their ID with UserXP relationship
func (r *Repository) GetUserByID(userID uint) (*model.EntityUsers, error) {
	var user model.EntityUsers
	err := r.db.Preload("UserXP").Where("user_id = ?", userID).First(&user).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		logrus.WithFields(logrus.Fields{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("Failed to get user by ID")
		return nil, err
	}
	
	return &user, nil
}