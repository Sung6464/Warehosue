package repository // The package for files in the 'repository' folder

import (
	"Inventory-Services/model" // FIXED: Matches go.mod module name
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// InventoryRepository defines the interface for inventory data operations.
type InventoryRepository interface {
	Create(ctx context.Context, item model.InventoryItem) error
	FindByID(ctx context.Context, id string) (model.InventoryItem, error)
	FindAll(ctx context.Context, filter bson.M) ([]model.InventoryItem, error)
	Update(ctx context.Context, id string, update interface{}) error
	UpdateQuantity(ctx context.Context, id string, newQuantity int) error
}

// mongoInventoryRepository implements InventoryRepository for MongoDB.
type mongoInventoryRepository struct {
	collection *mongo.Collection
}

// NewMongoInventoryRepository creates a new MongoDB repository for inventory items.
func NewMongoInventoryRepository(collection *mongo.Collection) InventoryRepository {
	return &mongoInventoryRepository{
		collection: collection,
	}
}

// Create inserts a new inventory item into the database.
func (r *mongoInventoryRepository) Create(ctx context.Context, item model.InventoryItem) error {
	_, err := r.collection.InsertOne(ctx, item)
	return err
}

// FindByID retrieves an inventory item by its ID.
func (r *mongoInventoryRepository) FindByID(ctx context.Context, id string) (model.InventoryItem, error) {
	var item model.InventoryItem
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&item)
	if err == mongo.ErrNoDocuments {
		return model.InventoryItem{}, mongo.ErrNoDocuments
	}
	return item, err
}

// FindAll retrieves inventory items based on a filter.
func (r *mongoInventoryRepository) FindAll(ctx context.Context, filter bson.M) ([]model.InventoryItem, error) {
	var items []model.InventoryItem
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

// Update updates an existing inventory item by its ID.
func (r *mongoInventoryRepository) Update(ctx context.Context, id string, update interface{}) error {
	result, err := r.collection.UpdateByID(ctx, id, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

// UpdateQuantity updates only the quantity of an inventory item.
func (r *mongoInventoryRepository) UpdateQuantity(ctx context.Context, id string, newQuantity int) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"quantity": newQuantity}}
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
