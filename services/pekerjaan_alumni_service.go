package services

import (
	"modul4crud/models"
	"modul4crud/repositories"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PekerjaanAlumniService struct {
	pekerjaanRepo repositories.PekerjaanAlumniRepository
}

func NewPekerjaanAlumniService(pekerjaanRepo repositories.PekerjaanAlumniRepository) *PekerjaanAlumniService {
	return &PekerjaanAlumniService{
		pekerjaanRepo: pekerjaanRepo,
	}
}

func (s *PekerjaanAlumniService) GetPekerjaanAlumnis(c *fiber.Ctx) error {
	// Parse pagination parameters from query
	var pagination models.PaginationRequest
	if err := c.QueryParser(&pagination); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid pagination parameters",
		})
	}

	pekerjaans, total, err := s.pekerjaanRepo.GetWithPagination(&pagination)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	response := models.NewPaginationResponse(pekerjaans, &pagination, total)
	return c.JSON(response)
}

// GetPekerjaanAlumnisLegacy endpoint untuk backward compatibility
func (s *PekerjaanAlumniService) GetPekerjaanAlumnisLegacy(c *fiber.Ctx) error {
	pekerjaans, err := s.pekerjaanRepo.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(pekerjaans)
}

func (s *PekerjaanAlumniService) CreatePekerjaanAlumni(c *fiber.Ctx) error {
	var pekerjaan models.PekerjaanAlumni
	if err := c.BodyParser(&pekerjaan); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	err := s.pekerjaanRepo.Create(&pekerjaan)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(pekerjaan)
}

func (s *PekerjaanAlumniService) GetPekerjaanAlumni(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	pekerjaan, err := s.pekerjaanRepo.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Pekerjaan not found"})
	}

	return c.JSON(pekerjaan)
}

func (s *PekerjaanAlumniService) GetPekerjaanByAlumni(c *fiber.Ctx) error {
	alumniID, err := strconv.ParseUint(c.Params("alumni_id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Alumni ID"})
	}

	pekerjaans, err := s.pekerjaanRepo.GetByAlumniID(uint(alumniID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(pekerjaans)
}

func (s *PekerjaanAlumniService) UpdatePekerjaanAlumni(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var updatedPekerjaan models.PekerjaanAlumni
	if err := c.BodyParser(&updatedPekerjaan); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Get existing pekerjaan
	pekerjaan, err := s.pekerjaanRepo.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Pekerjaan not found"})
	}

	// Update fields (business logic from usecase)
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

	err = s.pekerjaanRepo.Update(pekerjaan)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(pekerjaan)
}

func (s *PekerjaanAlumniService) DeletePekerjaanAlumni(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	err = s.pekerjaanRepo.Delete(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(204)
}

func (s *PekerjaanAlumniService) GetPekerjaanAlumniCount(c *fiber.Ctx) error {
	count, err := s.pekerjaanRepo.Count()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"total_pekerjaan_alumni": count,
	})
}

func (s *PekerjaanAlumniService) GetAlumniCountByCompany(c *fiber.Ctx) error {
	namaPerusahaan := c.Params("nama_perusahaan")
	if namaPerusahaan == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Nama perusahaan tidak boleh kosong"})
	}

	count, err := s.pekerjaanRepo.GetAlumniCountByCompany(namaPerusahaan)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"nama_perusahaan": namaPerusahaan,
		"jumlah_alumni":   count,
		"message":         "Data jumlah alumni di perusahaan berhasil diambil",
	})
}
