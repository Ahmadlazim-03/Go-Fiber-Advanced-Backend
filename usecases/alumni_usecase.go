package usecases

import (
	"modul4crud/models"
	"modul4crud/repositories"
)

type AlumniUsecase interface {
	GetAllAlumni() ([]models.Alumni, error)
	GetAlumniByID(id uint) (*models.Alumni, error)
	CreateAlumni(alumni *models.Alumni) (*models.Alumni, error)
	UpdateAlumni(id uint, alumni *models.Alumni) (*models.Alumni, error)
	DeleteAlumni(id uint) error
	CountAlumni() (int64, error)
}

type alumniUsecase struct {
	alumniRepo repositories.AlumniRepository
}

func NewAlumniUsecase(alumniRepo repositories.AlumniRepository) AlumniUsecase {
	return &alumniUsecase{
		alumniRepo: alumniRepo,
	}
}

func (u *alumniUsecase) GetAllAlumni() ([]models.Alumni, error) {
	return u.alumniRepo.GetAll()
}

func (u *alumniUsecase) GetAlumniByID(id uint) (*models.Alumni, error) {
	return u.alumniRepo.GetByID(id)
}

func (u *alumniUsecase) CreateAlumni(alumni *models.Alumni) (*models.Alumni, error) {
	err := u.alumniRepo.Create(alumni)
	if err != nil {
		return nil, err
	}
	return alumni, nil
}

func (u *alumniUsecase) UpdateAlumni(id uint, updatedAlumni *models.Alumni) (*models.Alumni, error) {
	alumni, err := u.alumniRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	alumni.Nama = updatedAlumni.Nama
	alumni.Jurusan = updatedAlumni.Jurusan
	alumni.Angkatan = updatedAlumni.Angkatan
	alumni.TahunLulus = updatedAlumni.TahunLulus
	alumni.Email = updatedAlumni.Email
	alumni.NoTelepon = updatedAlumni.NoTelepon
	alumni.Alamat = updatedAlumni.Alamat

	err = u.alumniRepo.Update(alumni)
	if err != nil {
		return nil, err
	}

	return alumni, nil
}

func (u *alumniUsecase) DeleteAlumni(id uint) error {
	return u.alumniRepo.Delete(id)
}

func (u *alumniUsecase) CountAlumni() (int64, error) {
	return u.alumniRepo.Count()
}
