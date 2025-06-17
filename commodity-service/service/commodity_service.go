package service

import (
	"commodity-service/model"
	"commodity-service/repository"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CommodityService defines the interface for commodity business logic.
type CommodityService interface {
	CreateCommodity(ctx context.Context, commodity *model.Commodity) (*model.Commodity, error)
	GetAllCommodities(ctx context.Context) ([]model.Commodity, error)
	GetCommodityByID(ctx context.Context, id string) (*model.Commodity, error)
	UpdateCommodity(ctx context.Context, id string, commodity *model.Commodity) (*model.Commodity, error)
	DeleteCommodity(ctx context.Context, id string) error
}

// commodityServiceImpl implements CommodityService.
type commodityServiceImpl struct {
	repository repository.CommodityRepository
}

// NewCommodityService creates a new instance of CommodityService.
func NewCommodityService() CommodityService {
	return &commodityServiceImpl{repository: repository.NewCommodityRepository()}
}

func (s *commodityServiceImpl) CreateCommodity(ctx context.Context, commodity *model.Commodity) (*model.Commodity, error) {
	return s.repository.CreateCommodity(ctx, commodity)
}

func (s *commodityServiceImpl) GetAllCommodities(ctx context.Context) ([]model.Commodity, error) {
	return s.repository.GetAllCommodities(ctx)
}

func (s *commodityServiceImpl) GetCommodityByID(ctx context.Context, id string) (*model.Commodity, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid commodity ID format")
	}
	return s.repository.GetCommodityByID(ctx, objID)
}

func (s *commodityServiceImpl) UpdateCommodity(ctx context.Context, id string, commodity *model.Commodity) (*model.Commodity, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid commodity ID format")
	}
	return s.repository.UpdateCommodity(ctx, objID, commodity)
}

func (s *commodityServiceImpl) DeleteCommodity(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid commodity ID format")
	}
	return s.repository.DeleteCommodity(ctx, objID)
}
