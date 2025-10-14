package services

import (
	"modul4crud/models"
	repo "modul4crud/repositories/interface"

	"github.com/gofiber/fiber/v2"
)

type TrashService struct {
	pekerjaanRepo repo.PekerjaanAlumniRepository
}

func NewTrashService(pekerjaanRepo repo.PekerjaanAlumniRepository) *TrashService {
	return &TrashService{
		pekerjaanRepo: pekerjaanRepo,
	}
}

func (s *TrashService) GetAllTrash(c *fiber.Ctx) error {
	userRole, ok := c.Locals("role").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Role tidak ditemukan"})
	}

	var pekerjaanAlumnis []models.PekerjaanAlumni
	var err error

	if userRole == "admin" {
		pekerjaanAlumnis, err = s.pekerjaanRepo.GetDeleted()
	} else {
		userID, ok := c.Locals("user_id").(int)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "User ID tidak ditemukan"})
		}
		pekerjaanAlumnis, err = s.pekerjaanRepo.GetDeletedByUserID(userID)
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch deleted pekerjaan alumni data",
		})
	}

	message := "Data trash berhasil diambil"
	if userRole != "admin" {
		message = "Data trash milik Anda berhasil diambil"
	}

	response := fiber.Map{
		"message": message,
		"data": fiber.Map{
			"pekerjaan_alumni": pekerjaanAlumnis,
		},
		"jumlah data pekerjaan yang di soft": fiber.Map{
			"total_pekerjaan_alumni": len(pekerjaanAlumnis),
		},
	}

	return c.JSON(response)
}
