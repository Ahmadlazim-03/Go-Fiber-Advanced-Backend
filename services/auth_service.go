package services

import (
	"fmt"
	"modul4crud/middleware"
	"modul4crud/models"
	repo "modul4crud/repositories/interface"
	"modul4crud/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AuthService struct {
	userRepo repo.UserRepository
}

func NewAuthService(userRepo repo.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
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

	// Business logic moved from usecase
	// Validasi input
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Username, email, dan password wajib diisi",
		})
	}

	// Cek apakah user sudah ada
	existingUser, _ := s.userRepo.GetByUsername(req.Username)
	if existingUser != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Username sudah digunakan",
		})
	}

	existingUser, _ = s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Email sudah digunakan",
		})
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal mengenkripsi password",
		})
	}

	// Create user
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
		IsActive: true,
	}

	// Validasi role - default ke user jika kosong atau tidak valid
	if user.Role != models.RoleAdmin && user.Role != models.RoleUser {
		user.Role = models.RoleUser
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal mendaftarkan user",
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

	// Debug logging
	fmt.Printf("LOGIN DEBUG - Email: %s, Password: %s\n", req.Email, req.Password)

	// Business logic moved from usecase
	// Validasi input
	if req.Email == "" || req.Password == "" {
		fmt.Println("LOGIN DEBUG - Email atau password kosong")
		return c.Status(400).JSON(fiber.Map{
			"error": "Email dan password wajib diisi",
		})
	}

	// Cari user berdasarkan email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil || user == nil {
		fmt.Printf("LOGIN DEBUG - User tidak ditemukan untuk email: %s, Error: %v\n", req.Email, err)
		return c.Status(401).JSON(fiber.Map{
			"error": "Email atau password salah",
		})
	}

	fmt.Printf("LOGIN DEBUG - User ditemukan: %+v\n", user)

	// Cek apakah user aktif
	if !user.IsActive {
		fmt.Println("LOGIN DEBUG - User tidak aktif")
		return c.Status(401).JSON(fiber.Map{
			"error": "Akun tidak aktif",
		})
	}

	// Verify password
	passwordValid := utils.CheckPassword(req.Password, user.Password)
	fmt.Printf("LOGIN DEBUG - Password valid: %v\n", passwordValid)
	
	if !passwordValid {
		fmt.Println("LOGIN DEBUG - Password tidak cocok")
		return c.Status(401).JSON(fiber.Map{
			"error": "Email atau password salah",
		})
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user)
	if err != nil {
		fmt.Println("LOGIN DEBUG - Gagal membuat token")
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal membuat token",
		})
	}

	response := &models.LoginResponse{
		User:  *user,
		Token: token,
	}

	fmt.Println("LOGIN DEBUG - Login berhasil, mengirim response")
	return c.JSON(fiber.Map{
		"message": "Login successful",
		"data": fiber.Map{
			"token": response.Token,
			"user":  response.User,
		},
	})
}

// GetProfile endpoint untuk mendapatkan profil user yang sedang login
func (s *AuthService) GetProfile(c *fiber.Ctx) error {
	userInfo := middleware.GetUserFromContext(c)

	user, err := s.userRepo.GetByID(userInfo.UserID)
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
	// Parse pagination parameters from query
	var pagination models.PaginationRequest
	if err := c.QueryParser(&pagination); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid pagination parameters",
		})
	}

	users, total, err := s.userRepo.GetWithPagination(&pagination)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := models.NewPaginationResponse(users, &pagination, total)
	return c.JSON(response)
}

// GetUsersLegacy endpoint untuk mendapatkan semua user tanpa pagination (backward compatibility)
func (s *AuthService) GetUsersLegacy(c *fiber.Ctx) error {
	users, err := s.userRepo.GetAll()
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

	user, err := s.userRepo.GetByID(id)
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

	var updatedUser models.User
	if err := c.BodyParser(&updatedUser); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Data request tidak valid",
		})
	}

	// Business logic moved from usecase
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "User tidak ditemukan",
		})
	}

	// Update fields yang diizinkan
	if updatedUser.Username != "" {
		user.Username = updatedUser.Username
	}
	if updatedUser.Email != "" {
		user.Email = updatedUser.Email
	}
	if updatedUser.Role != "" {
		// Validasi role - hanya admin dan user
		if updatedUser.Role != models.RoleUser && updatedUser.Role != models.RoleAdmin {
			return c.Status(400).JSON(fiber.Map{
				"error": "Role tidak valid. Hanya 'admin' dan 'user' yang diizinkan",
			})
		}
		user.Role = updatedUser.Role
	}

	user.IsActive = updatedUser.IsActive

	// Hash password baru jika ada
	if updatedUser.Password != "" {
		hashedPassword, err := utils.HashPassword(updatedUser.Password)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Gagal mengenkripsi password",
			})
		}
		user.Password = hashedPassword
	}

	err = s.userRepo.Update(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "User berhasil diupdate",
		"user":    user,
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

	err = s.userRepo.Delete(id)
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
	count, err := s.userRepo.Count()
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
