package migration

import (
	"context"
	"log"
	"modul4crud/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RunMongoDBMigrations membuat collections dan indexes untuk MongoDB jika belum ada
func RunMongoDBMigrations() {
	if !database.IsMongoDB() {
		log.Println("Skipping MongoDB migrations (not using MongoDB)")
		return
	}

	log.Println("Running MongoDB database migrations...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// List of collections to create
	collections := []string{
		"users",
		"mahasiswas",
		"alumnis",
		"pekerjaan_alumnis",
	}

	// Get existing collections
	existingCollections, err := database.MongoDB.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		log.Printf("Error listing collections: %v", err)
		return
	}

	// Create collections if they don't exist
	for _, collectionName := range collections {
		exists := false
		for _, existing := range existingCollections {
			if existing == collectionName {
				exists = true
				break
			}
		}

		if !exists {
			log.Printf("Creating collection: %s...", collectionName)
			err := database.MongoDB.CreateCollection(ctx, collectionName)
			if err != nil {
				log.Printf("Error creating collection %s: %v", collectionName, err)
			} else {
				log.Printf("✓ Collection %s created successfully", collectionName)
			}
		} else {
			log.Printf("✓ Collection %s already exists", collectionName)
		}
	}

	// Create indexes (non-blocking - akan tetap lanjut meskipun gagal)
	createMongoDBIndexes(ctx)

	log.Println("MongoDB database migrations completed successfully!")
	log.Println("⚠️  Note: If indexes failed due to disk space, the app will still work but queries may be slower.")
}

// createMongoDBIndexes membuat index yang diperlukan untuk performa MongoDB
func createMongoDBIndexes(ctx context.Context) {
	log.Println("Creating MongoDB indexes...")

	// Indexes untuk users collection
	usersCollection := database.MongoDB.Collection("users")
	createMongoIndex(ctx, usersCollection, "email", true, "idx_users_email")
	createMongoIndex(ctx, usersCollection, "username", true, "idx_users_username")

	// Indexes untuk mahasiswas collection
	mahasiswasCollection := database.MongoDB.Collection("mahasiswas")
	createMongoIndex(ctx, mahasiswasCollection, "nim", true, "idx_mahasiswas_nim")
	createMongoIndex(ctx, mahasiswasCollection, "email", true, "idx_mahasiswas_email")

	// Indexes untuk alumnis collection
	alumnisCollection := database.MongoDB.Collection("alumnis")
	createMongoIndex(ctx, alumnisCollection, "nim", true, "idx_alumnis_nim")
	createMongoIndex(ctx, alumnisCollection, "user_id", false, "idx_alumnis_user_id")

	// Indexes untuk pekerjaan_alumnis collection
	pekerjaanCollection := database.MongoDB.Collection("pekerjaan_alumnis")
	createMongoIndex(ctx, pekerjaanCollection, "alumni_id", false, "idx_pekerjaan_alumni_id")
	createMongoIndex(ctx, pekerjaanCollection, "deleted_at", false, "idx_pekerjaan_deleted_at")

	log.Println("MongoDB indexes creation completed!")
}

// createMongoIndex helper function untuk membuat index di MongoDB
func createMongoIndex(ctx context.Context, collection *mongo.Collection, field string, unique bool, indexName string) {
	// Check if index already exists
	cursor, err := collection.Indexes().List(ctx)
	if err != nil {
		log.Printf("Error listing indexes for %s: %v", collection.Name(), err)
		return
	}

	var indexes []bson.M
	if err = cursor.All(ctx, &indexes); err != nil {
		log.Printf("Error reading indexes for %s: %v", collection.Name(), err)
		return
	}

	// Check if index with this name already exists
	for _, idx := range indexes {
		if name, ok := idx["name"].(string); ok && name == indexName {
			log.Printf("✓ Index %s on %s.%s already exists", indexName, collection.Name(), field)
			return
		}
	}

	// Create the index
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: field, Value: 1}}, // 1 for ascending order
		Options: options.Index().SetUnique(unique).SetName(indexName),
	}

	_, err = collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Printf("Error creating index %s on %s.%s: %v", indexName, collection.Name(), field, err)
	} else {
		log.Printf("✓ Created index %s on %s.%s", indexName, collection.Name(), field)
	}
}

// DropMongoDBCollections fungsi untuk menghapus semua collections (gunakan dengan hati-hati!)
func DropMongoDBCollections() error {
	if !database.IsMongoDB() {
		log.Println("Not using MongoDB, skipping drop collections")
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collections := []string{"users", "mahasiswas", "alumnis", "pekerjaan_alumnis"}

	for _, collectionName := range collections {
		log.Printf("Dropping collection: %s...", collectionName)
		err := database.MongoDB.Collection(collectionName).Drop(ctx)
		if err != nil {
			log.Printf("Error dropping collection %s: %v", collectionName, err)
			return err
		}
		log.Printf("✓ Collection %s dropped successfully", collectionName)
	}

	return nil
}
