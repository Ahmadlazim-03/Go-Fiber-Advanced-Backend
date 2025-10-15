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

type userRepositoryMongo struct {
	collection *mongo.Collection
}

func NewUserRepositoryMongo(db *mongo.Database) repo.UserRepository {
	return &userRepositoryMongo{
		collection: db.Collection("users"),
	}
}

func (r *userRepositoryMongo) GetAll() ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepositoryMongo) GetWithPagination(pagination *models.PaginationRequest) ([]models.User, int64, error) {
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
				{"username": bson.M{"$regex": pagination.Search, "$options": "i"}},
				{"email": bson.M{"$regex": pagination.Search, "$options": "i"}},
				{"role": bson.M{"$regex": pagination.Search, "$options": "i"}},
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

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepositoryMongo) GetByID(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	filter := bson.M{"id": id}
	
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Return nil when no record found
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepositoryMongo) GetByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	filter := bson.M{"email": email}
	
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Return nil when no record found
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepositoryMongo) GetByUsername(username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	filter := bson.M{"username": username}
	
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Return nil when no record found
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepositoryMongo) Create(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Set timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Get next ID
	nextID, err := r.getNextSequenceID()
	if err != nil {
		return err
	}
	user.ID = nextID

	_, err = r.collection.InsertOne(ctx, user)
	return err
}

func (r *userRepositoryMongo) Update(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user.UpdatedAt = time.Now()

	filter := bson.M{"id": user.ID}
	update := bson.M{
		"$set": bson.M{
			"username":   user.Username,
			"email":      user.Email,
			"password":   user.Password,
			"role":       user.Role,
			"is_active":  user.IsActive,
			"updated_at": user.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepositoryMongo) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepositoryMongo) Count() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.collection.CountDocuments(ctx, bson.M{})
}

// AuthenticateWithPassword is not supported for MongoDB
// MongoDB uses bcrypt password verification, not API authentication
func (r *userRepositoryMongo) AuthenticateWithPassword(email, password string) (*models.User, error) {
	return nil, fmt.Errorf("AuthenticateWithPassword not supported for MongoDB - use GetByEmail + bcrypt verification")
}

// Helper function to get next sequence ID
func (r *userRepositoryMongo) getNextSequenceID() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find the document with the highest ID
	findOptions := options.FindOne().SetSort(bson.D{{Key: "id", Value: -1}})
	var result struct {
		ID int `bson:"id"`
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
