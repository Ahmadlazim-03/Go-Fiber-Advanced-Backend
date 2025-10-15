package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/pocketbase/pocketbase"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB             *gorm.DB
	MongoDB        *mongo.Database
	MongoClient    *mongo.Client
	PocketBaseApp  *pocketbase.PocketBase
	PocketBaseURL  string
	PocketBaseAuth string // Token untuk auth
	DBType         string
)

func ConnectDB() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	DBType = os.Getenv("DB_TYPE")
	if DBType == "" {
		DBType = "postgres" // Default to postgres
	}

	switch DBType {
	case "postgres":
		connectPostgres()
	case "mongodb":
		connectMongoDB()
	case "pocketbase":
		connectPocketBase()
	default:
		log.Fatalf("Unknown database type: %s. Use 'postgres', 'mongodb', or 'pocketbase'", DBType)
	}
}

func connectPostgres() {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		log.Fatal("POSTGRES_DSN environment variable is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}

	DB = db
	log.Println("✓ Connected to PostgreSQL successfully")
}

func connectMongoDB() {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MONGODB_URI environment variable is not set")
	}

	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = "railway" // Default database name
	}

	log.Printf("Connecting to MongoDB at: %s", uri)
	log.Printf("Database name: %s", dbName)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Configure client options with better settings for Railway
	clientOptions := options.Client().
		ApplyURI(uri).
		SetConnectTimeout(30 * time.Second).
		SetServerSelectionTimeout(30 * time.Second).
		SetSocketTimeout(30 * time.Second).
		SetMaxPoolSize(50).
		SetMinPoolSize(10)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Ping the database to verify connection
	log.Println("Pinging MongoDB to verify connection...")
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	MongoClient = client
	MongoDB = client.Database(dbName)
	log.Printf("✓ Connected to MongoDB successfully (Database: %s)\n", dbName)
}

func connectPocketBase() {
	url := os.Getenv("POCKETBASE_URL")
	if url == "" {
		log.Fatal("POCKETBASE_URL environment variable is not set")
	}

	PocketBaseURL = url
	
	log.Printf("Connecting to PocketBase at: %s", url)
	log.Println("✓ PocketBase URL configured successfully")
	log.Println("Note: PocketBase uses HTTP API - no persistent connection needed")
}

func DisconnectMongoDB() error {
	if MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return MongoClient.Disconnect(ctx)
	}
	return nil
}

func GetDBType() string {
	return DBType
}

func IsMongoDB() bool {
	return DBType == "mongodb"
}

func IsPostgres() bool {
	return DBType == "postgres"
}

func IsPocketBase() bool {
	return DBType == "pocketbase"
}

func CheckDatabaseConnection() error {
	if IsPostgres() && DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return fmt.Errorf("failed to get database instance: %v", err)
		}
		return sqlDB.Ping()
	} else if IsMongoDB() && MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return MongoClient.Ping(ctx, nil)
	} else if IsPocketBase() && PocketBaseURL != "" {
		// PocketBase menggunakan HTTP API, tidak ada koneksi persistent
		log.Println("PocketBase uses HTTP API - connection check skipped")
		return nil
	}
	return fmt.Errorf("no database connection available")
}
