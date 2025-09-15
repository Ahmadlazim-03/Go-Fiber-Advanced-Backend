package services

import (
	"modul4crud/models"
	"modul4crud/usecases"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AlumniService struct {
	alumniUsecase usecases.AlumniUsecase
}

func NewAlumniService(alumniUsecase usecases.AlumniUsecase) *AlumniService {
	return &AlumniService{
		alumniUsecase: alumniUsecase,
	}
}

func (s *AlumniService) GetAlumnis(c *fiber.Ctx) error {
	alumnis, err := s.alumniUsecase.GetAllAlumni()
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

	result, err := s.alumniUsecase.CreateAlumni(&alumni)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

func (s *AlumniService) GetAlumni(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	alumni, err := s.alumniUsecase.GetAlumniByID(uint(id))
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

	var alumni models.Alumni
	if err := c.BodyParser(&alumni); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := s.alumniUsecase.UpdateAlumni(uint(id), &alumni)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

func (s *AlumniService) DeleteAlumni(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	err = s.alumniUsecase.DeleteAlumni(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(204)
}

func (s *AlumniService) CountAlumni(c *fiber.Ctx) error {
	count, err := s.alumniUsecase.CountAlumni()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"count":   count,
		"message": "Total jumlah alumni",
	})
}

func (s *AlumniService) GetAlumniCount(c *fiber.Ctx) error {
	count, err := s.alumniUsecase.CountAlumni()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"total_alumni": count,
	})
}
