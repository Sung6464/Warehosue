package repository

import (
	"Customer-Services/model" // Fixed import path
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// CustomerRepository defines the interface for customer data operations.
type CustomerRepository interface {
	Create(ctx context.Context, customer model.Customer) error
	FindByID(ctx context.Context, id string) (model.Customer, error)
	FindAll(ctx context.Context) ([]model.Customer, error)
	Update(ctx context.Context, id string, update interface{}) error
	// New methods for array manipulation
	AddWarehouseToCustomer(ctx context.Context, customerID, warehouseID string) error
	RemoveWarehouseFromCustomer(ctx context.Context, customerID, warehouseID string) error
	FindCustomersByWarehouseID(ctx context.Context, warehouseID string) ([]model.Customer, error)
}

// mongoCustomerRepository implements CustomerRepository for MongoDB.
type mongoCustomerRepository struct {
	collection *mongo.Collection
}

// NewMongoCustomerRepository creates a new MongoDB repository for customers.
func NewMongoCustomerRepository(collection *mongo.Collection) CustomerRepository {
	return &mongoCustomerRepository{
		collection: collection,
	}
}

// Create inserts a new customer into the database.
func (r *mongoCustomerRepository) Create(ctx context.Context, customer model.Customer) error {
	_, err := r.collection.InsertOne(ctx, customer)
	return err
}

// FindByID retrieves a customer by its ID.
func (r *mongoCustomerRepository) FindByID(ctx context.Context, id string) (model.Customer, error) {
	var customer model.Customer
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&customer)
	if err == mongo.ErrNoDocuments {
		return model.Customer{}, mongo.ErrNoDocuments // Ensure this specific error is returned for 'not found'
	}
	return customer, err
}

// FindAll retrieves all customers from the database.
func (r *mongoCustomerRepository) FindAll(ctx context.Context) ([]model.Customer, error) {
	var customers []model.Customer
	cursor, err := r.collection.Find(ctx, bson.M{}) // Empty filter to get all
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx) // Always close the cursor

	if err = cursor.All(ctx, &customers); err != nil {
		return nil, err
	}
	return customers, nil
}

// Update updates an existing customer by its ID.
// This general update allows setting a new array of warehouse_ids directly, or other fields.
func (r *mongoCustomerRepository) Update(ctx context.Context, id string, update interface{}) error {
	// filter := bson.M{"_id": id} // Removed unused filter declaration
	result, err := r.collection.UpdateByID(ctx, id, bson.D{{Key: "$set", Value: update}}) // Uses 'id' directly
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments // Return specific error if not found
	}
	return nil
}

// AddWarehouseToCustomer adds a warehouse ID to a customer's list of associated warehouses.
// Uses $addToSet to add only if the ID is not already present, ensuring uniqueness within the array.
func (r *mongoCustomerRepository) AddWarehouseToCustomer(ctx context.Context, customerID, warehouseID string) error {
	filter := bson.M{"_id": customerID}
	update := bson.M{"$addToSet": bson.M{"warehouse_ids": warehouseID}}
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

// RemoveWarehouseFromCustomer removes a warehouse ID from a customer's list.
// Uses $pull to remove all instances of the specified ID.
func (r *mongoCustomerRepository) RemoveWarehouseFromCustomer(ctx context.Context, customerID, warehouseID string) error {
	filter := bson.M{"_id": customerID}
	update := bson.M{"$pull": bson.M{"warehouse_ids": warehouseID}}
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

// FindCustomersByWarehouseID finds all customers associated with a given warehouse ID.
// This is used by the Warehouse Service to query its associated customers.
func (r *mongoCustomerRepository) FindCustomersByWarehouseID(ctx context.Context, warehouseID string) ([]model.Customer, error) {
	var customers []model.Customer
	filter := bson.M{"warehouse_ids": warehouseID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &customers); err != nil {
		return nil, err
	}
	return customers, nil
}
