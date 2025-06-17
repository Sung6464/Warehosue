package repository

import (
	"Inventory-Services/database"
	"Inventory-Services/model"
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// InventoryRepository defines the interface for inventory data operations.
type InventoryRepository interface {
	CreateInventory(ctx context.Context, inventory *model.Inventory) (*model.Inventory, error)
	GetAllInventories(ctx context.Context) ([]model.Inventory, error)
	GetInventoryByID(ctx context.Context, id primitive.ObjectID) (*model.Inventory, error)
	UpdateInventory(ctx context.Context, id primitive.ObjectID, inventory *model.Inventory) (*model.Inventory, error)
	DeleteInventory(ctx context.Context, id primitive.ObjectID) error
}

// inventoryRepositoryImpl implements InventoryRepository.
type inventoryRepositoryImpl struct {
	collection *mongo.Collection
}

// NewInventoryRepository creates a new instance of InventoryRepository.
func NewInventoryRepository() InventoryRepository {
	// Ensure database.Client is initialized before calling GetCollection
	if database.Client == nil {
		log.Fatal("MongoDB client is not initialized. Call database.ConnectDB() first.")
	}
	collection := database.GetCollection(database.Client, "inventories")
	return &inventoryRepositoryImpl{collection: collection}
}

func (r *inventoryRepositoryImpl) CreateInventory(ctx context.Context, inventory *model.Inventory) (*model.Inventory, error) {
	result, err := r.collection.InsertOne(ctx, inventory)
	if err != nil {
		return nil, fmt.Errorf("failed to create inventory in repository: %w", err)
	}
	inventory.ID = result.InsertedID.(primitive.ObjectID)
	return inventory, nil
}

func (r *inventoryRepositoryImpl) GetAllInventories(ctx context.Context) ([]model.Inventory, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve inventories from repository: %w", err)
	}
	defer cursor.Close(ctx)

	var inventories []model.Inventory
	if err = cursor.All(ctx, &inventories); err != nil {
		return nil, fmt.Errorf("failed to decode inventories from cursor: %w", err)
	}
	return inventories, nil
}

func (r *inventoryRepositoryImpl) GetInventoryByID(ctx context.Context, id primitive.ObjectID) (*model.Inventory, error) {
	var inventory model.Inventory
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&inventory)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("inventory not found in repository")
		}
		return nil, fmt.Errorf("failed to retrieve inventory by ID from repository: %w", err)
	}
	return &inventory, nil
}

func (r *inventoryRepositoryImpl) UpdateInventory(ctx context.Context, id primitive.ObjectID, inventory *model.Inventory) (*model.Inventory, error) {
	updateDoc := bson.M{
		"$set": bson.M{
			"product_id":   inventory.ProductID,
			"quantity":     inventory.Quantity,
			"location":     inventory.Location,
			"last_updated": inventory.LastUpdated,
		},
	}

	result, err := r.collection.UpdateByID(ctx, id, updateDoc)
	if err != nil {
		return nil, fmt.Errorf("failed to update inventory in repository: %w", err)
	}
	if result.ModifiedCount == 0 {
		return nil, errors.New("inventory not found or no changes made in repository")
	}

	return r.GetInventoryByID(ctx, id)
}

func (r *inventoryRepositoryImpl) DeleteInventory(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete inventory from repository: %w", err)
	}
	if result.DeletedCount == 0 {
		return errors.New("inventory not found in repository")
	}
	return nil
}
