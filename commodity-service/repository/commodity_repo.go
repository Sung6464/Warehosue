package repository

import (
	"commodity-service/database"
	"commodity-service/model"
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CommodityRepository defines the interface for commodity data operations.
type CommodityRepository interface {
	CreateCommodity(ctx context.Context, commodity *model.Commodity) (*model.Commodity, error)
	GetAllCommodities(ctx context.Context) ([]model.Commodity, error)
	GetCommodityByID(ctx context.Context, id primitive.ObjectID) (*model.Commodity, error)
	UpdateCommodity(ctx context.Context, id primitive.ObjectID, commodity *model.Commodity) (*model.Commodity, error)
	DeleteCommodity(ctx context.Context, id primitive.ObjectID) error
}

// commodityRepositoryImpl implements CommodityRepository.
type commodityRepositoryImpl struct {
	collection *mongo.Collection
}

// NewCommodityRepository creates a new instance of CommodityRepository.
func NewCommodityRepository() CommodityRepository {
	if database.Client == nil {
		log.Fatal("MongoDB client is not initialized. Call database.ConnectDB() first.")
	}
	collection := database.GetCollection(database.Client, "commodities")
	return &commodityRepositoryImpl{collection: collection}
}

func (r *commodityRepositoryImpl) CreateCommodity(ctx context.Context, commodity *model.Commodity) (*model.Commodity, error) {
	result, err := r.collection.InsertOne(ctx, commodity)
	if err != nil {
		return nil, fmt.Errorf("failed to create commodity in repository: %w", err)
	}
	commodity.ID = result.InsertedID.(primitive.ObjectID)
	return commodity, nil
}

func (r *commodityRepositoryImpl) GetAllCommodities(ctx context.Context) ([]model.Commodity, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve commodities from repository: %w", err)
	}
	defer cursor.Close(ctx)

	var commodities []model.Commodity
	if err = cursor.All(ctx, &commodities); err != nil {
		return nil, fmt.Errorf("failed to decode commodities from cursor: %w", err)
	}
	return commodities, nil
}

func (r *commodityRepositoryImpl) GetCommodityByID(ctx context.Context, id primitive.ObjectID) (*model.Commodity, error) {
	var commodity model.Commodity
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&commodity)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("commodity not found in repository")
		}
		return nil, fmt.Errorf("failed to retrieve commodity by ID from repository: %w", err)
	}
	return &commodity, nil
}

func (r *commodityRepositoryImpl) UpdateCommodity(ctx context.Context, id primitive.ObjectID, commodity *model.Commodity) (*model.Commodity, error) {
	updateDoc := bson.M{
		"$set": bson.M{
			"name":   commodity.Name,
			"amount": commodity.Amount,
		},
	}

	result, err := r.collection.UpdateByID(ctx, id, updateDoc)
	if err != nil {
		return nil, fmt.Errorf("failed to update commodity in repository: %w", err)
	}
	if result.ModifiedCount == 0 {
		return nil, errors.New("commodity not found or no changes made in repository")
	}

	return r.GetCommodityByID(ctx, id)
}

func (r *commodityRepositoryImpl) DeleteCommodity(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete commodity from repository: %w", err)
	}
	if result.DeletedCount == 0 {
		return errors.New("commodity not found in repository")
	}
	return nil
}
