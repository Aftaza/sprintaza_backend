package repository

import (
	"github.com/Aftaza/sprintaza_backend/internal/model"
	"gorm.io/gorm"
)

type AchievementRepository interface {
	FindByName(name string) (*model.Achievement, error)
}

type achievementRepository struct {
	db *gorm.DB
}

func NewAchievementRepository(db *gorm.DB) AchievementRepository {
	return &achievementRepository{db: db}
}

func (r *achievementRepository) FindByName(name string) (*model.Achievement, error) {
	var achievement model.Achievement
	err := r.db.Where("name = ?", name).First(&achievement).Error
	return &achievement, err
}