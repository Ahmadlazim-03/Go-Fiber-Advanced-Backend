package repositories

import (
	"fmt"
	"modul4crud/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetAll() ([]models.User, error)
	GetWithPagination(pagination *models.PaginationRequest) ([]models.User, int64, error)
	GetByID(id int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id int) error
	Count() (int64, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAll() ([]models.User, error) {
	var users []models.User
	
	query := `
		SELECT id, username, email, role, is_active, created_at, updated_at
		FROM users
		ORDER BY id DESC
	`
	
	err := r.db.Raw(query).Scan(&users).Error
	return users, err
}

func (r *userRepository) GetWithPagination(pagination *models.PaginationRequest) ([]models.User, int64, error) {
	var users []models.User
	var total int64
	
	// Set default values
	pagination.SetDefaults()
	pagination.ValidateSortOrder()
	
	// Count query
	countQuery := `SELECT COUNT(*) FROM users`
	
	// Search filter
	searchCondition := ""
	searchArgs := []interface{}{}
	if pagination.Search != "" {
		searchPattern := "%" + pagination.Search + "%"
		searchCondition = ` WHERE (
			username ILIKE ? OR 
			email ILIKE ? OR 
			role ILIKE ?
		)`
		searchArgs = []interface{}{searchPattern, searchPattern, searchPattern}
	}
	
	// Execute count query
	err := r.db.Raw(countQuery+searchCondition, searchArgs...).Scan(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	// Data query
	dataQuery := `
		SELECT id, username, email, role, is_active, created_at, updated_at
		FROM users
	`
	
	// Add search condition to data query
	dataQuery += searchCondition
	
	// Add sorting and pagination
	dataQuery += fmt.Sprintf(" ORDER BY %s %s LIMIT ? OFFSET ?", pagination.SortBy, pagination.SortOrder)
	
	// Prepare arguments for data query
	dataArgs := append(searchArgs, pagination.Limit, pagination.GetOffset())
	
	err = r.db.Raw(dataQuery, dataArgs...).Scan(&users).Error
	return users, total, err
}

func (r *userRepository) GetByID(id int) (*models.User, error) {
	var user models.User
	
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil when no record found
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil when no record found
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil when no record found
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users 
		(username, email, password, role, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	
	return r.db.Raw(query,
		user.Username,
		user.Email,
		user.Password,
		user.Role,
		user.IsActive,
	).Scan(user).Error
}

func (r *userRepository) Update(user *models.User) error {
	query := `
		UPDATE users 
		SET username = ?, email = ?, password = ?, role = ?, is_active = ?, updated_at = NOW()
		WHERE id = ?
		RETURNING updated_at
	`
	
	return r.db.Raw(query,
		user.Username,
		user.Email,
		user.Password,
		user.Role,
		user.IsActive,
		user.ID,
	).Scan(user).Error
}

func (r *userRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = ?`
	result := r.db.Exec(query, id)
	return result.Error
}

func (r *userRepository) Count() (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM users`
	err := r.db.Raw(query).Scan(&count).Error
	return count, err
}
