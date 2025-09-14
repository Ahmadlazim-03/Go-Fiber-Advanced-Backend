package repositories

import (
	"modul4crud/models"

	"gorm.io/gorm"
)

type AlumniRepository interface {
	GetAll() ([]models.Alumni, error)
	GetByID(id uint) (*models.Alumni, error)
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
	err := r.db.Find(&alumnis).Error
	return alumnis, err
}

func (r *alumniRepository) GetByID(id uint) (*models.Alumni, error) {
	var alumni models.Alumni
	err := r.db.First(&alumni, id).Error
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
