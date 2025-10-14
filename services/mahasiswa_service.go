package services

import (
	"modul4crud/models"
	repo "modul4crud/repositories/interface"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type MahasiswaService struct {
	mahasiswaRepo repo.MahasiswaRepository
}

func NewMahasiswaService(mahasiswaRepo repo.MahasiswaRepository) *MahasiswaService {
	return &MahasiswaService{
		mahasiswaRepo: mahasiswaRepo,
	}
}

func (s *MahasiswaService) GetMahasiswas(c *fiber.Ctx) error {
	// Parse pagination parameters from query
	var pagination models.PaginationRequest
	if err := c.QueryParser(&pagination); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid pagination parameters",
		})
	}

	mahasiswas, total, err := s.mahasiswaRepo.GetWithPagination(&pagination)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	response := models.NewPaginationResponse(mahasiswas, &pagination, total)
	return c.JSON(response)
}

// GetMahasiswasLegacy endpoint untuk backward compatibility
func (s *MahasiswaService) GetMahasiswasLegacy(c *fiber.Ctx) error {
	mahasiswas, err := s.mahasiswaRepo.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(mahasiswas)
}

func (s *MahasiswaService) CreateMahasiswa(c *fiber.Ctx) error {
	var req models.CreateMahasiswaRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Convert request to model (business logic from usecase)
	mahasiswa := &models.Mahasiswa{
		NIM:      req.NIM,
		Nama:     req.Nama,
		Email:    req.Email,
		Jurusan:  req.Jurusan,
		Angkatan: req.Angkatan,
	}

	err := s.mahasiswaRepo.Create(mahasiswa)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(mahasiswa)
}

func (s *MahasiswaService) GetMahasiswa(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	mahasiswa, err := s.mahasiswaRepo.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Mahasiswa not found"})
	}

	return c.JSON(mahasiswa)
}

func (s *MahasiswaService) UpdateMahasiswa(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var req models.UpdateMahasiswaRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Get existing mahasiswa
	mahasiswa, err := s.mahasiswaRepo.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Mahasiswa not found"})
	}

	// Update fields (business logic from usecase)
	mahasiswa.Nama = req.Nama
	mahasiswa.Email = req.Email
	mahasiswa.Jurusan = req.Jurusan
	mahasiswa.Angkatan = req.Angkatan

	err = s.mahasiswaRepo.Update(mahasiswa)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(mahasiswa)
}

func (s *MahasiswaService) DeleteMahasiswa(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	err = s.mahasiswaRepo.Delete(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(204)
}

func (s *MahasiswaService) GetMahasiswaCount(c *fiber.Ctx) error {
	count, err := s.mahasiswaRepo.Count()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"total_mahasiswa": count,
	})
}
