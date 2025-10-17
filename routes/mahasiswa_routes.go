package routes

import (
	"modul4crud/middleware"
	"modul4crud/services"

	"github.com/gofiber/fiber/v2"
)

// SetupMahasiswaRoutes configures all mahasiswa-related routes
// User: Only GET operations
// Admin: Full CRUD operations
func SetupMahasiswaRoutes(api fiber.Router, mahasiswaService *services.MahasiswaService) {
	mahasiswa := api.Group("/mahasiswa")

	// Public GET routes - accessible by both User & Admin
	mahasiswa.Get("/count", mahasiswaService.GetMahasiswaCount)   // Get total count
	mahasiswa.Get("/search", mahasiswaService.GetMahasiswas)      // Search endpoint
	mahasiswa.Get("/filter", mahasiswaService.GetMahasiswas)      // Filter endpoint
	mahasiswa.Get("/", mahasiswaService.GetMahasiswas)            // Get all with pagination
	mahasiswa.Get("/:id", mahasiswaService.GetMahasiswa)          // Get by ID

	// Admin-only routes - requires admin role
	mahasiswa.Post("/", middleware.RequireAdmin(), mahasiswaService.CreateMahasiswa)      // Create new
	mahasiswa.Put("/:id", middleware.RequireAdmin(), mahasiswaService.UpdateMahasiswa)    // Update existing
	mahasiswa.Delete("/:id", middleware.RequireAdmin(), mahasiswaService.DeleteMahasiswa) // Delete
}
