package repositories

import (
	"modul4crud/models"

	"gorm.io/gorm"
)

type PekerjaanAlumniRepository interface {
	GetAll() ([]models.PekerjaanAlumni, error)
	GetWithPagination(pagination *models.PaginationRequest) ([]models.PekerjaanAlumni, int64, error)
	GetByID(id uint) (*models.PekerjaanAlumni, error)
	GetByAlumniID(alumniID uint) ([]models.PekerjaanAlumni, error)
	Create(pekerjaan *models.PekerjaanAlumni) error
	Update(pekerjaan *models.PekerjaanAlumni) error
	Delete(id uint) error
	Count() (int64, error)
	GetAlumniCountByCompany(namaPerusahaan string) (int64, error)
}

type pekerjaanAlumniRepository struct {
	db *gorm.DB
}

func NewPekerjaanAlumniRepository(db *gorm.DB) PekerjaanAlumniRepository {
	return &pekerjaanAlumniRepository{db: db}
}

func (r *pekerjaanAlumniRepository) GetAll() ([]models.PekerjaanAlumni, error) {
	var pekerjaans []models.PekerjaanAlumni
	err := r.db.Preload("Alumni").Find(&pekerjaans).Error
	return pekerjaans, err
}

func (r *pekerjaanAlumniRepository) GetWithPagination(pagination *models.PaginationRequest) ([]models.PekerjaanAlumni, int64, error) {
	var pekerjaans []models.PekerjaanAlumni
	var total int64
	
	// Set default values
	pagination.SetDefaults()
	pagination.ValidateSortOrder()
	
	// Base query
	query := r.db.Model(&models.PekerjaanAlumni{}).Preload("Alumni")
	
	// Apply search filter if provided
	if pagination.Search != "" {
		searchPattern := "%" + pagination.Search + "%"
		query = query.Joins("JOIN alumnis ON pekerjaan_alumnis.alumni_id = alumnis.id").
			Where("pekerjaan_alumnis.posisi ILIKE ? OR pekerjaan_alumnis.nama_perusahaan ILIKE ? OR pekerjaan_alumnis.alamat_perusahaan ILIKE ? OR alumnis.nama ILIKE ?", 
			searchPattern, searchPattern, searchPattern, searchPattern)
	}
	
	// Count total records with filters
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply sorting and pagination
	err := query.Order(pagination.SortBy + " " + pagination.SortOrder).
		Limit(pagination.Limit).
		Offset(pagination.GetOffset()).
		Find(&pekerjaans).Error
	
	return pekerjaans, total, err
}

func (r *pekerjaanAlumniRepository) GetByID(id uint) (*models.PekerjaanAlumni, error) {
	var pekerjaan models.PekerjaanAlumni
	err := r.db.First(&pekerjaan, id).Error
	if err != nil {
		return nil, err
	}
	return &pekerjaan, nil
}

func (r *pekerjaanAlumniRepository) GetByAlumniID(alumniID uint) ([]models.PekerjaanAlumni, error) {
	var pekerjaans []models.PekerjaanAlumni
	err := r.db.Where("alumni_id = ?", alumniID).Find(&pekerjaans).Error
	return pekerjaans, err
}

func (r *pekerjaanAlumniRepository) Create(pekerjaan *models.PekerjaanAlumni) error {
	return r.db.Create(pekerjaan).Error
}

func (r *pekerjaanAlumniRepository) Update(pekerjaan *models.PekerjaanAlumni) error {
	return r.db.Save(pekerjaan).Error
}

func (r *pekerjaanAlumniRepository) Delete(id uint) error {
	return r.db.Delete(&models.PekerjaanAlumni{}, id).Error
}

func (r *pekerjaanAlumniRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.PekerjaanAlumni{}).Count(&count).Error
	return count, err
}

func (r *pekerjaanAlumniRepository) GetAlumniCountByCompany(namaPerusahaan string) (int64, error) {
	var count int64

	err := r.db.Table("pekerjaan_alumnis").
		Where("nama_perusahaan = ?", namaPerusahaan).
		Select("COUNT(DISTINCT alumni_id)").
		Count(&count).Error

	return count, err
}
