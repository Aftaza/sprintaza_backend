package repository

import (

	"github.com/Aftaza/sprintaza_backend/internal/model"
	"gorm.io/gorm"
)

// UserRepository mendefinisikan operasi database untuk User.
type UserRepository interface {
	FindByEmail(email string) (*model.User, error)
	Create(user *model.User) (*model.User, error)
}

// userRepository adalah implementasi dari UserRepository menggunakan GORM.
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository membuat instance baru dari userRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// FindByEmail mencari pengguna berdasarkan alamat email mereka.
// Ini juga akan me-load (preload) data UserXP yang terkait.
func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	// Menggunakan Preload untuk secara otomatis mengambil data UserXP yang berelasi.
	err := r.db.Preload("UserXP").Where("email = ?", email).First(&user).Error
	// Jangan lagi menangani ErrRecordNotFound secara khusus.
	// Biarkan error ini "naik" ke service layer untuk ditangani di sana.
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create membuat entri pengguna baru beserta entri UserXP terkait dalam satu transaksi.
// --- PERBAIKAN DI SINI ---
func (r *userRepository) Create(user *model.User) (*model.User, error) {
	// Memulai transaksi untuk memastikan konsistensi data.
	// Jika ada error di tengah jalan, semua operasi akan dibatalkan (rollback).
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Buat entri UserXP terlebih dahulu.
		userXP := model.UserXP{XP: 0, InitialXP: 0} // Nilai awal
		if err := tx.Create(&userXP).Error; err != nil {
			return err // Rollback jika gagal
		}

		// 2. Hubungkan UserXP yang baru dibuat ke pengguna baru.
		user.UserXPID = userXP.ID
		user.UserXP = userXP

		// 3. Buat entri pengguna.
		if err := tx.Create(user).Error; err != nil {
			return err // Rollback jika gagal
		}
		
		// Jika tidak ada error, transaksi akan di-commit.
		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}