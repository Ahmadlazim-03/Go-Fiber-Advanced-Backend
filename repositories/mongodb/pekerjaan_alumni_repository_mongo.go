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

type pekerjaanAlumniRepositoryMongo struct {
	collection        *mongo.Collection
	alumniCollection  *mongo.Collection
	userCollection    *mongo.Collection
}

func NewPekerjaanAlumniRepositoryMongo(db *mongo.Database) repo.PekerjaanAlumniRepository {
	return &pekerjaanAlumniRepositoryMongo{
		collection:       db.Collection("pekerjaan_alumnis"),
		alumniCollection: db.Collection("alumnis"),
		userCollection:   db.Collection("users"),
	}
}

func (r *pekerjaanAlumniRepositoryMongo) GetAll() ([]models.PekerjaanAlumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Use aggregation pipeline to join with alumnis and users collections
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"deleted_at": bson.M{"$eq": nil}}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "alumnis"},
			{Key: "localField", Value: "alumni_id"},
			{Key: "foreignField", Value: "id"},
			{Key: "as", Value: "alumni"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$alumni"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "alumni.user_id"},
			{Key: "foreignField", Value: "id"},
			{Key: "as", Value: "alumni.user"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$alumni.user"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "id", Value: -1}}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pekerjaans []models.PekerjaanAlumni
	if err = cursor.All(ctx, &pekerjaans); err != nil {
		return nil, err
	}

	return pekerjaans, nil
}

func (r *pekerjaanAlumniRepositoryMongo) GetWithPagination(pagination *models.PaginationRequest) ([]models.PekerjaanAlumni, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set default values
	pagination.SetDefaults()
	pagination.ValidateSortOrder()

	// Build search filter
	matchStage := bson.D{{Key: "$match", Value: bson.M{"deleted_at": bson.M{"$eq": nil}}}}
	if pagination.Search != "" {
		matchStage = bson.D{
			{Key: "$match", Value: bson.M{
				"deleted_at": bson.M{"$eq": nil},
				"$or": []bson.M{
					{"posisi_jabatan": bson.M{"$regex": pagination.Search, "$options": "i"}},
					{"nama_perusahaan": bson.M{"$regex": pagination.Search, "$options": "i"}},
					{"bidang_industri": bson.M{"$regex": pagination.Search, "$options": "i"}},
					{"lokasi_kerja": bson.M{"$regex": pagination.Search, "$options": "i"}},
				},
			}},
		}
	}

	// Count pipeline
	countPipeline := mongo.Pipeline{matchStage, bson.D{{Key: "$count", Value: "total"}}}

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
	dataPipeline := mongo.Pipeline{
		matchStage,
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "alumnis"},
			{Key: "localField", Value: "alumni_id"},
			{Key: "foreignField", Value: "id"},
			{Key: "as", Value: "alumni"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$alumni"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "alumni.user_id"},
			{Key: "foreignField", Value: "id"},
			{Key: "as", Value: "alumni.user"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$alumni.user"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: pagination.SortBy, Value: sortOrder}}}},
		{{Key: "$skip", Value: pagination.GetOffset()}},
		{{Key: "$limit", Value: pagination.Limit}},
	}

	cursor, err := r.collection.Aggregate(ctx, dataPipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var pekerjaans []models.PekerjaanAlumni
	if err = cursor.All(ctx, &pekerjaans); err != nil {
		return nil, 0, err
	}

	return pekerjaans, total, nil
}

func (r *pekerjaanAlumniRepositoryMongo) GetByID(id uint) (*models.PekerjaanAlumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"id": id, "deleted_at": bson.M{"$eq": nil}}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "alumnis"},
			{Key: "localField", Value: "alumni_id"},
			{Key: "foreignField", Value: "id"},
			{Key: "as", Value: "alumni"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$alumni"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "alumni.user_id"},
			{Key: "foreignField", Value: "id"},
			{Key: "as", Value: "alumni.user"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$alumni.user"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pekerjaans []models.PekerjaanAlumni
	if err = cursor.All(ctx, &pekerjaans); err != nil {
		return nil, err
	}

	if len(pekerjaans) == 0 {
		return nil, fmt.Errorf("pekerjaan alumni not found")
	}

	return &pekerjaans[0], nil
}

func (r *pekerjaanAlumniRepositoryMongo) GetByAlumniID(alumniID uint) ([]models.PekerjaanAlumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"alumni_id": alumniID, "deleted_at": bson.M{"$eq": nil}}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "alumnis"},
			{Key: "localField", Value: "alumni_id"},
			{Key: "foreignField", Value: "id"},
			{Key: "as", Value: "alumni"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$alumni"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "alumni.user_id"},
			{Key: "foreignField", Value: "id"},
			{Key: "as", Value: "alumni.user"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$alumni.user"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "id", Value: -1}}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pekerjaans []models.PekerjaanAlumni
	if err = cursor.All(ctx, &pekerjaans); err != nil {
		return nil, err
	}

	return pekerjaans, nil
}

func (r *pekerjaanAlumniRepositoryMongo) GetByUserID(userID int) ([]models.PekerjaanAlumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// First, get alumni IDs for this user
	var alumni []struct {
		ID uint `bson:"id"`
	}
	alumniCursor, err := r.alumniCollection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer alumniCursor.Close(ctx)
	
	if err = alumniCursor.All(ctx, &alumni); err != nil {
		return nil, err
	}

	if len(alumni) == 0 {
		return []models.PekerjaanAlumni{}, nil
	}

	alumniIDs := make([]uint, len(alumni))
	for i, a := range alumni {
		alumniIDs[i] = a.ID
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"alumni_id": bson.M{"$in": alumniIDs}, "deleted_at": bson.M{"$eq": nil}}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "alumnis"},
			{Key: "localField", Value: "alumni_id"},
			{Key: "foreignField", Value: "id"},
			{Key: "as", Value: "alumni"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$alumni"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "alumni.user_id"},
			{Key: "foreignField", Value: "id"},
			{Key: "as", Value: "alumni.user"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$alumni.user"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "id", Value: -1}}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pekerjaans []models.PekerjaanAlumni
	if err = cursor.All(ctx, &pekerjaans); err != nil {
		return nil, err
	}

	return pekerjaans, nil
}

func (r *pekerjaanAlumniRepositoryMongo) Create(pekerjaan *models.PekerjaanAlumni) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Set timestamps
	now := time.Now()
	pekerjaan.CreatedAt = now
	pekerjaan.UpdatedAt = now

	// Get next ID
	nextID, err := r.getNextSequenceID()
	if err != nil {
		return err
	}
	pekerjaan.ID = nextID

	_, err = r.collection.InsertOne(ctx, pekerjaan)
	return err
}

func (r *pekerjaanAlumniRepositoryMongo) Update(pekerjaan *models.PekerjaanAlumni) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pekerjaan.UpdatedAt = time.Now()

	filter := bson.M{"id": pekerjaan.ID, "deleted_at": bson.M{"$eq": nil}}
	update := bson.M{
		"$set": bson.M{
			"nama_perusahaan":       pekerjaan.NamaPerusahaan,
			"posisi_jabatan":        pekerjaan.PosisiJabatan,
			"bidang_industri":       pekerjaan.BidangIndustri,
			"lokasi_kerja":          pekerjaan.LokasiKerja,
			"gaji_range":            pekerjaan.GajiRange,
			"tanggal_mulai_kerja":   pekerjaan.TanggalMulaiKerja,
			"tanggal_selesai_kerja": pekerjaan.TanggalSelesaiKerja,
			"status_pekerjaan":      pekerjaan.StatusPekerjaan,
			"deskripsi_pekerjaan":   pekerjaan.DeskripsiPekerjaan,
			"updated_at":            pekerjaan.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("pekerjaan alumni not found")
	}

	return nil
}

func (r *pekerjaanAlumniRepositoryMongo) Delete(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if already soft deleted
	var pekerjaan struct {
		DeletedAt *time.Time `bson:"deleted_at"`
	}
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&pekerjaan)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("data pekerjaan alumni tidak ditemukan")
		}
		return err
	}

	if pekerjaan.DeletedAt == nil {
		return fmt.Errorf("tidak bisa hard delete: data belum di-soft delete terlebih dahulu")
	}

	filter := bson.M{"id": id, "deleted_at": bson.M{"$ne": nil}}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("tidak ada data yang dihapus - pastikan data sudah di-soft delete")
	}

	return nil
}

func (r *pekerjaanAlumniRepositoryMongo) SoftDelete(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now()
	filter := bson.M{"id": id, "deleted_at": bson.M{"$eq": nil}}
	update := bson.M{"$set": bson.M{"deleted_at": now}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("pekerjaan alumni not found or already deleted")
	}

	return nil
}

func (r *pekerjaanAlumniRepositoryMongo) SoftDeleteByAlumniID(alumniID uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now()
	filter := bson.M{"alumni_id": alumniID, "deleted_at": bson.M{"$eq": nil}}
	update := bson.M{"$set": bson.M{"deleted_at": now}}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}

func (r *pekerjaanAlumniRepositoryMongo) Restore(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"deleted_at": nil}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("pekerjaan alumni not found")
	}

	return nil
}

func (r *pekerjaanAlumniRepositoryMongo) GetDeleted() ([]models.PekerjaanAlumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"deleted_at": bson.M{"$ne": nil}}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "alumnis"},
			{Key: "localField", Value: "alumni_id"},
			{Key: "foreignField", Value: "id"},
			{Key: "as", Value: "alumni"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$alumni"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "alumni.user_id"},
			{Key: "foreignField", Value: "id"},
			{Key: "as", Value: "alumni.user"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$alumni.user"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "deleted_at", Value: -1}}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pekerjaans []models.PekerjaanAlumni
	if err = cursor.All(ctx, &pekerjaans); err != nil {
		return nil, err
	}

	return pekerjaans, nil
}

func (r *pekerjaanAlumniRepositoryMongo) GetDeletedByUserID(userID int) ([]models.PekerjaanAlumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// First, get alumni IDs for this user
	var alumni []struct {
		ID uint `bson:"id"`
	}
	alumniCursor, err := r.alumniCollection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer alumniCursor.Close(ctx)
	
	if err = alumniCursor.All(ctx, &alumni); err != nil {
		return nil, err
	}

	if len(alumni) == 0 {
		return []models.PekerjaanAlumni{}, nil
	}

	alumniIDs := make([]uint, len(alumni))
	for i, a := range alumni {
		alumniIDs[i] = a.ID
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"alumni_id": bson.M{"$in": alumniIDs}, "deleted_at": bson.M{"$ne": nil}}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "alumnis"},
			{Key: "localField", Value: "alumni_id"},
			{Key: "foreignField", Value: "id"},
			{Key: "as", Value: "alumni"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$alumni"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "alumni.user_id"},
			{Key: "foreignField", Value: "id"},
			{Key: "as", Value: "alumni.user"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$alumni.user"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "deleted_at", Value: -1}}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pekerjaans []models.PekerjaanAlumni
	if err = cursor.All(ctx, &pekerjaans); err != nil {
		return nil, err
	}

	return pekerjaans, nil
}

func (r *pekerjaanAlumniRepositoryMongo) Count() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.collection.CountDocuments(ctx, bson.M{"deleted_at": bson.M{"$eq": nil}})
}

func (r *pekerjaanAlumniRepositoryMongo) GetAlumniCountByCompany(namaPerusahaan string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"nama_perusahaan": namaPerusahaan, "deleted_at": bson.M{"$eq": nil}}}},
		{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$alumni_id"}}}},
		{{Key: "$count", Value: "total"}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result []struct {
		Total int64 `bson:"total"`
	}
	if err = cursor.All(ctx, &result); err != nil {
		return 0, err
	}

	if len(result) > 0 {
		return result[0].Total, nil
	}

	return 0, nil
}

// Helper function to get next sequence ID
func (r *pekerjaanAlumniRepositoryMongo) getNextSequenceID() (uint, error) {
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
