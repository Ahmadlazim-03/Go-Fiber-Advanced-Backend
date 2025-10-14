package services

import (
	"modul4crud/models"
	repo "modul4crud/repositories/interface"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AlumniService struct {
	alumniRepo repo.AlumniRepository
}

func NewAlumniService(alumniRepo repo.AlumniRepository) *AlumniService {
	return &AlumniService{
		alumniRepo: alumniRepo,
	}
}

func (s *AlumniService) GetAlumnis(c *fiber.Ctx) error {
	// Parse pagination parameters from query
	var pagination models.PaginationRequest
	if err := c.QueryParser(&pagination); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid pagination parameters",
		})
	}

	alumnis, total, err := s.alumniRepo.GetWithPagination(&pagination)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	response := models.NewPaginationResponse(alumnis, &pagination, total)
	return c.JSON(response)
}

// GetAlumnisLegacy endpoint untuk backward compatibility
func (s *AlumniService) GetAlumnisLegacy(c *fiber.Ctx) error {
	alumnis, err := s.alumniRepo.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(alumnis)
}

func (s *AlumniService) CreateAlumni(c *fiber.Ctx) error {
	var req models.CreateAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	alumni := models.Alumni{
		UserID:     req.UserID,
		NIM:        req.NIM,
		Nama:       req.Nama,
		Jurusan:    req.Jurusan,
		Angkatan:   req.Angkatan,
		TahunLulus: req.TahunLulus,
		NoTelepon:  req.NoTelepon,
		Alamat:     req.Alamat,
	}

	err := s.alumniRepo.Create(&alumni)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(alumni)
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
	
	var req models.UpdateAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	
	// Get existing alumni
	alumni, err := s.alumniRepo.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Alumni not found"})
	}
	
	// Update fields (business logic from usecase)
	alumni.Nama = req.Nama
	alumni.Jurusan = req.Jurusan
	alumni.Angkatan = req.Angkatan
	alumni.TahunLulus = req.TahunLulus
	alumni.NoTelepon = req.NoTelepon
	alumni.Alamat = req.Alamat
	
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

// Get alumni by user ID (untuk alumni melihat profil mereka sendiri)
func (s *AlumniService) GetAlumniByUser(c *fiber.Ctx) error {
	// Get user info from middleware locals
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "User ID tidak ditemukan"})
	}

	alumni, err := s.alumniRepo.GetByUserID(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Alumni profile not found"})
	}

	return c.JSON(alumni)
}
