package usecases

import (
	"modul4crud/models"
	"modul4crud/repositories"
)

type PekerjaanAlumniUsecase interface {
	GetAllPekerjaanAlumni() ([]models.PekerjaanAlumni, error)
	GetPekerjaanAlumniByID(id uint) (*models.PekerjaanAlumni, error)
	GetPekerjaanByAlumniID(alumniID uint) ([]models.PekerjaanAlumni, error)
	CreatePekerjaanAlumni(pekerjaan *models.PekerjaanAlumni) (*models.PekerjaanAlumni, error)
	UpdatePekerjaanAlumni(id uint, pekerjaan *models.PekerjaanAlumni) (*models.PekerjaanAlumni, error)
	DeletePekerjaanAlumni(id uint) error
}

type pekerjaanAlumniUsecase struct {
	pekerjaanRepo repositories.PekerjaanAlumniRepository
}

func NewPekerjaanAlumniUsecase(pekerjaanRepo repositories.PekerjaanAlumniRepository) PekerjaanAlumniUsecase {
	return &pekerjaanAlumniUsecase{
		pekerjaanRepo: pekerjaanRepo,
	}
}

func (u *pekerjaanAlumniUsecase) GetAllPekerjaanAlumni() ([]models.PekerjaanAlumni, error) {
	return u.pekerjaanRepo.GetAll()
}

func (u *pekerjaanAlumniUsecase) GetPekerjaanAlumniByID(id uint) (*models.PekerjaanAlumni, error) {
	return u.pekerjaanRepo.GetByID(id)
}

func (u *pekerjaanAlumniUsecase) GetPekerjaanByAlumniID(alumniID uint) ([]models.PekerjaanAlumni, error) {
	return u.pekerjaanRepo.GetByAlumniID(alumniID)
}

func (u *pekerjaanAlumniUsecase) CreatePekerjaanAlumni(pekerjaan *models.PekerjaanAlumni) (*models.PekerjaanAlumni, error) {
	err := u.pekerjaanRepo.Create(pekerjaan)
	if err != nil {
		return nil, err
	}
	return pekerjaan, nil
}

func (u *pekerjaanAlumniUsecase) UpdatePekerjaanAlumni(id uint, updatedPekerjaan *models.PekerjaanAlumni) (*models.PekerjaanAlumni, error) {
	pekerjaan, err := u.pekerjaanRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	pekerjaan.AlumniID = updatedPekerjaan.AlumniID
	pekerjaan.NamaPerusahaan = updatedPekerjaan.NamaPerusahaan
	pekerjaan.PosisiJabatan = updatedPekerjaan.PosisiJabatan
	pekerjaan.BidangIndustri = updatedPekerjaan.BidangIndustri
	pekerjaan.LokasiKerja = updatedPekerjaan.LokasiKerja
	pekerjaan.GajiRange = updatedPekerjaan.GajiRange
	pekerjaan.TanggalMulaiKerja = updatedPekerjaan.TanggalMulaiKerja
	pekerjaan.TanggalSelesaiKerja = updatedPekerjaan.TanggalSelesaiKerja
	pekerjaan.StatusPekerjaan = updatedPekerjaan.StatusPekerjaan
	pekerjaan.DeskripsiPekerjaan = updatedPekerjaan.DeskripsiPekerjaan

	err = u.pekerjaanRepo.Update(pekerjaan)
	if err != nil {
		return nil, err
	}

	return pekerjaan, nil
}

func (u *pekerjaanAlumniUsecase) DeletePekerjaanAlumni(id uint) error {
	return u.pekerjaanRepo.Delete(id)
}
