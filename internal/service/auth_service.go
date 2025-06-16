package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Aftaza/sprintaza_backend/internal/model"
	"github.com/Aftaza/sprintaza_backend/internal/repository"
	"github.com/Aftaza/sprintaza_backend/internal/utils"
	"golang.org/x/oauth2"
	googleOAuth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

// AuthService mendefinisikan interface untuk logika autentikasi.
type AuthService interface {
	ProcessGoogleCallback(ctx context.Context, code string) (string, error)
	LoginWithPassword(email, password string) (string, error)
}

type authService struct {
	userRepo        repository.UserRepository
	googleOAuthConfig *oauth2.Config
}

// NewAuthService membuat instance baru dari AuthService.
func NewAuthService(userRepo repository.UserRepository, googleOAuthConfig *oauth2.Config) AuthService {
	return &authService{
		userRepo:        userRepo,
		googleOAuthConfig: googleOAuthConfig,
	}
}

// ProcessGoogleCallback menangani logika setelah menerima callback dari Google.
func (s *authService) ProcessGoogleCallback(ctx context.Context, code string) (string, error) {
	// 1. Tukar 'code' dengan token akses dari Google
	token, err := s.googleOAuthConfig.Exchange(ctx, code)
	if err != nil {
		return "", fmt.Errorf("gagal menukar kode oauth: %w", err)
	}

	// 2. Gunakan token untuk mendapatkan informasi pengguna dari Google
	client := s.googleOAuthConfig.Client(ctx, token)
	oauth2Service, err := googleOAuth2.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return "", fmt.Errorf("gagal membuat service oauth2: %w", err)
	}

	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		return "", fmt.Errorf("gagal mendapatkan info pengguna: %w", err)
	}

	// 3. Cek apakah user sudah terdaftar di database kita
	user, err := s.userRepo.FindByEmail(userInfo.Email)
	if err != nil {
		// Cek secara spesifik apakah errornya adalah 'record not found'.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Pengguna dengan email %s belum terdaftar, membuat akun baru...", userInfo.Email)
			
			// Logika pembuatan user baru dipindahkan ke sini.
			randomPassword, randErr := util.GenerateRandomPassword(16)
			if randErr != nil { return "", fmt.Errorf("gagal membuat password acak: %w", randErr) }

			hashedPassword, hashErr := util.HashPassword(randomPassword)
			if hashErr != nil { return "", fmt.Errorf("gagal melakukan hash password: %w", hashErr) }

			newUser := &model.User{
				Name: userInfo.Name, Email: userInfo.Email,
				PasswordHash: hashedPassword, AuthProvider: "google",
			}

			createdUser, createErr := s.userRepo.Create(newUser)
			if createErr != nil { return "", fmt.Errorf("gokal membuat user baru: %w", createErr) }
			
			user = createdUser // Set user saat ini ke user yang baru dibuat
			
			log.Printf("Password sementara untuk %s: %s (HANYA UNTUK DEV)", user.Email, randomPassword)

		} else {
			// Jika errornya bukan 'record not found', berarti ini masalah database lain.
			return "", fmt.Errorf("error saat mencari user: %w", err)
		}
	}

	// Jika tidak ada error, 'user' pasti sudah terisi (baik dari hasil find atau create).
	jwtToken, err := util.GenerateJWT(user.ID)
	if err != nil {
		return "", fmt.Errorf("gagal membuat token JWT: %w", err)
	}

	return jwtToken, nil
}

// LoginWithPassword memverifikasi kredensial email/password.
func (s *authService) LoginWithPassword(email, password string) (string, error) {
	// Cari user berdasarkan email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		// Jangan beri tahu penyerang apakah email ada atau tidak.
		return "", errors.New("kredensial tidak valid")
	}

	// Bandingkan password
	if err := util.ComparePassword(user.PasswordHash, password); err != nil {
		return "", errors.New("kredensial tidak valid")
	}

	// Buat JWT jika berhasil
	jwtToken, err := util.GenerateJWT(user.ID)
	if err != nil {
		return "", fmt.Errorf("gagal membuat token JWT: %w", err)
	}

	return jwtToken, nil
}