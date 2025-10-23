package routes

import (
    "modul4crud/services"
    "github.com/gofiber/fiber/v2"
)

func SetupFileRoutes(app *fiber.App, service services.FileService) {
    api := app.Group("/api")
    files := api.Group("/files")

    files.Post("/upload", service.UploadFile)
    files.Get("/", service.GetAllFiles)
    files.Get("/:id", service.GetFileByID)
    files.Delete("/:id", service.DeleteFile)
}
