package repository

import (

	"github.com/Aftaza/sprintaza_backend/internal/model"
	"gorm.io/gorm"
)

// UserRepository mendefinisikan operasi database untuk User.
type UserRepository interface {
	FindByEmail(email string) (*model.User, error)
	FindByID(id uint) (*model.User, error)
	Create(user *model.User) (*model.User, error)
	Update(user *model.User) error
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

// FindByID mencari pengguna berdasarkan ID primary key.
// Ini akan preload semua data terkait yang dibutuhkan untuk profil.
func (r *userRepository) FindByID(id uint) (*model.User, error) {
    var user model.User
    // Preload UserXP dan UserAchievements beserta detail Achievement-nya.
    err := r.db.Preload("UserXP").Preload("UserAchievements.Achievement").Where("id = ?", id).First(&user).Error
    if err != nil {
        return nil, err // Kembalikan error apa adanya, termasuk gorm.ErrRecordNotFound
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

// Update menyimpan perubahan pada model User ke database.
func (r *userRepository) Update(user *model.User) error {
    // .Save akan memperbarui semua kolom, termasuk yang nilainya nol/kosong.
    return r.db.Save(user).Error
}