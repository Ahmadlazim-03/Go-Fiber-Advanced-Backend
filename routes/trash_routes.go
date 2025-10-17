package routes

import (
	"modul4crud/middleware"
	"modul4crud/services"

	"github.com/gofiber/fiber/v2"
)

// SetupTrashRoutes configures all trash/recycle bin related routes
// All trash routes are admin-only for managing soft-deleted items
func SetupTrashRoutes(api fiber.Router, pekerjaanService *services.PekerjaanAlumniService, trashService *services.TrashService) {
	trash := api.Group("/trash", middleware.RequireAdmin())

	// Pekerjaan trash management
	trash.Get("/pekerjaan", pekerjaanService.GetDeletedPekerjaan)                  // Get all trashed pekerjaan
	trash.Post("/pekerjaan/:id/restore", pekerjaanService.RestorePekerjaanAlumni)  // Restore specific pekerjaan
	trash.Delete("/pekerjaan/:id", pekerjaanService.DeletePekerjaanAlumni)         // Permanent delete pekerjaan

	// General trash operations
	trash.Get("/", trashService.GetAllTrash)                                       // Get all trashed items (all types)
}
