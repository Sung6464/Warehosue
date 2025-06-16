package repository

import (
	"commodity-service/model" // Fixed import path
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// CommodityRepository defines the interface for commodity data operations.
type CommodityRepository interface {
	Create(ctx context.Context, commodity model.Commodity) error
	FindByID(ctx context.Context, id string) (model.Commodity, error)
	FindAll(ctx context.Context) ([]model.Commodity, error)
	Update(ctx context.Context, id string, update interface{}) error
}

// mongoCommodityRepository implements CommodityRepository for MongoDB.
type mongoCommodityRepository struct {
	collection *mongo.Collection
}

// NewMongoCommodityRepository creates a new MongoDB repository for commodities.
func NewMongoCommodityRepository(collection *mongo.Collection) CommodityRepository {
	return &mongoCommodityRepository{
		collection: collection,
	}
}

// Create inserts a new commodity into the database.
func (r *mongoCommodityRepository) Create(ctx context.Context, commodity model.Commodity) error {
	_, err := r.collection.InsertOne(ctx, commodity)
	return err
}

// FindByID retrieves a commodity by its ID.
func (r *mongoCommodityRepository) FindByID(ctx context.Context, id string) (model.Commodity, error) {
	var commodity model.Commodity
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&commodity)
	if err == mongo.ErrNoDocuments {
		return model.Commodity{}, mongo.ErrNoDocuments
	}
	return commodity, err
}

// FindAll retrieves all commodities from the database.
func (r *mongoCommodityRepository) FindAll(ctx context.Context) ([]model.Commodity, error) {
	var commodities []model.Commodity
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &commodities); err != nil {
		return nil, err
	}
	return commodities, nil
}

// Update updates an existing commodity by its ID.
func (r *mongoCommodityRepository) Update(ctx context.Context, id string, update interface{}) error {
	result, err := r.collection.UpdateByID(ctx, id, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
