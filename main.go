package main

import (
	"log"
	"modul4crud/database"
	"modul4crud/database/migration"
	"modul4crud/models"
	repo "modul4crud/repositories/interface"
	"modul4crud/repositories/mongodb"
	"modul4crud/repositories/pocketbase"
	"modul4crud/repositories/postgres"
	"modul4crud/routes"
	"modul4crud/services"
	"modul4crud/utils"

	"github.com/gofiber/fiber/v2"
)

// createDefaultAdmin membuat user admin default jika belum ada
func createDefaultAdmin(userRepo repo.UserRepository) {
	log.Println("Checking for default admin user...")

	// Check if admin user exists
	existingUser, err := userRepo.GetByEmail("admin@example.com")
	if err != nil {
		log.Printf("Error checking admin user: %v", err)
		return
	}

	if existingUser == nil {
		// Hash password admin123
		hashedPassword, err := utils.HashPassword("admin123")
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			return
		}

		// Admin user belum ada, buat yang baru
		adminUser := &models.User{
			Username: "admin",
			Email:    "admin@example.com",
			Password: hashedPassword,
			Role:     "admin",
			IsActive: true,
		}

		err = userRepo.Create(adminUser)
		if err != nil {
			log.Printf("Warning: Could not create default admin user: %v", err)
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
				"error":   "Failed to fetch users",
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

	// Run database migrations
	migration.RunMigrations()

	// Initialize repositories based on database type
	var userRepo repo.UserRepository
	var mahasiswaRepo repo.MahasiswaRepository
	var alumniRepo repo.AlumniRepository
	var pekerjaanRepo repo.PekerjaanAlumniRepository
	var fileRepo repo.FileRepository

	if database.IsPostgres() {
		userRepo = postgre.NewUserRepository(database.DB)
		mahasiswaRepo = postgre.NewMahasiswaRepository(database.DB)
		alumniRepo = postgre.NewAlumniRepository(database.DB)
		pekerjaanRepo = postgre.NewPekerjaanAlumniRepository(database.DB)
		// TODO: Tambahkan fileRepo Postgres jika ada
	} else if database.IsMongoDB() {
		userRepo = mongodb.NewUserRepositoryMongo(database.MongoDB)
		mahasiswaRepo = mongodb.NewMahasiswaRepositoryMongo(database.MongoDB)
		alumniRepo = mongodb.NewAlumniRepositoryMongo(database.MongoDB)
		pekerjaanRepo = mongodb.NewPekerjaanAlumniRepositoryMongo(database.MongoDB)
		fileRepo = mongodb.NewFileRepository(database.MongoDB)
	} else if database.IsPocketBase() {
		userRepo = pocketbase.NewUserRepository(database.PocketBaseURL)
		mahasiswaRepo = pocketbase.NewMahasiswaRepository(database.PocketBaseURL)
		alumniRepo = pocketbase.NewAlumniRepository(database.PocketBaseURL)
		pekerjaanRepo = pocketbase.NewPekerjaanAlumniRepository(database.PocketBaseURL)
		// TODO: Tambahkan fileRepo PocketBase jika ada
		log.Println("✓ All PocketBase repositories initialized successfully")
	}

	// Create default admin user
	createDefaultAdmin(userRepo)

	// Initialize services - all with direct repository access
	authService := services.NewAuthService(userRepo)
	mahasiswaService := services.NewMahasiswaService(mahasiswaRepo)       // Direct repository
	alumniService := services.NewAlumniService(alumniRepo)                // Direct repository
	pekerjaanService := services.NewPekerjaanAlumniService(pekerjaanRepo) // Direct repository
	trashService := services.NewTrashService(pekerjaanRepo)               // Trash service untuk data soft deleted
	fileService := services.NewFileService(fileRepo, "./uploads")        // Path upload file

	// Protected dashboard route - perlu autentikasi JWT
	app.Get("/dashboard", func(c *fiber.Ctx) error {
		// Untuk halaman HTML, kita tidak bisa validate JWT di server side
		// Validation akan dilakukan di client side via JavaScript
		return c.SendFile("./templates/index.html")
	})

	// Setup API routes with dependency injection
	routes.SetupRoutes(app, mahasiswaService, alumniService, pekerjaanService, authService, trashService, fileService)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(app.Listen(":8080"))
}
