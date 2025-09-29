package repositories

import (
	"fmt"
	"modul4crud/models"

	"gorm.io/gorm"
)

type MahasiswaRepository interface {
	GetAll() ([]models.Mahasiswa, error)
	GetWithPagination(pagination *models.PaginationRequest) ([]models.Mahasiswa, int64, error)
	GetByID(id uint) (*models.Mahasiswa, error)
	Create(mahasiswa *models.Mahasiswa) error
	Update(mahasiswa *models.Mahasiswa) error
	Delete(id uint) error
	Count() (int64, error)
}

type mahasiswaRepository struct {
	db *gorm.DB
}

func NewMahasiswaRepository(db *gorm.DB) MahasiswaRepository {
	return &mahasiswaRepository{db: db}
}

func (r *mahasiswaRepository) GetAll() ([]models.Mahasiswa, error) {
	var mahasiswas []models.Mahasiswa
	
	query := `
		SELECT id, nim, nama, jurusan, angkatan, email, created_at, updated_at
		FROM mahasiswas
		ORDER BY id DESC
	`
	
	err := r.db.Raw(query).Scan(&mahasiswas).Error
	return mahasiswas, err
}

func (r *mahasiswaRepository) GetWithPagination(pagination *models.PaginationRequest) ([]models.Mahasiswa, int64, error) {
	var mahasiswas []models.Mahasiswa
	var total int64
	
	// Set default values
	pagination.SetDefaults()
	pagination.ValidateSortOrder()
	
	// Count query
	countQuery := `SELECT COUNT(*) FROM mahasiswas`
	
	// Search filter
	searchCondition := ""
	searchArgs := []interface{}{}
	if pagination.Search != "" {
		searchPattern := "%" + pagination.Search + "%"
		searchCondition = ` WHERE (
			nim ILIKE ? OR 
			nama ILIKE ? OR 
			jurusan ILIKE ? OR 
			CAST(angkatan AS TEXT) ILIKE ? OR 
			email ILIKE ?
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
		SELECT id, nim, nama, jurusan, angkatan, email, created_at, updated_at
		FROM mahasiswas
	`
	
	// Add search condition to data query
	dataQuery += searchCondition
	
	// Add sorting and pagination
	dataQuery += fmt.Sprintf(" ORDER BY %s %s LIMIT ? OFFSET ?", pagination.SortBy, pagination.SortOrder)
	
	// Prepare arguments for data query
	dataArgs := append(searchArgs, pagination.Limit, pagination.GetOffset())
	
	err = r.db.Raw(dataQuery, dataArgs...).Scan(&mahasiswas).Error
	return mahasiswas, total, err
}

func (r *mahasiswaRepository) GetByID(id uint) (*models.Mahasiswa, error) {
	var mahasiswa models.Mahasiswa
	
	query := `
		SELECT id, nim, nama, jurusan, angkatan, email, created_at, updated_at
		FROM mahasiswas
		WHERE id = ?
	`
	
	err := r.db.Raw(query, id).Scan(&mahasiswa).Error
	if err != nil {
		return nil, err
	}
	return &mahasiswa, nil
}

func (r *mahasiswaRepository) Create(mahasiswa *models.Mahasiswa) error {
	query := `
		INSERT INTO mahasiswas 
		(nim, nama, jurusan, angkatan, email, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	
	return r.db.Raw(query,
		mahasiswa.NIM,
		mahasiswa.Nama,
		mahasiswa.Jurusan,
		mahasiswa.Angkatan,
		mahasiswa.Email,
	).Scan(mahasiswa).Error
}

func (r *mahasiswaRepository) Update(mahasiswa *models.Mahasiswa) error {
	query := `
		UPDATE mahasiswas 
		SET nim = ?, nama = ?, jurusan = ?, angkatan = ?, email = ?, updated_at = NOW()
		WHERE id = ?
		RETURNING updated_at
	`
	
	return r.db.Raw(query,
		mahasiswa.NIM,
		mahasiswa.Nama,
		mahasiswa.Jurusan,
		mahasiswa.Angkatan,
		mahasiswa.Email,
		mahasiswa.ID,
	).Scan(mahasiswa).Error
}

func (r *mahasiswaRepository) Delete(id uint) error {
	query := `DELETE FROM mahasiswas WHERE id = ?`
	result := r.db.Exec(query, id)
	return result.Error
}

func (r *mahasiswaRepository) Count() (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM mahasiswas`
	err := r.db.Raw(query).Scan(&count).Error
	return count, err
}
