package repository

import (
	"context"

	"Warehouse-Services/model" // Import the model package

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// WarehouseRepository defines the interface for warehouse data operations.
type WarehouseRepository interface {
	Create(ctx context.Context, warehouse model.Warehouse) error
	FindByID(ctx context.Context, id string) (model.Warehouse, error)
	FindAll(ctx context.Context) ([]model.Warehouse, error)
	Update(ctx context.Context, id string, update interface{}) error
}

// mongoWarehouseRepository implements WarehouseRepository for MongoDB.
type mongoWarehouseRepository struct {
	collection *mongo.Collection
}

// NewMongoWarehouseRepository creates a new MongoDB repository for warehouses.
func NewMongoWarehouseRepository(collection *mongo.Collection) WarehouseRepository {
	return &mongoWarehouseRepository{
		collection: collection,
	}
}

// Create inserts a new warehouse into the database.
func (r *mongoWarehouseRepository) Create(ctx context.Context, warehouse model.Warehouse) error {
	_, err := r.collection.InsertOne(ctx, warehouse)
	return err
}

// FindByID retrieves a warehouse by its ID.
func (r *mongoWarehouseRepository) FindByID(ctx context.Context, id string) (model.Warehouse, error) {
	var warehouse model.Warehouse
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&warehouse)
	return warehouse, err
}

// FindAll retrieves all warehouses from the database.
func (r *mongoWarehouseRepository) FindAll(ctx context.Context) ([]model.Warehouse, error) {
	var warehouses []model.Warehouse
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &warehouses); err != nil {
		return nil, err
	}
	return warehouses, nil
}

// Update updates an existing warehouse by its ID.
func (r *mongoWarehouseRepository) Update(ctx context.Context, id string, update interface{}) error {
	// The 'filter' variable was declared but not used by UpdateByID,
	// as UpdateByID takes the ID directly. Removed for cleanliness.
	result, err := r.collection.UpdateByID(ctx, id, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments // Return specific error if not found
	}
	return nil
}
