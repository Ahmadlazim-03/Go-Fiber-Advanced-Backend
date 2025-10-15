package migration

import (
	"log"
	"modul4crud/database"
)

// RunMigrations adalah wrapper function yang menjalankan migrations sesuai tipe database
func RunMigrations() {
	log.Printf("Database type: %s", database.GetDBType())

	switch database.GetDBType() {
	case "postgres":
		RunPostgresMigrations()
	case "mongodb":
		RunMongoDBMigrations()
	case "pocketbase":
		RunPocketBaseMigrations()
	default:
		log.Printf("Unknown database type: %s", database.GetDBType())
	}
}
