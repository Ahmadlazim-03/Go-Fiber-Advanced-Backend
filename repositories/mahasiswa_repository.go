package repositories

import (
	"modul4crud/models"

	"gorm.io/gorm"
)

type MahasiswaRepository interface {
	GetAll() ([]models.Mahasiswa, error)
	GetByID(id uint) (*models.Mahasiswa, error)
	Create(mahasiswa *models.Mahasiswa) error
	Update(mahasiswa *models.Mahasiswa) error
	Delete(id uint) error
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
