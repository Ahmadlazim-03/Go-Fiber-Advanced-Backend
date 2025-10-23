package routes

import (
	"modul4crud/middleware"
	"modul4crud/services"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all application routes
// Routes are organized by domain/module in separate files:
// - auth_routes.go: Authentication & user management
// - mahasiswa_routes.go: Student management
// - alumni_routes.go: Alumni management
// - pekerjaan_routes.go: Job/employment management
// - trash_routes.go: Soft delete/recycle bin management
func SetupRoutes(
	app *fiber.App,
	mahasiswaService *services.MahasiswaService,
	alumniService *services.AlumniService,
	pekerjaanService *services.PekerjaanAlumniService,
	authService *services.AuthService,
	trashService *services.TrashService,
	fileService services.FileService,
) {
	// Global variable for API status
	var isAPIActive = true

	// ========================================
	// PUBLIC ROUTES - No authentication required
	// These MUST be defined BEFORE the protected API group
	// ========================================
	
	// Public authentication routes
	app.Post("/api/register", authService.Register)
	app.Post("/api/login", authService.Login)
	
	auth := app.Group("/auth")
	auth.Post("/register", authService.Register)
	auth.Post("/login", authService.Login)

	// ========================================
	// PROTECTED API GROUP - JWT authentication required
	// All routes under /api/* (except register/login above) need JWT
	// ========================================
	api := app.Group("/api", middleware.ValidateJWT())

	// API Status routes - admin only
	// Allows admin to enable/disable API temporarily
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

	// ========================================
	// SETUP MODULAR ROUTES
	// Each function handles its own domain/module
	// ========================================
	
	// Authentication & user management (protected routes only)
	// Public routes already defined above
	api.Get("/profile", authService.GetProfile)
	api.Post("/logout", authService.Logout)
	users := api.Group("/users", middleware.RequireAdmin())
	users.Get("/", authService.GetUsers)
	users.Get("/count", authService.GetUsersCount)
	users.Get("/:id", authService.GetUser)
	users.Put("/:id", authService.UpdateUser)
	users.Delete("/:id", authService.DeleteUser)
	
	SetupMahasiswaRoutes(api, mahasiswaService)          // Student management
	SetupAlumniRoutes(api, alumniService)                // Alumni management
	SetupPekerjaanRoutes(api, pekerjaanService)          // Job/employment management
	SetupTrashRoutes(api, pekerjaanService, trashService) // Trash/recycle bin
	SetupFileRoutes(api, fileService)                    // File management
}
