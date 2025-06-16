package util

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims adalah custom claims yang kita tambahkan ke token,
// selain claims standar dari RegisteredClaims.
type JWTClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateJWT membuat token JWT baru untuk user ID yang diberikan.
func GenerateJWT(userID uint) (string, error) {
	// Ambil konfigurasi dari environment variables
	secretKey := os.Getenv("JWT_SECRET_KEY")
	expirationHoursStr := os.Getenv("JWT_EXPIRATION_HOURS")
	
	expirationHours, err := strconv.Atoi(expirationHoursStr)
	if err != nil {
		// Fallback ke nilai default jika env var tidak valid atau tidak ada
		expirationHours = 72 
	}

	// Tentukan waktu kedaluwarsa token
	expirationTime := time.Now().Add(time.Duration(expirationHours) * time.Hour)

	// Buat claims
	claims := &JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Buat token baru dengan claims dan metode signing HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Tanda tangani token dengan secret key untuk menghasilkan string token final
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ValidateJWT memvalidasi string token dan mengembalikan claims jika valid.
func ValidateJWT(tokenString string) (*JWTClaims, error) {
	secretKey := os.Getenv("JWT_SECRET_KEY")

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validasi metode signing. Pastikan metodenya adalah yang kita harapkan (HS256).
		// Ini penting untuk mencegah serangan 'alg:none'.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, fmt.Errorf("invalid token")
}