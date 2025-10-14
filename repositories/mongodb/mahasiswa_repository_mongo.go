package mongodb

import (
	"context"
	"fmt"
	"modul4crud/models"
	repo "modul4crud/repositories/interface"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mahasiswaRepositoryMongo struct {
	collection *mongo.Collection
}

func NewMahasiswaRepositoryMongo(db *mongo.Database) repo.MahasiswaRepository {
	return &mahasiswaRepositoryMongo{
		collection: db.Collection("mahasiswas"),
	}
}

func (r *mahasiswaRepositoryMongo) GetAll() ([]models.Mahasiswa, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var mahasiswas []models.Mahasiswa
	if err = cursor.All(ctx, &mahasiswas); err != nil {
		return nil, err
	}

	return mahasiswas, nil
}

func (r *mahasiswaRepositoryMongo) GetWithPagination(pagination *models.PaginationRequest) ([]models.Mahasiswa, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set default values
	pagination.SetDefaults()
	pagination.ValidateSortOrder()

	// Build search filter
	filter := bson.M{}
	if pagination.Search != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"nim": bson.M{"$regex": pagination.Search, "$options": "i"}},
				{"nama": bson.M{"$regex": pagination.Search, "$options": "i"}},
				{"jurusan": bson.M{"$regex": pagination.Search, "$options": "i"}},
				{"email": bson.M{"$regex": pagination.Search, "$options": "i"}},
			},
		}
	}

	// Count total documents
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Build sort order
	sortOrder := 1
	if pagination.SortOrder == "DESC" {
		sortOrder = -1
	}

	// Query options with pagination and sorting
	findOptions := options.Find().
		SetLimit(int64(pagination.Limit)).
		SetSkip(int64(pagination.GetOffset())).
		SetSort(bson.D{{Key: pagination.SortBy, Value: sortOrder}})

	// Execute query
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var mahasiswas []models.Mahasiswa
	if err = cursor.All(ctx, &mahasiswas); err != nil {
		return nil, 0, err
	}

	return mahasiswas, total, nil
}

func (r *mahasiswaRepositoryMongo) GetByID(id uint) (*models.Mahasiswa, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var mahasiswa models.Mahasiswa
	filter := bson.M{"id": id}
	
	err := r.collection.FindOne(ctx, filter).Decode(&mahasiswa)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("mahasiswa not found")
		}
		return nil, err
	}

	return &mahasiswa, nil
}

func (r *mahasiswaRepositoryMongo) Create(mahasiswa *models.Mahasiswa) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Set timestamps
	now := time.Now()
	mahasiswa.CreatedAt = now
	mahasiswa.UpdatedAt = now

	// Get next ID
	nextID, err := r.getNextSequenceID()
	if err != nil {
		return err
	}
	mahasiswa.ID = nextID

	_, err = r.collection.InsertOne(ctx, mahasiswa)
	return err
}

func (r *mahasiswaRepositoryMongo) Update(mahasiswa *models.Mahasiswa) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mahasiswa.UpdatedAt = time.Now()

	filter := bson.M{"id": mahasiswa.ID}
	update := bson.M{
		"$set": bson.M{
			"nim":        mahasiswa.NIM,
			"nama":       mahasiswa.Nama,
			"jurusan":    mahasiswa.Jurusan,
			"angkatan":   mahasiswa.Angkatan,
			"email":      mahasiswa.Email,
			"updated_at": mahasiswa.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("mahasiswa not found")
	}

	return nil
}

func (r *mahasiswaRepositoryMongo) Delete(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("mahasiswa not found")
	}

	return nil
}

func (r *mahasiswaRepositoryMongo) Count() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.collection.CountDocuments(ctx, bson.M{})
}

// Helper function to get next sequence ID
func (r *mahasiswaRepositoryMongo) getNextSequenceID() (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find the document with the highest ID
	findOptions := options.FindOne().SetSort(bson.D{{Key: "id", Value: -1}})
	var result struct {
		ID uint `bson:"id"`
	}

	err := r.collection.FindOne(ctx, bson.M{}, findOptions).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 1, nil // Start from 1 if no documents exist
		}
		return 0, err
	}

	return result.ID + 1, nil
}
