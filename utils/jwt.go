package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"modul4crud/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWT Secret Key - dalam production harus di environment variable
var jwtSecret = []byte(getJWTSecret())

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Default secret untuk development - GANTI untuk production!
		return "12345678"
	}
	return secret
}

// generateRandomJTI membuat JTI (JWT ID) yang unik
func generateRandomJTI() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// GenerateJWT membuat token JWT untuk user dengan claims yang unik
func GenerateJWT(user *models.User) (string, error) {
	now := time.Now()
	
	claims := models.JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "CRUD-Go-Fiber-App",                           // Penerbit token
			Subject:   fmt.Sprintf("user_%d", user.ID),               // Subject berdasarkan user ID
			Audience:  []string{"CRUD-Go-Fiber-Users"},               // Audience untuk token
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24)),   // Token berlaku 24 jam
			NotBefore: jwt.NewNumericDate(now),                       // Token valid mulai sekarang
			IssuedAt:  jwt.NewNumericDate(now),                       // Waktu token dibuat
			ID:        generateRandomJTI(),                           // Unique JWT ID
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateJWT memvalidasi token JWT dan mengembalikan claims
func ValidateJWT(tokenString string) (*models.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validasi signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.JWTClaims); ok && token.Valid {
		// Validasi tambahan untuk memastikan token belum expired
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			return nil, fmt.Errorf("token has expired")
		}
		
		// Validasi NotBefore jika ada
		if claims.NotBefore != nil && claims.NotBefore.Time.After(time.Now()) {
			return nil, fmt.Errorf("token not valid yet")
		}
		
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
