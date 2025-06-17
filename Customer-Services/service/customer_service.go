package service

import (
	"Customer-Services/database" // Corrected import path
	"Customer-Services/model"    // Corrected import path
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CustomerService defines the interface for customer business logic.
type CustomerService interface {
	CreateCustomer(ctx context.Context, customer *model.Customer) (*model.Customer, error)
	GetAllCustomers(ctx context.Context) ([]model.Customer, error)
	GetCustomerByID(ctx context.Context, id string) (*model.Customer, error)
	UpdateCustomer(ctx context.Context, id string, customer *model.Customer) (*model.Customer, error)
	DeleteCustomer(ctx context.Context, id string) error
}

// customerServiceImpl implements CustomerService.
type customerServiceImpl struct {
	collection *mongo.Collection
}

// NewCustomerService creates a new instance of CustomerService.
func NewCustomerService() CustomerService {
	// Ensure database.Client is initialized before calling GetCollection
	if database.Client == nil {
		log.Fatal("MongoDB client is not initialized. Call database.ConnectDB() first.")
	}
	collection := database.GetCollection(database.Client, "customers")
	return &customerServiceImpl{collection: collection}
}

func (s *customerServiceImpl) CreateCustomer(ctx context.Context, customer *model.Customer) (*model.Customer, error) {
	result, err := s.collection.InsertOne(ctx, customer)
	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}
	customer.ID = result.InsertedID.(primitive.ObjectID)
	return customer, nil
}

func (s *customerServiceImpl) GetAllCustomers(ctx context.Context) ([]model.Customer, error) {
	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve customers: %w", err)
	}
	defer cursor.Close(ctx)

	var customers []model.Customer
	if err = cursor.All(ctx, &customers); err != nil {
		return nil, fmt.Errorf("failed to decode customers: %w", err)
	}
	return customers, nil
}

func (s *customerServiceImpl) GetCustomerByID(ctx context.Context, id string) (*model.Customer, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid customer ID format")
	}

	var customer model.Customer
	err = s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&customer)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("customer not found")
		}
		return nil, fmt.Errorf("failed to retrieve customer by ID: %w", err)
	}
	return &customer, nil
}

func (s *customerServiceImpl) UpdateCustomer(ctx context.Context, id string, customer *model.Customer) (*model.Customer, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid customer ID format")
	}

	updateDoc := bson.M{
		"$set": bson.M{
			"first_name": customer.FirstName,
			"last_name":  customer.LastName,
			"email":      customer.Email,
			"phone":      customer.Phone,
			"address":    customer.Address,
		},
	}

	result, err := s.collection.UpdateByID(ctx, objID, updateDoc)
	if err != nil {
		return nil, fmt.Errorf("failed to update customer: %w", err)
	}
	if result.ModifiedCount == 0 {
		return nil, errors.New("customer not found or no changes made")
	}

	return s.GetCustomerByID(ctx, id)
}

func (s *customerServiceImpl) DeleteCustomer(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid customer ID format")
	}

	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}
	if result.DeletedCount == 0 {
		return errors.New("customer not found")
	}
	return nil
}
