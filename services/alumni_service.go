package services

import (
	"modul4crud/models"
	"modul4crud/repositories"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AlumniService struct {
	alumniRepo repositories.AlumniRepository
}

func NewAlumniService(alumniRepo repositories.AlumniRepository) *AlumniService {
	return &AlumniService{
		alumniRepo: alumniRepo,
	}
}

func (s *AlumniService) GetAlumnis(c *fiber.Ctx) error {
	alumnis, err := s.alumniRepo.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(alumnis)
}

func (s *AlumniService) CreateAlumni(c *fiber.Ctx) error {
	var alumni models.Alumni
	if err := c.BodyParser(&alumni); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	err := s.alumniRepo.Create(&alumni)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(alumni)
}

func (s *AlumniService) GetAlumni(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}
	alumni, err := s.alumniRepo.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Alumni not found"})
	}
	return c.JSON(alumni)
}

func (s *AlumniService) UpdateAlumni(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}
	var updatedAlumni models.Alumni
	if err := c.BodyParser(&updatedAlumni); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	
	// Get existing alumni
	alumni, err := s.alumniRepo.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Alumni not found"})
	}
	
	// Update fields (business logic from usecase)
	alumni.Nama = updatedAlumni.Nama
	alumni.Jurusan = updatedAlumni.Jurusan
	alumni.Angkatan = updatedAlumni.Angkatan
	alumni.TahunLulus = updatedAlumni.TahunLulus
	alumni.Email = updatedAlumni.Email
	alumni.NoTelepon = updatedAlumni.NoTelepon
	alumni.Alamat = updatedAlumni.Alamat
	
	err = s.alumniRepo.Update(alumni)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(alumni)
}

func (s *AlumniService) DeleteAlumni(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}
	err = s.alumniRepo.Delete(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}

func (s *AlumniService) CountAlumni(c *fiber.Ctx) error {
	count, err := s.alumniRepo.Count()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"count":   count,
		"message": "Total jumlah alumni",
	})
}

func (s *AlumniService) GetAlumniCount(c *fiber.Ctx) error {
	count, err := s.alumniRepo.Count()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"total_alumni": count,
	})
}
