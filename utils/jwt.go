package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Struktur payload token
type JWTClaim struct {
	UserID               string `json:"user_id"` // UUID string
	Email                string `json:"email"`   // email user
	jwt.RegisteredClaims        // claim bawaan (exp, iss, dll)
}

// Generate token JWT
// GenerateToken → membuat token JWT untuk user
func GenerateToken(userId uuid.UUID, email string) (string, error) {

	// expired token 24 jam dari sekarang
	expirationTime := time.Now().Add(24 * time.Hour)

	// isi payload
	claims := &JWTClaim{
		UserID: userId.String(), // konversi uuid → string
		Email:  email,

		// claim bawaan JWT
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime), // batas waktu
		},
	}

	// generate token dengan metode HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// ambil secret dari ENV
	secret := os.Getenv("JWT_SECRET")

	// hasilkan token string
	return token.SignedString([]byte(secret))
}

// Validasi JWT
func ValidateToken(tokenString string) (*JWTClaim, error) {
	claims := &JWTClaim{}

	secret := os.Getenv("JWT_SECRET")

	// Parse dan cek validitas
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}
