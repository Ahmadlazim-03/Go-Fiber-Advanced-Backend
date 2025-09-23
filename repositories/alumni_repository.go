package repositories

import (
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
	err := r.db.Preload("User").Preload("Pekerjaan").Find(&alumnis).Error
	return alumnis, err
}

func (r *alumniRepository) GetWithPagination(pagination *models.PaginationRequest) ([]models.Alumni, int64, error) {
	var alumnis []models.Alumni
	var total int64
	
	// Set default values
	pagination.SetDefaults()
	pagination.ValidateSortOrder()
	
	// Base query
	query := r.db.Model(&models.Alumni{}).Preload("User").Preload("Pekerjaan")
	
	// Apply search filter if provided
	if pagination.Search != "" {
		searchPattern := "%" + pagination.Search + "%"
		query = query.Joins("JOIN users ON alumnis.user_id = users.id").
			Where("nim ILIKE ? OR nama ILIKE ? OR jurusan ILIKE ? OR CAST(tahun_lulus AS TEXT) ILIKE ? OR users.email ILIKE ?", 
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
		Find(&alumnis).Error
	
	return alumnis, total, err
}

func (r *alumniRepository) GetByID(id uint) (*models.Alumni, error) {
	var alumni models.Alumni
	err := r.db.Preload("User").Preload("Pekerjaan").First(&alumni, id).Error
	if err != nil {
		return nil, err
	}
	return &alumni, nil
}

func (r *alumniRepository) GetByUserID(userID int) (*models.Alumni, error) {
	var alumni models.Alumni
	err := r.db.Preload("User").Preload("Pekerjaan").Where("user_id = ?", userID).First(&alumni).Error
	if err != nil {
		return nil, err
	}
	return &alumni, nil
}

func (r *alumniRepository) Create(alumni *models.Alumni) error {
	return r.db.Create(alumni).Error
}

func (r *alumniRepository) Update(alumni *models.Alumni) error {
	return r.db.Save(alumni).Error
}

func (r *alumniRepository) Delete(id uint) error {
	return r.db.Delete(&models.Alumni{}, id).Error
}

func (r *alumniRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Alumni{}).Count(&count).Error
	return count, err
}
