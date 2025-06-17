package repository

import (
	"Warehouse-Services/database"
	"Warehouse-Services/model"
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// WarehouseRepository defines the interface for warehouse data operations.
type WarehouseRepository interface {
	CreateWarehouse(ctx context.Context, warehouse *model.Warehouse) (*model.Warehouse, error)
	GetAllWarehouses(ctx context.Context) ([]model.Warehouse, error)
	GetWarehouseByID(ctx context.Context, id primitive.ObjectID) (*model.Warehouse, error)
	UpdateWarehouse(ctx context.Context, id primitive.ObjectID, warehouse *model.Warehouse) (*model.Warehouse, error)
	DeleteWarehouse(ctx context.Context, id primitive.ObjectID) error
}

// warehouseRepositoryImpl implements WarehouseRepository.
type warehouseRepositoryImpl struct {
	collection *mongo.Collection
}

// NewWarehouseRepository creates a new instance of WarehouseRepository.
func NewWarehouseRepository() WarehouseRepository {
	if database.Client == nil {
		log.Fatal("MongoDB client is not initialized. Call database.ConnectDB() first.")
	}
	collection := database.GetCollection(database.Client, "warehouses")
	return &warehouseRepositoryImpl{collection: collection}
}

func (r *warehouseRepositoryImpl) CreateWarehouse(ctx context.Context, warehouse *model.Warehouse) (*model.Warehouse, error) {
	result, err := r.collection.InsertOne(ctx, warehouse)
	if err != nil {
		return nil, fmt.Errorf("failed to create warehouse in repository: %w", err)
	}
	warehouse.ID = result.InsertedID.(primitive.ObjectID)
	return warehouse, nil
}

func (r *warehouseRepositoryImpl) GetAllWarehouses(ctx context.Context) ([]model.Warehouse, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve warehouses from repository: %w", err)
	}
	defer cursor.Close(ctx)

	var warehouses []model.Warehouse
	if err = cursor.All(ctx, &warehouses); err != nil {
		return nil, fmt.Errorf("failed to decode warehouses from cursor: %w", err)
	}
	return warehouses, nil
}

func (r *warehouseRepositoryImpl) GetWarehouseByID(ctx context.Context, id primitive.ObjectID) (*model.Warehouse, error) {
	var warehouse model.Warehouse
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&warehouse)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("warehouse not found in repository")
		}
		return nil, fmt.Errorf("failed to retrieve warehouse by ID from repository: %w", err)
	}
	return &warehouse, nil
}

func (r *warehouseRepositoryImpl) UpdateWarehouse(ctx context.Context, id primitive.ObjectID, warehouse *model.Warehouse) (*model.Warehouse, error) {
	updateDoc := bson.M{
		"$set": bson.M{
			"name":     warehouse.Name,
			"location": warehouse.Location,
			"storage":  warehouse.Storage,
		},
	}

	result, err := r.collection.UpdateByID(ctx, id, updateDoc)
	if err != nil {
		return nil, fmt.Errorf("failed to update warehouse: %w", err)
	}
	if result.ModifiedCount == 0 {
		return nil, errors.New("warehouse not found or no changes made")
	}

	return r.GetWarehouseByID(ctx, id)
}

func (r *warehouseRepositoryImpl) DeleteWarehouse(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete warehouse: %w", err)
	}
	if result.DeletedCount == 0 {
		return errors.New("warehouse not found")
	}
	return nil
}
