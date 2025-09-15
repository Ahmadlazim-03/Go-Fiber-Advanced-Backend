package services

import (
	"modul4crud/middleware"
	"modul4crud/models"
	"modul4crud/usecases"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AuthService struct {
	authUsecase usecases.AuthUsecase
}

func NewAuthService(authUsecase usecases.AuthUsecase) *AuthService {
	return &AuthService{
		authUsecase: authUsecase,
	}
}

// Register endpoint untuk registrasi user baru
func (s *AuthService) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Data request tidak valid",
		})
	}

	user, err := s.authUsecase.Register(&req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "User berhasil didaftarkan",
		"user":    user,
	})
}

// Login endpoint untuk autentikasi user
func (s *AuthService) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Data request tidak valid",
		})
	}

	response, err := s.authUsecase.Login(&req)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"token":   response.Token,
		"user":    response.User,
	})
}

// GetProfile endpoint untuk mendapatkan profil user yang sedang login
func (s *AuthService) GetProfile(c *fiber.Ctx) error {
	userInfo := middleware.GetUserFromContext(c)

	user, err := s.authUsecase.GetUserByID(userInfo.UserID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "User tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"user": user,
	})
}

// GetUsers endpoint untuk mendapatkan semua user (admin only)
func (s *AuthService) GetUsers(c *fiber.Ctx) error {
	users, err := s.authUsecase.GetAllUsers()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"users": users,
		"total": len(users),
	})
}

// GetUser endpoint untuk mendapatkan user berdasarkan ID (admin only)
func (s *AuthService) GetUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "ID tidak valid",
		})
	}

	user, err := s.authUsecase.GetUserByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "User tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"user": user,
	})
}

// UpdateUser endpoint untuk update user (admin only)
func (s *AuthService) UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "ID tidak valid",
		})
	}

	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Data request tidak valid",
		})
	}

	updatedUser, err := s.authUsecase.UpdateUser(id, &user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "User berhasil diupdate",
		"user":    updatedUser,
	})
}

// DeleteUser endpoint untuk hapus user (admin only)
func (s *AuthService) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "ID tidak valid",
		})
	}

	err = s.authUsecase.DeleteUser(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "User berhasil dihapus",
	})
}

// GetUsersCount endpoint untuk mendapatkan jumlah user (admin only)
func (s *AuthService) GetUsersCount(c *fiber.Ctx) error {
	count, err := s.authUsecase.CountUsers()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"total_users": count,
	})
}

// Logout endpoint (optional - untuk blacklist token di implementasi advanced)
func (s *AuthService) Logout(c *fiber.Ctx) error {
	// Untuk implementasi sederhana, logout hanya mengembalikan pesan
	// Dalam implementasi advanced, token bisa ditambahkan ke blacklist
	return c.JSON(fiber.Map{
		"message": "Logout berhasil",
	})
}
