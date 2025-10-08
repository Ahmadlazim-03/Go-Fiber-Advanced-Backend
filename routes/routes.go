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
	mahasiswa.Get("/", mahasiswaService.GetMahasiswas)                                    // User & Admin
	mahasiswa.Get("/:id", mahasiswaService.GetMahasiswa)                                  // User & Admin
	mahasiswa.Post("/", middleware.RequireAdmin(), mahasiswaService.CreateMahasiswa)      // Admin only
	mahasiswa.Put("/:id", middleware.RequireAdmin(), mahasiswaService.UpdateMahasiswa)    // Admin only
	mahasiswa.Delete("/:id", middleware.RequireAdmin(), mahasiswaService.DeleteMahasiswa) // Admin only

	// Alumni routes - User: hanya GET, Admin: semua operasi
	alumni := api.Group("/alumni")
	alumni.Get("/count", alumniService.GetAlumniCount)                           // User & Admin
	alumni.Get("/", alumniService.GetAlumnis)                                    // User & Admin
	alumni.Get("/my-profile", alumniService.GetAlumniByUser)                     // User untuk melihat profil sendiri
	alumni.Get("/:id", alumniService.GetAlumni)                                  // User & Admin
	alumni.Post("/", middleware.RequireAdmin(), alumniService.CreateAlumni)      // Admin only
	alumni.Put("/:id", middleware.RequireAdmin(), alumniService.UpdateAlumni)    // Admin only
	alumni.Delete("/:id", middleware.RequireAdmin(), alumniService.DeleteAlumni) // Admin only

	// Pekerjaan Alumni routes - User: hanya GET, Admin: semua operasi
	pekerjaan := api.Group("/pekerjaan")
	pekerjaan.Get("/count", pekerjaanService.GetPekerjaanAlumniCount)
	pekerjaan.Get("/", pekerjaanService.GetPekerjaanAlumnis)
	pekerjaan.Get("/my-jobs", pekerjaanService.GetPekerjaanByUser)
	pekerjaan.Get("/deleted", pekerjaanService.GetDeletedPekerjaan)
	pekerjaan.Get("/:id", pekerjaanService.GetPekerjaanAlumni)
	pekerjaan.Get("/alumni/:alumni_id", pekerjaanService.GetPekerjaanByAlumni)
	pekerjaan.Post("/", middleware.RequireAdmin(), pekerjaanService.CreatePekerjaanAlumni)
	pekerjaan.Put("/:id", middleware.RequireAdmin(), pekerjaanService.UpdatePekerjaanAlumni)
	api.Get("/perusahaan/:nama_perusahaan", pekerjaanService.GetAlumniCountByCompany)



	
	// Endpoint Soal 1
	api.Get("/trash", middleware.RequireAdmin(), trashService.GetAllTrash)
	// Endpoint Soal 2
	pekerjaan.Delete("/soft/:id", pekerjaanService.SoftDeletePekerjaanAlumni)
	pekerjaan.Delete("/soft/alumni/:alumni_id", pekerjaanService.SoftDeletePekerjaanByAlumni)
	pekerjaan.Post("/restore/:id", pekerjaanService.RestorePekerjaanAlumni)
	// Endpoint Soal 3
	pekerjaan.Delete("/:id", middleware.RequireAdmin(), pekerjaanService.DeletePekerjaanAlumni)
}
