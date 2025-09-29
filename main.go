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
	log.Println("Checking for default admin user...")

	// Check if admin user exists using raw SQL
	var count int64
	checkQuery := `SELECT COUNT(*) FROM users WHERE email = ?`
	err := database.DB.Raw(checkQuery, "admin@example.com").Scan(&count).Error
	
	if err != nil {
		log.Printf("Error checking admin user: %v", err)
		return
	}

	if count == 0 {
		// Hash password admin123
		hashedPassword, err := utils.HashPassword("admin123")
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			return
		}

		// Admin user belum ada, buat yang baru
		insertQuery := `
			INSERT INTO users (username, email, password, role, is_active, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, NOW(), NOW())
		`
		
		result := database.DB.Exec(insertQuery, "admin", "admin@example.com", hashedPassword, "admin", true)
		if result.Error != nil {
			log.Printf("Warning: Could not create default admin user: %v", result.Error)
		} else {
			log.Println("✓ Default admin user created: admin@example.com / admin123")
		}
	} else {
		log.Println("✓ Admin user already exists")
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
		query := `SELECT id, username, email, role, is_active, created_at, updated_at FROM users ORDER BY id`
		err := database.DB.Raw(query).Scan(&users).Error
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to fetch users",
				"details": err.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"users": users,
			"count": len(users),
		})
	})

	// Initialize database connection
	database.ConnectDB()

	// Check database connection health
	if err := database.CheckDatabaseConnection(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Run database migrations (only creates tables if they don't exist)
	database.RunMigrations()

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
