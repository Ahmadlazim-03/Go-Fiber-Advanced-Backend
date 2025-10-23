package routes

import (
    "modul4crud/services"
    "github.com/gofiber/fiber/v2"
)

func SetupFileRoutes(router fiber.Router, service services.FileService) {
    files := router.Group("/files")

    files.Post("/upload", service.UploadFile)
    files.Get("/", service.GetAllFiles)
    files.Get("/:id", service.GetFileByID)
    files.Delete("/:id", service.DeleteFile)
}
