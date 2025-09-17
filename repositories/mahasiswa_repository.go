package repositories

import (
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
	err := r.db.Find(&mahasiswas).Error
	return mahasiswas, err
}

func (r *mahasiswaRepository) GetWithPagination(pagination *models.PaginationRequest) ([]models.Mahasiswa, int64, error) {
	var mahasiswas []models.Mahasiswa
	var total int64
	
	// Set default values
	pagination.SetDefaults()
	pagination.ValidateSortOrder()
	
	// Base query
	query := r.db.Model(&models.Mahasiswa{})
	
	// Apply search filter if provided
	if pagination.Search != "" {
		searchPattern := "%" + pagination.Search + "%"
		query = query.Where("nim ILIKE ? OR nama ILIKE ? OR jurusan ILIKE ? OR angkatan ILIKE ? OR email ILIKE ?", 
			searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
	}
	
	// Count total records with filters
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply sorting and pagination
	err := query.Order(pagination.SortBy + " " + pagination.SortOrder).
		Limit(pagination.Limit).
		Offset(pagination.GetOffset()).
		Find(&mahasiswas).Error
	
	return mahasiswas, total, err
}

func (r *mahasiswaRepository) GetByID(id uint) (*models.Mahasiswa, error) {
	var mahasiswa models.Mahasiswa
	err := r.db.First(&mahasiswa, id).Error
	if err != nil {
		return nil, err
	}
	return &mahasiswa, nil
}

func (r *mahasiswaRepository) Create(mahasiswa *models.Mahasiswa) error {
	return r.db.Create(mahasiswa).Error
}

func (r *mahasiswaRepository) Update(mahasiswa *models.Mahasiswa) error {
	return r.db.Save(mahasiswa).Error
}

func (r *mahasiswaRepository) Delete(id uint) error {
	return r.db.Delete(&models.Mahasiswa{}, id).Error
}

func (r *mahasiswaRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Mahasiswa{}).Count(&count).Error
	return count, err
}
