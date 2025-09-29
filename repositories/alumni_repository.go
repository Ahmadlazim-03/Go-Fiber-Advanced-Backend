package repositories

import (
	"fmt"
	"modul4crud/models"

	"gorm.io/gorm"
)

type AlumniRepository interface {
	GetAll() ([]models.Alumni, error)
	GetWithPagination(pagination *models.PaginationRequest) ([]models.Alumni, int64, error)
	GetByID(id uint) (*models.Alumni, error)
	GetByUserID(userID int) (*models.Alumni, error)
	Create(alumni *models.Alumni) error
	Update(alumni *models.Alumni) error
	Delete(id uint) error
	Count() (int64, error)
}

type alumniRepository struct {
	db *gorm.DB
}

func NewAlumniRepository(db *gorm.DB) AlumniRepository {
	return &alumniRepository{db: db}
}

func (r *alumniRepository) GetAll() ([]models.Alumni, error) {
	var alumnis []models.Alumni
	
	query := `
		SELECT 
			a.id, a.user_id, a.nim, a.nama, a.jurusan, 
			a.angkatan, a.tahun_lulus, a.no_telepon, a.alamat, 
			a.created_at, a.updated_at,
			u.id as "User__id", u.username as "User__username", 
			u.email as "User__email", u.role as "User__role", 
			u.is_active as "User__is_active", u.created_at as "User__created_at", 
			u.updated_at as "User__updated_at"
		FROM alumnis a
		LEFT JOIN users u ON a.user_id = u.id
		ORDER BY a.id DESC
	`
	
	err := r.db.Raw(query).Scan(&alumnis).Error
	return alumnis, err
}

func (r *alumniRepository) GetWithPagination(pagination *models.PaginationRequest) ([]models.Alumni, int64, error) {
	var alumnis []models.Alumni
	var total int64
	
	// Set default values
	pagination.SetDefaults()
	pagination.ValidateSortOrder()
	
	// Count query
	countQuery := `
		SELECT COUNT(*) 
		FROM alumnis a
		LEFT JOIN users u ON a.user_id = u.id
	`
	
	// Search filter
	searchCondition := ""
	searchArgs := []interface{}{}
	if pagination.Search != "" {
		searchPattern := "%" + pagination.Search + "%"
		searchCondition = ` WHERE (
			a.nim ILIKE ? OR 
			a.nama ILIKE ? OR 
			a.jurusan ILIKE ? OR 
			CAST(a.tahun_lulus AS TEXT) ILIKE ? OR 
			u.email ILIKE ?
		)`
		searchArgs = []interface{}{searchPattern, searchPattern, searchPattern, searchPattern, searchPattern}
	}
	
	// Execute count query
	err := r.db.Raw(countQuery+searchCondition, searchArgs...).Scan(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	// Data query
	dataQuery := `
		SELECT 
			a.id, a.user_id, a.nim, a.nama, a.jurusan, 
			a.angkatan, a.tahun_lulus, a.no_telepon, a.alamat, 
			a.created_at, a.updated_at,
			u.id as "User__id", u.username as "User__username", 
			u.email as "User__email", u.role as "User__role", 
			u.is_active as "User__is_active", u.created_at as "User__created_at", 
			u.updated_at as "User__updated_at"
		FROM alumnis a
		LEFT JOIN users u ON a.user_id = u.id
	`
	
	// Add search condition to data query
	dataQuery += searchCondition
	
	// Add sorting and pagination
	dataQuery += fmt.Sprintf(" ORDER BY a.%s %s LIMIT ? OFFSET ?", pagination.SortBy, pagination.SortOrder)
	
	// Prepare arguments for data query
	dataArgs := append(searchArgs, pagination.Limit, pagination.GetOffset())
	
	err = r.db.Raw(dataQuery, dataArgs...).Scan(&alumnis).Error
	return alumnis, total, err
}

func (r *alumniRepository) GetByID(id uint) (*models.Alumni, error) {
	var alumni models.Alumni
	
	query := `
		SELECT 
			a.id, a.user_id, a.nim, a.nama, a.jurusan, 
			a.angkatan, a.tahun_lulus, a.no_telepon, a.alamat, 
			a.created_at, a.updated_at,
			u.id as "User__id", u.username as "User__username", 
			u.email as "User__email", u.role as "User__role", 
			u.is_active as "User__is_active", u.created_at as "User__created_at", 
			u.updated_at as "User__updated_at"
		FROM alumnis a
		LEFT JOIN users u ON a.user_id = u.id
		WHERE a.id = ?
	`
	
	err := r.db.Raw(query, id).Scan(&alumni).Error
	if err != nil {
		return nil, err
	}
	return &alumni, nil
}

func (r *alumniRepository) GetByUserID(userID int) (*models.Alumni, error) {
	var alumni models.Alumni
	
	query := `
		SELECT 
			a.id, a.user_id, a.nim, a.nama, a.jurusan, 
			a.angkatan, a.tahun_lulus, a.no_telepon, a.alamat, 
			a.created_at, a.updated_at,
			u.id as "User__id", u.username as "User__username", 
			u.email as "User__email", u.role as "User__role", 
			u.is_active as "User__is_active", u.created_at as "User__created_at", 
			u.updated_at as "User__updated_at"
		FROM alumnis a
		LEFT JOIN users u ON a.user_id = u.id
		WHERE a.user_id = ?
	`
	
	err := r.db.Raw(query, userID).Scan(&alumni).Error
	if err != nil {
		return nil, err
	}
	return &alumni, nil
}

func (r *alumniRepository) Create(alumni *models.Alumni) error {
	query := `
		INSERT INTO alumnis 
		(user_id, nim, nama, jurusan, angkatan, tahun_lulus, no_telepon, alamat, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	
	return r.db.Raw(query,
		alumni.UserID,
		alumni.NIM,
		alumni.Nama,
		alumni.Jurusan,
		alumni.Angkatan,
		alumni.TahunLulus,
		alumni.NoTelepon,
		alumni.Alamat,
	).Scan(alumni).Error
}

func (r *alumniRepository) Update(alumni *models.Alumni) error {
	query := `
		UPDATE alumnis 
		SET nim = ?, nama = ?, jurusan = ?, angkatan = ?, 
		    tahun_lulus = ?, no_telepon = ?, alamat = ?, updated_at = NOW()
		WHERE id = ?
		RETURNING updated_at
	`
	
	return r.db.Raw(query,
		alumni.NIM,
		alumni.Nama,
		alumni.Jurusan,
		alumni.Angkatan,
		alumni.TahunLulus,
		alumni.NoTelepon,
		alumni.Alamat,
		alumni.ID,
	).Scan(alumni).Error
}

func (r *alumniRepository) Delete(id uint) error {
	query := `DELETE FROM alumnis WHERE id = ?`
	result := r.db.Exec(query, id)
	return result.Error
}

func (r *alumniRepository) Count() (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM alumnis`
	err := r.db.Raw(query).Scan(&count).Error
	return count, err
}
