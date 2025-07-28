package config

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	model "github.com/Aftaza/sprintaza_backend/models"
	// util "github.com/Aftaza/sprintaza_backend/utils"
)

// SeedAchievements memeriksa apakah tabel ACHIEVEMENTS kosong dan mengisinya dengan data default.
func SeedAchievements(db *gorm.DB) {
	var count int64
	db.Model(&model.EntityAchievement{}).Count(&count)

	if count == 0 {
		logrus.Info("Seeding default achievements...")
		for _, achievement := range model.DefaultAchievements {
			result := db.Create(&achievement)
			if result.Error != nil {
				logrus.Fatalf("Failed to seed achievement '%s': %v\n", achievement.Name, result.Error)
				return // Hentikan seeding jika ada error
			}
		}
		logrus.Info("Default achievements seeded successfully.")
	} else {
		logrus.Info("achievements table already has data. Skipping seeding.")
	}
}

func SeedRoles(db *gorm.DB){
	var count int64
	db.Model(&model.EntityRole{}).Count(&count)

	if count == 0 {
		logrus.Info("Seeding default roles...")
		for _, role := range model.DefaultRoles {
			result := db.Create(&role)
			if result.Error != nil {
				logrus.Fatalf("Failed to seed achievement '%s': %v\n", role.Name, result.Error)
				return // Hentikan seeding jika ada error
			}
		}
		logrus.Info("Default roles seeded successfully.")
	} else {
		logrus.Info("roles table already has data. Skipping seeding.")
	}
}