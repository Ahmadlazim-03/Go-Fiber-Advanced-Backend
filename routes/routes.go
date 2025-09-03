package routes

import (
	"modul4crud/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(
	app *fiber.App,
	mahasiswaController *controllers.MahasiswaController,
	alumniController *controllers.AlumniController,
	pekerjaanController *controllers.PekerjaanAlumniController,
) {
	// Mahasiswa routes
	app.Get("/mahasiswa", mahasiswaController.GetMahasiswas)
	app.Post("/mahasiswa", mahasiswaController.CreateMahasiswa)
	app.Get("/mahasiswa/:id", mahasiswaController.GetMahasiswa)
	app.Put("/mahasiswa/:id", mahasiswaController.UpdateMahasiswa)
	app.Delete("/mahasiswa/:id", mahasiswaController.DeleteMahasiswa)

	// Alumni routes
	app.Get("/alumni", alumniController.GetAlumnis)
	app.Get("/alumni/:id", alumniController.GetAlumni)
	app.Post("/alumni", alumniController.CreateAlumni)
	app.Put("/alumni/:id", alumniController.UpdateAlumni)
	app.Delete("/alumni/:id", alumniController.DeleteAlumni)

	// Pekerjaan Alumni routes
	app.Get("/pekerjaan", pekerjaanController.GetPekerjaanAlumnis)
	app.Get("/pekerjaan/:id", pekerjaanController.GetPekerjaanAlumni)
	app.Get("/pekerjaan/alumni/:alumni_id", pekerjaanController.GetPekerjaanByAlumni)
	app.Post("/pekerjaan", pekerjaanController.CreatePekerjaanAlumni)
	app.Put("/pekerjaan/:id", pekerjaanController.UpdatePekerjaanAlumni)
	app.Delete("/pekerjaan/:id", pekerjaanController.DeletePekerjaanAlumni)
}
