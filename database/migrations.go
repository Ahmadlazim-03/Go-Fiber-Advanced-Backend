package database

import (
	"log"
	"modul4crud/models"
)

// RunMigrations membuat tabel jika belum ada
func RunMigrations() {
	log.Println("Running database migrations...")

	// Check and create users table
	if !DB.Migrator().HasTable(&models.User{}) {
		log.Println("Creating users table...")
		if err := DB.Migrator().CreateTable(&models.User{}); err != nil {
			log.Printf("Error creating users table: %v", err)
		} else {
			log.Println("✓ Users table created successfully")
		}
	} else {
		log.Println("✓ Users table already exists")
	}

	// Check and create mahasiswas table
	if !DB.Migrator().HasTable(&models.Mahasiswa{}) {
		log.Println("Creating mahasiswas table...")
		if err := DB.Migrator().CreateTable(&models.Mahasiswa{}); err != nil {
			log.Printf("Error creating mahasiswas table: %v", err)
		} else {
			log.Println("✓ Mahasiswas table created successfully")
		}
	} else {
		log.Println("✓ Mahasiswas table already exists")
	}

	// Check and create alumnis table
	if !DB.Migrator().HasTable(&models.Alumni{}) {
		log.Println("Creating alumnis table...")
		if err := DB.Migrator().CreateTable(&models.Alumni{}); err != nil {
			log.Printf("Error creating alumnis table: %v", err)
		} else {
			log.Println("✓ Alumnis table created successfully")
		}
	} else {
		log.Println("✓ Alumnis table already exists")
	}

	// Check and create pekerjaan_alumnis table
	if !DB.Migrator().HasTable(&models.PekerjaanAlumni{}) {
		log.Println("Creating pekerjaan_alumnis table...")
		if err := DB.Migrator().CreateTable(&models.PekerjaanAlumni{}); err != nil {
			log.Printf("Error creating pekerjaan_alumnis table: %v", err)
		} else {
			log.Println("✓ Pekerjaan_alumnis table created successfully")
		}
	} else {
		log.Println("✓ Pekerjaan_alumnis table already exists")
	}

	// Create indexes if they don't exist
	createIndexes()

	log.Println("Database migrations completed successfully!")
}

// createIndexes membuat index yang diperlukan untuk performa
func createIndexes() {
	log.Println("Creating database indexes...")

	// Index untuk users table
	if !DB.Migrator().HasIndex(&models.User{}, "idx_users_email") {
		DB.Migrator().CreateIndex(&models.User{}, "email")
		log.Println("✓ Created index on users.email")
	}

	if !DB.Migrator().HasIndex(&models.User{}, "idx_users_username") {
		DB.Migrator().CreateIndex(&models.User{}, "username")
		log.Println("✓ Created index on users.username")
	}

	// Index untuk mahasiswas table
	if !DB.Migrator().HasIndex(&models.Mahasiswa{}, "idx_mahasiswas_nim") {
		DB.Migrator().CreateIndex(&models.Mahasiswa{}, "nim")
		log.Println("✓ Created index on mahasiswas.nim")
	}

	// Index untuk alumnis table
	if !DB.Migrator().HasIndex(&models.Alumni{}, "idx_alumnis_nim") {
		DB.Migrator().CreateIndex(&models.Alumni{}, "nim")
		log.Println("✓ Created index on alumnis.nim")
	}

	if !DB.Migrator().HasIndex(&models.Alumni{}, "idx_alumnis_user_id") {
		DB.Migrator().CreateIndex(&models.Alumni{}, "user_id")
		log.Println("✓ Created index on alumnis.user_id")
	}

	// Index untuk pekerjaan_alumnis table
	if !DB.Migrator().HasIndex(&models.PekerjaanAlumni{}, "idx_pekerjaan_alumni_id") {
		DB.Migrator().CreateIndex(&models.PekerjaanAlumni{}, "alumni_id")
		log.Println("✓ Created index on pekerjaan_alumnis.alumni_id")
	}

	if !DB.Migrator().HasIndex(&models.PekerjaanAlumni{}, "idx_pekerjaan_deleted_at") {
		DB.Migrator().CreateIndex(&models.PekerjaanAlumni{}, "deleted_at")
		log.Println("✓ Created index on pekerjaan_alumnis.deleted_at")
	}

	log.Println("Database indexes creation completed!")
}

// CheckDatabaseConnection memastikan koneksi database berfungsi
func CheckDatabaseConnection() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Ping(); err != nil {
		return err
	}

	log.Println("✓ Database connection is healthy")
	return nil
}
