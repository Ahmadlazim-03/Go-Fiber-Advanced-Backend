package usecases

import (
	"modul4crud/models"
	"modul4crud/repositories"
)

type MahasiswaUsecase interface {
	GetAllMahasiswa() ([]models.Mahasiswa, error)
	GetMahasiswaByID(id uint) (*models.Mahasiswa, error)
	CreateMahasiswa(req *models.CreateMahasiswaRequest) (*models.Mahasiswa, error)
	UpdateMahasiswa(id uint, req *models.UpdateMahasiswaRequest) (*models.Mahasiswa, error)
	DeleteMahasiswa(id uint) error
}

type mahasiswaUsecase struct {
	mahasiswaRepo repositories.MahasiswaRepository
}

func NewMahasiswaUsecase(mahasiswaRepo repositories.MahasiswaRepository) MahasiswaUsecase {
	return &mahasiswaUsecase{
		mahasiswaRepo: mahasiswaRepo,
	}
}

func (u *mahasiswaUsecase) GetAllMahasiswa() ([]models.Mahasiswa, error) {
	return u.mahasiswaRepo.GetAll()
}

func (u *mahasiswaUsecase) GetMahasiswaByID(id uint) (*models.Mahasiswa, error) {
	return u.mahasiswaRepo.GetByID(id)
}

func (u *mahasiswaUsecase) CreateMahasiswa(req *models.CreateMahasiswaRequest) (*models.Mahasiswa, error) {
	mahasiswa := &models.Mahasiswa{
		NIM:      req.NIM,
		Nama:     req.Nama,
		Jurusan:  req.Jurusan,
		Angkatan: req.Angkatan,
		Email:    req.Email,
	}

	err := u.mahasiswaRepo.Create(mahasiswa)
	if err != nil {
		return nil, err
	}

	return mahasiswa, nil
}

func (u *mahasiswaUsecase) UpdateMahasiswa(id uint, req *models.UpdateMahasiswaRequest) (*models.Mahasiswa, error) {
	mahasiswa, err := u.mahasiswaRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	mahasiswa.Nama = req.Nama
	mahasiswa.Jurusan = req.Jurusan
	mahasiswa.Angkatan = req.Angkatan
	mahasiswa.Email = req.Email

	err = u.mahasiswaRepo.Update(mahasiswa)
	if err != nil {
		return nil, err
	}

	return mahasiswa, nil
}

func (u *mahasiswaUsecase) DeleteMahasiswa(id uint) error {
	return u.mahasiswaRepo.Delete(id)
}
