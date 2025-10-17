package routes

import (
	"modul4crud/middleware"
	"modul4crud/services"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(
	app *fiber.App,
	mahasiswaService *services.MahasiswaService,
	alumniService *services.AlumniService,
	pekerjaanService *services.PekerjaanAlumniService,
	authService *services.AuthService,
	trashService *services.TrashService,
) {
	// Global variable untuk status API
	var isAPIActive = true

	// Public routes - tidak perlu autentikasi
	auth := app.Group("/auth")
	auth.Post("/register", authService.Register)
	auth.Post("/login", authService.Login)

	// API public routes (alias untuk compatibility)
	app.Post("/api/register", authService.Register)
	app.Post("/api/login", authService.Login)

	// Protected routes - perlu autentikasi
	api := app.Group("/api", middleware.ValidateJWT())

	// API Status routes - admin only
	api.Post("/status", middleware.RequireAdmin(), func(c *fiber.Ctx) error {
		type StatusRequest struct {
			Active bool `json:"active"`
		}
		var req StatusRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}
		isAPIActive = req.Active
		return c.JSON(fiber.Map{"active": isAPIActive})
	})

	api.Get("/status", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"active": isAPIActive})
	})

	// Profile routes - user bisa akses profil sendiri
	api.Get("/profile", authService.GetProfile)
	api.Post("/logout", authService.Logout)

	// User management routes - admin only
	users := api.Group("/users", middleware.RequireAdmin())
	users.Get("/", authService.GetUsers)
	users.Get("/count", authService.GetUsersCount)
	users.Get("/:id", authService.GetUser)
	users.Put("/:id", authService.UpdateUser)
	users.Delete("/:id", authService.DeleteUser)

	// Mahasiswa routes - User: hanya GET, Admin: semua operasi
	mahasiswa := api.Group("/mahasiswa")
	mahasiswa.Get("/count", mahasiswaService.GetMahasiswaCount)                           // User & Admin
	mahasiswa.Get("/search", mahasiswaService.GetMahasiswas)                              // Search endpoint
	mahasiswa.Get("/filter", mahasiswaService.GetMahasiswas)                              // Filter endpoint
	mahasiswa.Get("/", mahasiswaService.GetMahasiswas)                                    // User & Admin
	mahasiswa.Get("/:id", mahasiswaService.GetMahasiswa)                                  // User & Admin
	mahasiswa.Post("/", middleware.RequireAdmin(), mahasiswaService.CreateMahasiswa)      // Admin only
	mahasiswa.Put("/:id", middleware.RequireAdmin(), mahasiswaService.UpdateMahasiswa)    // Admin only
	mahasiswa.Delete("/:id", middleware.RequireAdmin(), mahasiswaService.DeleteMahasiswa) // Admin only

	// Alumni routes - User: hanya GET, Admin: semua operasi
	alumni := api.Group("/alumni")
	alumni.Get("/count", alumniService.GetAlumniCount)                           // User & Admin
	alumni.Get("/my-profile", alumniService.GetAlumniByUser)                     // User untuk melihat profil sendiri
	alumni.Get("/search", alumniService.GetAlumnis)                              // Search endpoint
	alumni.Get("/filter", alumniService.GetAlumnis)                              // Filter endpoint
	alumni.Get("/stats/by-year", alumniService.GetAlumniStatsByYear)             // Statistics by year
	alumni.Get("/stats/by-jurusan", alumniService.GetAlumniStatsByJurusan)       // Statistics by jurusan
	alumni.Get("/", alumniService.GetAlumnis)                                    // User & Admin
	alumni.Get("/:id", alumniService.GetAlumni)                                  // User & Admin
	alumni.Post("/", middleware.RequireAdmin(), alumniService.CreateAlumni)      // Admin only
	alumni.Put("/:id", middleware.RequireAdmin(), alumniService.UpdateAlumni)    // Admin only
	alumni.Delete("/:id", middleware.RequireAdmin(), alumniService.DeleteAlumni) // Admin only

	// Pekerjaan Alumni routes - User: hanya GET, Admin: semua operasi
	pekerjaan := api.Group("/pekerjaan")
	pekerjaan.Get("/count", pekerjaanService.GetPekerjaanAlumniCount)
	pekerjaan.Get("/my-jobs", pekerjaanService.GetPekerjaanByUser)
	pekerjaan.Get("/deleted", pekerjaanService.GetDeletedPekerjaan)
	pekerjaan.Get("/search", pekerjaanService.GetPekerjaanAlumnis)                        // Search endpoint
	pekerjaan.Get("/filter", pekerjaanService.GetPekerjaanAlumnis)                        // Filter endpoint
	pekerjaan.Get("/stats/by-industry", pekerjaanService.GetPekerjaanStatsByIndustry)     // Statistics by industry
	pekerjaan.Get("/stats/by-location", pekerjaanService.GetPekerjaanStatsByLocation)     // Statistics by location
	pekerjaan.Get("/alumni/:alumni_id", pekerjaanService.GetPekerjaanByAlumni)
	pekerjaan.Get("/", pekerjaanService.GetPekerjaanAlumnis)
	pekerjaan.Get("/:id", pekerjaanService.GetPekerjaanAlumni)
	pekerjaan.Post("/", middleware.RequireAdmin(), pekerjaanService.CreatePekerjaanAlumni)
	pekerjaan.Put("/:id", middleware.RequireAdmin(), pekerjaanService.UpdatePekerjaanAlumni)
	pekerjaan.Delete("/soft/alumni/:alumni_id", pekerjaanService.SoftDeletePekerjaanByAlumni)
	pekerjaan.Delete("/soft/:id", pekerjaanService.SoftDeletePekerjaanAlumni)             // Soft delete
	pekerjaan.Post("/restore/:id", pekerjaanService.RestorePekerjaanAlumni)               // Restore
	pekerjaan.Delete("/:id", pekerjaanService.DeletePekerjaanAlumni)                      // Hard delete
	
	api.Get("/perusahaan/:nama_perusahaan", pekerjaanService.GetAlumniCountByCompany)

	// Trash routes
	trash := api.Group("/trash", middleware.RequireAdmin())
	trash.Get("/pekerjaan", pekerjaanService.GetDeletedPekerjaan)                         // Get trashed pekerjaan
	trash.Post("/pekerjaan/:id/restore", pekerjaanService.RestorePekerjaanAlumni)         // Restore pekerjaan
	trash.Delete("/pekerjaan/:id", pekerjaanService.DeletePekerjaanAlumni)                // Permanent delete
	trash.Get("/", trashService.GetAllTrash)                                              // Get all trash
}
