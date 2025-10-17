package routes

import (
	"modul4crud/middleware"
	"modul4crud/services"

	"github.com/gofiber/fiber/v2"
)

// SetupAlumniRoutes configures all alumni-related routes
// User: Only GET operations + own profile
// Admin: Full CRUD operations
func SetupAlumniRoutes(api fiber.Router, alumniService *services.AlumniService) {
	alumni := api.Group("/alumni")

	// Public GET routes - accessible by both User & Admin
	alumni.Get("/count", alumniService.GetAlumniCount)                       // Get total count
	alumni.Get("/my-profile", alumniService.GetAlumniByUser)                 // User's own profile
	alumni.Get("/search", alumniService.GetAlumnis)                          // Search endpoint
	alumni.Get("/filter", alumniService.GetAlumnis)                          // Filter endpoint
	alumni.Get("/stats/by-year", alumniService.GetAlumniStatsByYear)         // Statistics by graduation year
	alumni.Get("/stats/by-jurusan", alumniService.GetAlumniStatsByJurusan)   // Statistics by department
	alumni.Get("/", alumniService.GetAlumnis)                                // Get all with pagination
	alumni.Get("/:id", alumniService.GetAlumni)                              // Get by ID

	// Admin-only routes - requires admin role
	alumni.Post("/", middleware.RequireAdmin(), alumniService.CreateAlumni)      // Create new
	alumni.Put("/:id", middleware.RequireAdmin(), alumniService.UpdateAlumni)    // Update existing
	alumni.Delete("/:id", middleware.RequireAdmin(), alumniService.DeleteAlumni) // Delete
}
