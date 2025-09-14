package services

import (
	"modul4crud/models"
	"modul4crud/usecases"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type MahasiswaService struct {
	mahasiswaUsecase usecases.MahasiswaUsecase
}

func NewMahasiswaService(mahasiswaUsecase usecases.MahasiswaUsecase) *MahasiswaService {
	return &MahasiswaService{
		mahasiswaUsecase: mahasiswaUsecase,
	}
}

func (s *MahasiswaService) GetMahasiswas(c *fiber.Ctx) error {
	mahasiswas, err := s.mahasiswaUsecase.GetAllMahasiswa()
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

	mahasiswa, err := s.mahasiswaUsecase.CreateMahasiswa(&req)
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

	mahasiswa, err := s.mahasiswaUsecase.GetMahasiswaByID(uint(id))
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

	mahasiswa, err := s.mahasiswaUsecase.UpdateMahasiswa(uint(id), &req)
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

	err = s.mahasiswaUsecase.DeleteMahasiswa(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(204)
}

func (s *MahasiswaService) GetMahasiswaCount(c *fiber.Ctx) error {
	count, err := s.mahasiswaUsecase.CountMahasiswa()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"total_mahasiswa": count,
	})
}
