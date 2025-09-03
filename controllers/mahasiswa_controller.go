package controllers

import (
	"modul4crud/models"
	"modul4crud/usecases"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type MahasiswaController struct {
	mahasiswaUsecase usecases.MahasiswaUsecase
}

func NewMahasiswaController(mahasiswaUsecase usecases.MahasiswaUsecase) *MahasiswaController {
	return &MahasiswaController{
		mahasiswaUsecase: mahasiswaUsecase,
	}
}

func (ctrl *MahasiswaController) GetMahasiswas(c *fiber.Ctx) error {
	mahasiswas, err := ctrl.mahasiswaUsecase.GetAllMahasiswa()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(mahasiswas)
}

func (ctrl *MahasiswaController) CreateMahasiswa(c *fiber.Ctx) error {
	var req models.CreateMahasiswaRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	mahasiswa, err := ctrl.mahasiswaUsecase.CreateMahasiswa(&req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(mahasiswa)
}

func (ctrl *MahasiswaController) GetMahasiswa(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	mahasiswa, err := ctrl.mahasiswaUsecase.GetMahasiswaByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Mahasiswa not found"})
	}

	return c.JSON(mahasiswa)
}

func (ctrl *MahasiswaController) UpdateMahasiswa(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var req models.UpdateMahasiswaRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	mahasiswa, err := ctrl.mahasiswaUsecase.UpdateMahasiswa(uint(id), &req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(mahasiswa)
}

func (ctrl *MahasiswaController) DeleteMahasiswa(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	err = ctrl.mahasiswaUsecase.DeleteMahasiswa(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(204)
}
