package services

import (
	"modul4crud/repositories"

	"github.com/gofiber/fiber/v2"
)

type TrashService struct {
	pekerjaanRepo repositories.PekerjaanAlumniRepository
}

func NewTrashService(pekerjaanRepo repositories.PekerjaanAlumniRepository) *TrashService {
	return &TrashService{
		pekerjaanRepo: pekerjaanRepo,
	}
}

func (s *TrashService) GetAllTrash(c *fiber.Ctx) error {
	pekerjaanAlumnis, err := s.pekerjaanRepo.GetDeleted()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch deleted pekerjaan alumni data",
		})
	}
	response := fiber.Map{
		"message": "Data trash berhasil diambil",
		"data": fiber.Map{
			"pekerjaan_alumni": pekerjaanAlumnis,
		},
		"summary": fiber.Map{
			"total_pekerjaan_alumni": len(pekerjaanAlumnis),
		},
	}

	return c.JSON(response)
}
