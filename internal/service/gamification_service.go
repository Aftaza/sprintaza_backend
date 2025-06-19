package service

import (
	"log"
	
	"github.com/Aftaza/sprintaza_backend/internal/model"
	"github.com/Aftaza/sprintaza_backend/internal/repository"
)

type GamificationService interface {
	GrantWelcomeAchievement(user *model.User)
}

type gamificationService struct {
	userRepo   repository.UserRepository
	achieveRepo repository.AchievementRepository
}

func NewGamificationService(userRepo repository.UserRepository, achieveRepo repository.AchievementRepository) GamificationService {
	return &gamificationService{userRepo: userRepo, achieveRepo: achieveRepo}
}

func (s *gamificationService) GrantWelcomeAchievement(user *model.User) {
	// Cek apakah user sudah punya achievement "Selamat Datang!"
	// Ini untuk memastikan achievement hanya diberikan sekali.
	for _, ach := range user.UserAchievements {
		if ach.Achievement.Name == "Selamat Datang!" {
			return // User sudah punya, tidak perlu melakukan apa-apa.
		}
	}

	// Jika belum, ambil data achievement-nya dari DB
	welcomeAchieve, err := s.achieveRepo.FindByName("Selamat Datang!")
	if err != nil {
		log.Printf("Error: Gagal menemukan achievement 'Selamat Datang!': %v", err)
		return
	}

	// Berikan achievement dan XP
	err = s.userRepo.GrantAchievement(user.ID, welcomeAchieve.ID, welcomeAchieve.XpReward)
	if err != nil {
		log.Printf("Error: Gagal memberikan achievement ke user %d: %v", user.ID, err)
		return
	}

	log.Printf("Achievement 'Selamat Datang!' dan %d XP diberikan kepada user %s", welcomeAchieve.XpReward, user.Email)
}