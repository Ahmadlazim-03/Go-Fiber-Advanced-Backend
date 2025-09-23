package main

import (
	"log"
	"modul4crud/database"
	"modul4crud/models"
	"modul4crud/repositories"
	"modul4crud/routes"
	"modul4crud/services"
	"modul4crud/utils"

	"github.com/gofiber/fiber/v2"
)

// createDefaultAdmin membuat user admin default jika belum ada
func createDefaultAdmin() {

	var user models.User
	result := database.DB.Where("email = ?", "admin@example.com").First(&user)

	if result.Error != nil {
		// Hash password admin123
		hashedPassword, err := utils.HashPassword("admin123")
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			return
		}

		// Admin user belum ada, buat yang baru
		adminUser := models.User{
			Username: "admin",
			Email:    "admin@example.com",
			Password: hashedPassword,
			Role:     "admin",
			IsActive: true,
		}

		if err := database.DB.Create(&adminUser).Error; err != nil {
			log.Printf("Warning: Could not create default admin user: %v", err)
		} else {
			log.Println("Default admin user created: admin@example.com / admin123")
		}
	} else {
		log.Println("Admin user already exists")
	}
}

func main() {
	app := fiber.New()

	// Static files middleware
	app.Static("/static", "./static")

	// Public routes - tidak perlu autentikasi
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./templates/welcome.html")
	})

	app.Get("/login", func(c *fiber.Ctx) error {
		return c.SendFile("./templates/login.html")
	})

	app.Get("/register", func(c *fiber.Ctx) error {
		return c.SendFile("./templates/register.html")
	})

	// Debug route for testing
	app.Get("/debug", func(c *fiber.Ctx) error {
		return c.SendFile("./templates/debug.html")
	})

	// Debug route untuk melihat semua users
	app.Get("/debug/users", func(c *fiber.Ctx) error {
		var users []models.User
		database.DB.Find(&users)
		return c.JSON(fiber.Map{
			"users": users,
			"count": len(users),
		})
	})

	// Initialize database connection
	database.ConnectDB()

	// Auto migrate tables
	database.DB.AutoMigrate(&models.User{}, &models.Mahasiswa{}, &models.Alumni{}, &models.PekerjaanAlumni{})

	// Create default admin user
	createDefaultAdmin()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(database.DB)
	mahasiswaRepo := repositories.NewMahasiswaRepository(database.DB)
	alumniRepo := repositories.NewAlumniRepository(database.DB)
	pekerjaanRepo := repositories.NewPekerjaanAlumniRepository(database.DB)

	// Initialize services - all with direct repository access
	authService := services.NewAuthService(userRepo)
	mahasiswaService := services.NewMahasiswaService(mahasiswaRepo) // Direct repository
	alumniService := services.NewAlumniService(alumniRepo) // Direct repository
	pekerjaanService := services.NewPekerjaanAlumniService(pekerjaanRepo) // Direct repository

	// Protected dashboard route - perlu autentikasi JWT
	app.Get("/dashboard", func(c *fiber.Ctx) error {
		// Untuk halaman HTML, kita tidak bisa validate JWT di server side
		// Validation akan dilakukan di client side via JavaScript
		return c.SendFile("./templates/index.html")
	})

	// Setup API routes with dependency injection
	routes.SetupRoutes(app, mahasiswaService, alumniService, pekerjaanService, authService)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(app.Listen(":8080"))
}
