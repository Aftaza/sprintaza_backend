package service

import (
	"fmt"

	"github.com/Aftaza/sprintaza_backend/internal/model"
	"github.com/Aftaza/sprintaza_backend/internal/repository"
	"github.com/Aftaza/sprintaza_backend/internal/utils"
)

type UserService interface {
	GetProfile(userID uint) (*model.User, error)
	UpdateName(userID uint, newName string) (*model.User, error)
	UpdatePassword(userID uint, oldPassword, newPassword string) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetProfile(userID uint) (*model.User, error) {
	// Cukup panggil repository, karena data sudah di-preload di sana.
	return s.userRepo.FindByID(userID)
}

func (s *userService) UpdateName(userID uint, newName string) (*model.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err // User tidak ditemukan atau error DB
	}

	user.Name = newName

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) UpdatePassword(userID uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	// 1. Verifikasi password lama
	if err := util.ComparePassword(user.PasswordHash, oldPassword); err != nil {
		return fmt.Errorf("password lama tidak sesuai")
	}

	// 2. Hash password baru
	newHashedPassword, err := util.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("gagal melakukan hash password baru: %w", err)
	}

	// 3. Update password di model dan simpan
	user.PasswordHash = newHashedPassword
	return s.userRepo.Update(user)
}