package utils

import (
	"context"
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
type contextKey string

const userIDKey contextKey = "user_id"

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

func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}
