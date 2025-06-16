package service

import (
	"context"
	"fmt"

	"commodity-service/model"      // Fixed import path
	"commodity-service/repository" // Fixed import path

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

// CommodityService defines the interface for commodity business logic.
type CommodityService interface {
	CreateCommodity(ctx context.Context, newCommodity model.Commodity) (model.Commodity, error)
	GetCommodityByID(ctx context.Context, id string) (model.Commodity, error)
	GetAllCommodities(ctx context.Context) ([]model.Commodity, error)
	UpdateCommodity(ctx context.Context, id string, updatedData map[string]interface{}) (model.Commodity, error)
}

// commodityServiceImpl implements CommodityService.
type commodityServiceImpl struct {
	repo repository.CommodityRepository
}

// NewCommodityService creates a new instance of CommodityService.
func NewCommodityService(repo repository.CommodityRepository) CommodityService {
	return &commodityServiceImpl{
		repo: repo,
	}
}

// CreateCommodity implements CommodityService.
func (s *commodityServiceImpl) CreateCommodity(ctx context.Context, newCommodity model.Commodity) (model.Commodity, error) {
	if newCommodity.Name == "" {
		return model.Commodity{}, fmt.Errorf("commodity name is required")
	}
	if newCommodity.Amount <= 0 {
		return model.Commodity{}, fmt.Errorf("commodity amount must be positive")
	}

	newCommodity.ID = uuid.New().String()
	if err := s.repo.Create(ctx, newCommodity); err != nil {
		return model.Commodity{}, fmt.Errorf("failed to create commodity: %w", err)
	}
	return newCommodity, nil
}

// GetCommodityByID implements CommodityService.
func (s *commodityServiceImpl) GetCommodityByID(ctx context.Context, id string) (model.Commodity, error) {
	commodity, err := s.repo.FindByID(ctx, id)
	if err == mongo.ErrNoDocuments {
		return model.Commodity{}, fmt.Errorf("commodity not found")
	}
	return commodity, err
}

// GetAllCommodities implements CommodityService.
func (s *commodityServiceImpl) GetAllCommodities(ctx context.Context) ([]model.Commodity, error) {
	commodities, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all commodities: %w", err)
	}
	return commodities, nil
}

// UpdateCommodity implements CommodityService.
func (s *commodityServiceImpl) UpdateCommodity(ctx context.Context, id string, updatedData map[string]interface{}) (model.Commodity, error) {
	if len(updatedData) == 0 {
		return model.Commodity{}, fmt.Errorf("no fields provided for update")
	}

	if amount, ok := updatedData["amount"].(float64); ok {
		if amount < 0 {
			return model.Commodity{}, fmt.Errorf("commodity amount cannot be negative")
		}
		updatedData["amount"] = int(amount)
	}

	if err := s.repo.Update(ctx, id, updatedData); err != nil {
		if err == mongo.ErrNoDocuments {
			return model.Commodity{}, fmt.Errorf("commodity not found")
		}
		return model.Commodity{}, fmt.Errorf("failed to update commodity: %w", err)
	}

	updatedCommodity, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return model.Commodity{}, fmt.Errorf("failed to fetch updated commodity: %w", err)
	}
	return updatedCommodity, nil
}
