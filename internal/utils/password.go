package util

import (
	"crypto/rand"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword menggunakan bcrypt untuk membuat hash dari password.
func HashPassword(password string) (string, error) {
	// bcrypt.DefaultCost adalah pilihan yang baik dan seimbang antara keamanan dan kecepatan.
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// ComparePassword membandingkan hash password dari database dengan password mentah yang diinput pengguna.
// Mengembalikan nil jika cocok, dan error jika tidak cocok.
func ComparePassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

// GenerateRandomPassword membuat password acak yang aman secara kriptografis.
// Ini digunakan untuk pendaftaran pertama kali via Google OAuth.
func GenerateRandomPassword(length int) (string, error) {
	// Kumpulan karakter yang akan digunakan untuk password.
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	password := make([]byte, length)
	
	for i := 0; i < length; i++ {
		// Menggunakan crypto/rand untuk keamanan, bukan math/rand.
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		password[i] = chars[num.Int64()]
	}

	return string(password), nil
}