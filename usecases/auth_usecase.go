package usecases

import (
	"errors"
	"modul4crud/models"
	"modul4crud/repositories"
	"modul4crud/utils"

	"gorm.io/gorm"
)

type AuthUsecase interface {
	Register(req *models.RegisterRequest) (*models.User, error)
	Login(req *models.LoginRequest) (*models.LoginResponse, error)
	GetUserByID(id int) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	UpdateUser(id int, user *models.User) (*models.User, error)
	DeleteUser(id int) error
	CountUsers() (int64, error)
}

type authUsecase struct {
	userRepo repositories.UserRepository
}

func NewAuthUsecase(userRepo repositories.UserRepository) AuthUsecase {
	return &authUsecase{
		userRepo: userRepo,
	}
}

func (u *authUsecase) Register(req *models.RegisterRequest) (*models.User, error) {
	// Check apakah email sudah digunakan
	existingUser, err := u.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email sudah digunakan")
	}

	// Check apakah username sudah digunakan
	existingUser, err = u.userRepo.GetByUsername(req.Username)
	if err == nil && existingUser != nil {
		return nil, errors.New("username sudah digunakan")
	}

	// Set default role jika tidak disediakan
	role := req.Role
	if role == "" {
		role = models.RoleUser
	}

	// Validasi role - hanya admin dan user
	if role != models.RoleUser && role != models.RoleAdmin {
		return nil, errors.New("role tidak valid. Hanya 'admin' dan 'user' yang diizinkan")
	}

	// Buat user baru
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     role,
		IsActive: true,
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, errors.New("gagal mengenkripsi password")
	}
	user.Password = hashedPassword

	// Simpan ke database
	if err := u.userRepo.Create(user); err != nil {
		return nil, errors.New("gagal membuat user")
	}

	return user, nil
}

func (u *authUsecase) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	// Cari user berdasarkan email
	user, err := u.userRepo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("email atau password salah")
		}
		return nil, errors.New("gagal mengakses database")
	}

	// Check apakah user aktif
	if !user.IsActive {
		return nil, errors.New("akun tidak aktif")
	}

	// Validasi password
	if !utils.CheckPassword(user.Password, req.Password) {
		return nil, errors.New("email atau password salah")
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user)
	if err != nil {
		return nil, errors.New("gagal membuat token")
	}

	response := &models.LoginResponse{
		User:  *user,
		Token: token,
	}

	return response, nil
}

func (u *authUsecase) GetUserByID(id int) (*models.User, error) {
	return u.userRepo.GetByID(id)
}

func (u *authUsecase) GetAllUsers() ([]models.User, error) {
	return u.userRepo.GetAll()
}

func (u *authUsecase) UpdateUser(id int, updatedUser *models.User) (*models.User, error) {
	user, err := u.userRepo.GetByID(id)
	if err != nil {
		return nil, err
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
			return nil, errors.New("role tidak valid. Hanya 'admin' dan 'user' yang diizinkan")
		}
		user.Role = updatedUser.Role
	}

	user.IsActive = updatedUser.IsActive

	// Hash password baru jika ada
	if updatedUser.Password != "" {
		hashedPassword, err := utils.HashPassword(updatedUser.Password)
		if err != nil {
			return nil, errors.New("gagal mengenkripsi password")
		}
		user.Password = hashedPassword
	}

	err = u.userRepo.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *authUsecase) DeleteUser(id int) error {
	return u.userRepo.Delete(id)
}

func (u *authUsecase) CountUsers() (int64, error) {
	return u.userRepo.Count()
}
