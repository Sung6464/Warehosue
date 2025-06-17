package service

import (
	"Warehouse-Services/model"
	"Warehouse-Services/repository"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// WarehouseService defines the interface for warehouse business logic.
type WarehouseService interface {
	CreateWarehouse(ctx context.Context, warehouse *model.Warehouse) (*model.Warehouse, error)
	GetAllWarehouses(ctx context.Context) ([]model.Warehouse, error)
	GetWarehouseByID(ctx context.Context, id string) (*model.Warehouse, error)
	UpdateWarehouse(ctx context.Context, id string, warehouse *model.Warehouse) (*model.Warehouse, error)
	DeleteWarehouse(ctx context.Context, id string) error
}

// warehouseServiceImpl implements WarehouseService.
type warehouseServiceImpl struct {
	repository repository.WarehouseRepository
}

// NewWarehouseService creates a new instance of WarehouseService.
func NewWarehouseService() WarehouseService {
	return &warehouseServiceImpl{repository: repository.NewWarehouseRepository()}
}

func (s *warehouseServiceImpl) CreateWarehouse(ctx context.Context, warehouse *model.Warehouse) (*model.Warehouse, error) {
	return s.repository.CreateWarehouse(ctx, warehouse)
}

func (s *warehouseServiceImpl) GetAllWarehouses(ctx context.Context) ([]model.Warehouse, error) {
	return s.repository.GetAllWarehouses(ctx)
}

func (s *warehouseServiceImpl) GetWarehouseByID(ctx context.Context, id string) (*model.Warehouse, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid warehouse ID format")
	}
	return s.repository.GetWarehouseByID(ctx, objID)
}

func (s *warehouseServiceImpl) UpdateWarehouse(ctx context.Context, id string, warehouse *model.Warehouse) (*model.Warehouse, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid warehouse ID format")
	}
	return s.repository.UpdateWarehouse(ctx, objID, warehouse)
}

func (s *warehouseServiceImpl) DeleteWarehouse(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid warehouse ID format")
	}
	return s.repository.DeleteWarehouse(ctx, objID)
}
