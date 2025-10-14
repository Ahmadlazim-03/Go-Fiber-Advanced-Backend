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

type alumniRepositoryMongo struct {
	collection     *mongo.Collection
	userCollection *mongo.Collection
}

func NewAlumniRepositoryMongo(db *mongo.Database) repo.AlumniRepository {
	return &alumniRepositoryMongo{
		collection:     db.Collection("alumnis"),
		userCollection: db.Collection("users"),
	}
}

func (r *alumniRepositoryMongo) GetAll() ([]models.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Use aggregation pipeline to join with users collection
	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "user_id"},
			{Key: "foreignField", Value: "id"},
			{Key: "as", Value: "user"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$user"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "id", Value: -1}}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var alumnis []models.Alumni
	if err = cursor.All(ctx, &alumnis); err != nil {
		return nil, err
	}

	return alumnis, nil
}

func (r *alumniRepositoryMongo) GetWithPagination(pagination *models.PaginationRequest) ([]models.Alumni, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set default values
	pagination.SetDefaults()
	pagination.ValidateSortOrder()

	// Build search filter
	matchStage := bson.D{}
	if pagination.Search != "" {
		matchStage = bson.D{
			{Key: "$match", Value: bson.M{
				"$or": []bson.M{
					{"nim": bson.M{"$regex": pagination.Search, "$options": "i"}},
					{"nama": bson.M{"$regex": pagination.Search, "$options": "i"}},
					{"jurusan": bson.M{"$regex": pagination.Search, "$options": "i"}},
				},
			}},
		}
	}

	// Count pipeline
	countPipeline := mongo.Pipeline{}
	if len(matchStage) > 0 {
		countPipeline = append(countPipeline, matchStage)
	}
	countPipeline = append(countPipeline, bson.D{{Key: "$count", Value: "total"}})

	// Get total count
	var total int64
	countCursor, err := r.collection.Aggregate(ctx, countPipeline)
	if err != nil {
		return nil, 0, err
	}
	defer countCursor.Close(ctx)

	var countResult []struct {
		Total int64 `bson:"total"`
	}
	if err = countCursor.All(ctx, &countResult); err != nil {
		return nil, 0, err
	}
	if len(countResult) > 0 {
		total = countResult[0].Total
	}

	// Build sort order
	sortOrder := 1
	if pagination.SortOrder == "DESC" {
		sortOrder = -1
	}

	// Data pipeline with lookup
	dataPipeline := mongo.Pipeline{}
	if len(matchStage) > 0 {
		dataPipeline = append(dataPipeline, matchStage)
	}
	dataPipeline = append(dataPipeline,
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "user_id"},
			{Key: "foreignField", Value: "id"},
			{Key: "as", Value: "user"},
		}}},
		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$user"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: pagination.SortBy, Value: sortOrder}}}},
		bson.D{{Key: "$skip", Value: pagination.GetOffset()}},
		bson.D{{Key: "$limit", Value: pagination.Limit}},
	)

	cursor, err := r.collection.Aggregate(ctx, dataPipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var alumnis []models.Alumni
	if err = cursor.All(ctx, &alumnis); err != nil {
		return nil, 0, err
	}

	return alumnis, total, nil
}

func (r *alumniRepositoryMongo) GetByID(id uint) (*models.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"id": id}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "user_id"},
			{Key: "foreignField", Value: "id"},
			{Key: "as", Value: "user"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$user"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var alumnis []models.Alumni
	if err = cursor.All(ctx, &alumnis); err != nil {
		return nil, err
	}

	if len(alumnis) == 0 {
		return nil, fmt.Errorf("alumni not found")
	}

	return &alumnis[0], nil
}

func (r *alumniRepositoryMongo) GetByUserID(userID int) (*models.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"user_id": userID}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "user_id"},
			{Key: "foreignField", Value: "id"},
			{Key: "as", Value: "user"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$user"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var alumnis []models.Alumni
	if err = cursor.All(ctx, &alumnis); err != nil {
		return nil, err
	}

	if len(alumnis) == 0 {
		return nil, fmt.Errorf("alumni not found")
	}

	return &alumnis[0], nil
}

func (r *alumniRepositoryMongo) Create(alumni *models.Alumni) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Set timestamps
	now := time.Now()
	alumni.CreatedAt = now
	alumni.UpdatedAt = now

	// Get next ID
	nextID, err := r.getNextSequenceID()
	if err != nil {
		return err
	}
	alumni.ID = nextID

	_, err = r.collection.InsertOne(ctx, alumni)
	return err
}

func (r *alumniRepositoryMongo) Update(alumni *models.Alumni) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	alumni.UpdatedAt = time.Now()

	filter := bson.M{"id": alumni.ID}
	update := bson.M{
		"$set": bson.M{
			"nim":         alumni.NIM,
			"nama":        alumni.Nama,
			"jurusan":     alumni.Jurusan,
			"angkatan":    alumni.Angkatan,
			"tahun_lulus": alumni.TahunLulus,
			"no_telepon":  alumni.NoTelepon,
			"alamat":      alumni.Alamat,
			"updated_at":  alumni.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("alumni not found")
	}

	return nil
}

func (r *alumniRepositoryMongo) Delete(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("alumni not found")
	}

	return nil
}

func (r *alumniRepositoryMongo) Count() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.collection.CountDocuments(ctx, bson.M{})
}

// Helper function to get next sequence ID
func (r *alumniRepositoryMongo) getNextSequenceID() (uint, error) {
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
