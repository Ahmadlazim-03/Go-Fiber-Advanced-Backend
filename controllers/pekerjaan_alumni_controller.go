package controllers

import (
	"modul4crud/models"
	"modul4crud/usecases"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PekerjaanAlumniController struct {
	pekerjaanUsecase usecases.PekerjaanAlumniUsecase
}

func NewPekerjaanAlumniController(pekerjaanUsecase usecases.PekerjaanAlumniUsecase) *PekerjaanAlumniController {
	return &PekerjaanAlumniController{
		pekerjaanUsecase: pekerjaanUsecase,
	}
}

func (ctrl *PekerjaanAlumniController) GetPekerjaanAlumnis(c *fiber.Ctx) error {
	pekerjaans, err := ctrl.pekerjaanUsecase.GetAllPekerjaanAlumni()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(pekerjaans)
}

func (ctrl *PekerjaanAlumniController) CreatePekerjaanAlumni(c *fiber.Ctx) error {
	var pekerjaan models.PekerjaanAlumni
	if err := c.BodyParser(&pekerjaan); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := ctrl.pekerjaanUsecase.CreatePekerjaanAlumni(&pekerjaan)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

func (ctrl *PekerjaanAlumniController) GetPekerjaanAlumni(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	pekerjaan, err := ctrl.pekerjaanUsecase.GetPekerjaanAlumniByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Pekerjaan not found"})
	}

	return c.JSON(pekerjaan)
}

func (ctrl *PekerjaanAlumniController) GetPekerjaanByAlumni(c *fiber.Ctx) error {
	alumniID, err := strconv.ParseUint(c.Params("alumni_id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Alumni ID"})
	}

	pekerjaans, err := ctrl.pekerjaanUsecase.GetPekerjaanByAlumniID(uint(alumniID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(pekerjaans)
}

func (ctrl *PekerjaanAlumniController) UpdatePekerjaanAlumni(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var pekerjaan models.PekerjaanAlumni
	if err := c.BodyParser(&pekerjaan); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := ctrl.pekerjaanUsecase.UpdatePekerjaanAlumni(uint(id), &pekerjaan)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

func (ctrl *PekerjaanAlumniController) DeletePekerjaanAlumni(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	err = ctrl.pekerjaanUsecase.DeletePekerjaanAlumni(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(204)
}
