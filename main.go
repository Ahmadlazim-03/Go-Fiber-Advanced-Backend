package main

import (
	"log"
	"modul4crud/controllers"
	"modul4crud/database"
	"modul4crud/models"
	"modul4crud/repositories"
	"modul4crud/routes"
	"modul4crud/usecases"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Static files middleware
	app.Static("/static", "./static")

	// Serve HTML template
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./templates/index.html")
	})

	// Initialize database connection
	database.ConnectDB()

	// Auto migrate tables
	database.DB.AutoMigrate(&models.Mahasiswa{}, &models.Alumni{}, &models.PekerjaanAlumni{})

	// Initialize repositories
	mahasiswaRepo := repositories.NewMahasiswaRepository(database.DB)
	alumniRepo := repositories.NewAlumniRepository(database.DB)
	pekerjaanRepo := repositories.NewPekerjaanAlumniRepository(database.DB)

	// Initialize usecases
	mahasiswaUsecase := usecases.NewMahasiswaUsecase(mahasiswaRepo)
	alumniUsecase := usecases.NewAlumniUsecase(alumniRepo)
	pekerjaanUsecase := usecases.NewPekerjaanAlumniUsecase(pekerjaanRepo)

	// Initialize controllers
	mahasiswaController := controllers.NewMahasiswaController(mahasiswaUsecase)
	alumniController := controllers.NewAlumniController(alumniUsecase)
	pekerjaanController := controllers.NewPekerjaanAlumniController(pekerjaanUsecase)

	// Setup routes with dependency injection
	routes.SetupRoutes(app, mahasiswaController, alumniController, pekerjaanController)

	log.Println("Server running on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}
