package routes

import (
	"modul4crud/middleware"
	"modul4crud/services"

	"github.com/gofiber/fiber/v2"
)

// SetupPekerjaanRoutes configures all pekerjaan alumni (job) related routes
// User: Only GET operations + own jobs
// Admin: Full CRUD operations + soft delete management
func SetupPekerjaanRoutes(api fiber.Router, pekerjaanService *services.PekerjaanAlumniService) {
	pekerjaan := api.Group("/pekerjaan")

	// Public GET routes - accessible by both User & Admin
	pekerjaan.Get("/count", pekerjaanService.GetPekerjaanAlumniCount)                    // Get total count
	pekerjaan.Get("/my-jobs", pekerjaanService.GetPekerjaanByUser)                       // User's own jobs
	pekerjaan.Get("/deleted", pekerjaanService.GetDeletedPekerjaan)                      // Get soft-deleted items
	pekerjaan.Get("/search", pekerjaanService.GetPekerjaanAlumnis)                       // Search endpoint
	pekerjaan.Get("/filter", pekerjaanService.GetPekerjaanAlumnis)                       // Filter endpoint
	pekerjaan.Get("/stats/by-industry", pekerjaanService.GetPekerjaanStatsByIndustry)    // Statistics by industry
	pekerjaan.Get("/stats/by-location", pekerjaanService.GetPekerjaanStatsByLocation)    // Statistics by location
	pekerjaan.Get("/alumni/:alumni_id", pekerjaanService.GetPekerjaanByAlumni)           // Get jobs by alumni ID
	pekerjaan.Get("/", pekerjaanService.GetPekerjaanAlumnis)                             // Get all with pagination
	pekerjaan.Get("/:id", pekerjaanService.GetPekerjaanAlumni)                           // Get by ID

	// Admin-only routes - requires admin role
	pekerjaan.Post("/", middleware.RequireAdmin(), pekerjaanService.CreatePekerjaanAlumni)     // Create new
	pekerjaan.Put("/:id", middleware.RequireAdmin(), pekerjaanService.UpdatePekerjaanAlumni)   // Update existing
	
	// Soft delete operations - admin only
	pekerjaan.Delete("/soft/alumni/:alumni_id", pekerjaanService.SoftDeletePekerjaanByAlumni)  // Soft delete by alumni
	pekerjaan.Delete("/soft/:id", pekerjaanService.SoftDeletePekerjaanAlumni)                  // Soft delete by ID
	pekerjaan.Post("/restore/:id", pekerjaanService.RestorePekerjaanAlumni)                    // Restore soft-deleted
	
	// Hard delete - admin only
	pekerjaan.Delete("/:id", pekerjaanService.DeletePekerjaanAlumni)                           // Permanent delete

	// Company statistics - accessible by all authenticated users
	api.Get("/perusahaan/:nama_perusahaan", pekerjaanService.GetAlumniCountByCompany)          // Alumni count by company
}
