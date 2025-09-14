package services

import (
	"modul4crud/models"
	"modul4crud/usecases"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PekerjaanAlumniService struct {
	pekerjaanUsecase usecases.PekerjaanAlumniUsecase
}

func NewPekerjaanAlumniService(pekerjaanUsecase usecases.PekerjaanAlumniUsecase) *PekerjaanAlumniService {
	return &PekerjaanAlumniService{
		pekerjaanUsecase: pekerjaanUsecase,
	}
}

func (s *PekerjaanAlumniService) GetPekerjaanAlumnis(c *fiber.Ctx) error {
	pekerjaans, err := s.pekerjaanUsecase.GetAllPekerjaanAlumni()
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

	result, err := s.pekerjaanUsecase.CreatePekerjaanAlumni(&pekerjaan)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

func (s *PekerjaanAlumniService) GetPekerjaanAlumni(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	pekerjaan, err := s.pekerjaanUsecase.GetPekerjaanAlumniByID(uint(id))
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

	pekerjaans, err := s.pekerjaanUsecase.GetPekerjaanByAlumniID(uint(alumniID))
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

	var pekerjaan models.PekerjaanAlumni
	if err := c.BodyParser(&pekerjaan); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := s.pekerjaanUsecase.UpdatePekerjaanAlumni(uint(id), &pekerjaan)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

func (s *PekerjaanAlumniService) DeletePekerjaanAlumni(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	err = s.pekerjaanUsecase.DeletePekerjaanAlumni(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(204)
}

func (s *PekerjaanAlumniService) GetPekerjaanAlumniCount(c *fiber.Ctx) error {
	count, err := s.pekerjaanUsecase.CountPekerjaanAlumni()
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

	count, err := s.pekerjaanUsecase.GetAlumniCountByCompany(namaPerusahaan)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"nama_perusahaan": namaPerusahaan,
		"jumlah_alumni": count,
		"message": "Data jumlah alumni di perusahaan berhasil diambil",
	})
}
