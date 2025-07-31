package loginAuth

import (
	model "github.com/Aftaza/sprintaza_backend/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetUserByEmail(email string) (*model.EntityUsers, error) {
	var user model.EntityUsers
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
