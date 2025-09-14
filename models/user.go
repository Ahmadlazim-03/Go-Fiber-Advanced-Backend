package models

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

// User model untuk autentikasi
type User struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"type:varchar(50);unique;not null" json:"username"`
	Email     string    `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"` // Hide password in JSON
	Role      string    `gorm:"type:varchar(20);default:'user'" json:"role"`
	IsActive  bool      `gorm:"default:true" json:"-"` // Hide in JSON
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"-"` // Hide in JSON
}

// Request struct untuk registrasi
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role,omitempty"` // Optional, default to 'user'
}

// Request struct untuk login
type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Response struct untuk login
type LoginResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

// JWT Claims struct
type JWTClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Role constants - Hanya Admin dan User
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)
