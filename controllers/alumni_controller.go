package controllers

import (
	"modul4crud/models"
	"modul4crud/usecases"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AlumniController struct {
	alumniUsecase usecases.AlumniUsecase
}

func NewAlumniController(alumniUsecase usecases.AlumniUsecase) *AlumniController {
	return &AlumniController{
		alumniUsecase: alumniUsecase,
	}
}

func (ctrl *AlumniController) GetAlumnis(c *fiber.Ctx) error {
	alumnis, err := ctrl.alumniUsecase.GetAllAlumni()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(alumnis)
}

func (ctrl *AlumniController) CreateAlumni(c *fiber.Ctx) error {
	var alumni models.Alumni
	if err := c.BodyParser(&alumni); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := ctrl.alumniUsecase.CreateAlumni(&alumni)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

func (ctrl *AlumniController) GetAlumni(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	alumni, err := ctrl.alumniUsecase.GetAlumniByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Alumni not found"})
	}

	return c.JSON(alumni)
}

func (ctrl *AlumniController) UpdateAlumni(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var alumni models.Alumni
	if err := c.BodyParser(&alumni); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := ctrl.alumniUsecase.UpdateAlumni(uint(id), &alumni)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

func (ctrl *AlumniController) DeleteAlumni(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	err = ctrl.alumniUsecase.DeleteAlumni(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(204)
}
