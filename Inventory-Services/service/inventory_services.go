package service

import (
	"Inventory-Services/model"
	"Inventory-Services/repository" // Added this import
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InventoryService defines the interface for inventory business logic.
type InventoryService interface {
	CreateInventory(ctx context.Context, inventory *model.Inventory) (*model.Inventory, error)
	GetAllInventories(ctx context.Context) ([]model.Inventory, error)
	GetInventoryByID(ctx context.Context, id string) (*model.Inventory, error)
	UpdateInventory(ctx context.Context, id string, inventory *model.Inventory) (*model.Inventory, error)
	DeleteInventory(ctx context.Context, id string) error
}

// inventoryServiceImpl implements InventoryService.
type inventoryServiceImpl struct {
	repository repository.InventoryRepository // Changed to use repository
}

// NewInventoryService creates a new instance of InventoryService.
func NewInventoryService() InventoryService {
	// We now create the repository and pass it to the service
	return &inventoryServiceImpl{repository: repository.NewInventoryRepository()}
}

func (s *inventoryServiceImpl) CreateInventory(ctx context.Context, inventory *model.Inventory) (*model.Inventory, error) {
	return s.repository.CreateInventory(ctx, inventory)
}

func (s *inventoryServiceImpl) GetAllInventories(ctx context.Context) ([]model.Inventory, error) {
	return s.repository.GetAllInventories(ctx)
}

func (s *inventoryServiceImpl) GetInventoryByID(ctx context.Context, id string) (*model.Inventory, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid inventory ID format")
	}
	return s.repository.GetInventoryByID(ctx, objID)
}

func (s *inventoryServiceImpl) UpdateInventory(ctx context.Context, id string, inventory *model.Inventory) (*model.Inventory, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid inventory ID format")
	}
	return s.repository.UpdateInventory(ctx, objID, inventory)
}

func (s *inventoryServiceImpl) DeleteInventory(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid inventory ID format")
	}
	return s.repository.DeleteInventory(ctx, objID)
}
