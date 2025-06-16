package service // The package for files in the 'service' folder

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"Inventory-Services/config"
	"Inventory-Services/model"
	"Inventory-Services/repository"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// InventoryService defines the interface for inventory business logic.
type InventoryService interface {
	CreateInventoryItem(ctx context.Context, newItem model.InventoryItem) (model.InventoryItem, error)
	GetInventoryItemByID(ctx context.Context, id string) (model.InventoryItem, error)
	GetAllInventoryItems(ctx context.Context, warehouseID, commodityID, customerID string) ([]model.InventoryItem, error)
	UpdateInventoryItem(ctx context.Context, id string, updatedData map[string]interface{}) (model.InventoryItem, error)
	AdjustInventoryQuantity(ctx context.Context, id string, quantityChange int) (model.InventoryItem, error)
}

// inventoryServiceImpl implements InventoryService.
type inventoryServiceImpl struct {
	repo       repository.InventoryRepository
	httpClient *http.Client
}

// NewInventoryService creates a new instance of InventoryService.
func NewInventoryService(repo repository.InventoryRepository) InventoryService {
	return &inventoryServiceImpl{
		repo: repo,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// checkCommodityExists calls the Commodities Service to verify a commodity by its ID.
func (s *inventoryServiceImpl) checkCommodityExists(ctx context.Context, commodityID string) (model.Commodity, error) {
	if commodityID == "" {
		return model.Commodity{}, fmt.Errorf("commodity ID cannot be empty")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, config.Cfg.CommoditiesServiceURL+"/commodities/"+commodityID, nil)
	if err != nil {
		return model.Commodity{}, fmt.Errorf("failed to create request for commodities service: %w", err)
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return model.Commodity{}, fmt.Errorf("failed to call commodities service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResponse map[string]string
		if decodeErr := json.NewDecoder(resp.Body).Decode(&errResponse); decodeErr == nil {
			return model.Commodity{}, fmt.Errorf("commodities service responded with status %d: %s", resp.StatusCode, errResponse["error"])
		}
		return model.Commodity{}, fmt.Errorf("commodities service responded with status %d", resp.StatusCode)
	}
	var commodity model.Commodity
	if err := json.NewDecoder(resp.Body).Decode(&commodity); err != nil {
		return model.Commodity{}, fmt.Errorf("failed to decode commodity response: %w", err)
	}
	return commodity, nil
}

// checkWarehouseExists calls the Warehouse Service to verify a warehouse by its ID.
func (s *inventoryServiceImpl) checkWarehouseExists(ctx context.Context, warehouseID string) (model.Warehouse, error) {
	if warehouseID == "" {
		return model.Warehouse{}, fmt.Errorf("warehouse ID cannot be empty")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, config.Cfg.WarehouseServiceURL+"/warehouses/"+warehouseID, nil)
	if err != nil {
		return model.Warehouse{}, fmt.Errorf("failed to create request for warehouse service: %w", err)
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return model.Warehouse{}, fmt.Errorf("failed to call warehouse service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResponse map[string]string
		if decodeErr := json.NewDecoder(resp.Body).Decode(&errResponse); decodeErr == nil {
			return model.Warehouse{}, fmt.Errorf("warehouse service responded with status %d: %s", resp.StatusCode, errResponse["error"])
		}
		return model.Warehouse{}, fmt.Errorf("warehouse service responded with status %d", resp.StatusCode)
	}
	var warehouse model.Warehouse
	if err := json.NewDecoder(resp.Body).Decode(&warehouse); err != nil {
		return model.Warehouse{}, fmt.Errorf("failed to decode warehouse response: %w", err)
	}
	return warehouse, nil
}

// checkCustomerExists calls the Customers Service to verify a customer by its ID.
func (s *inventoryServiceImpl) checkCustomerExists(ctx context.Context, customerID string) (model.Customer, error) {
	if customerID == "" {
		return model.Customer{}, fmt.Errorf("customer ID cannot be empty")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, config.Cfg.CustomerServiceURL+"/customers/"+customerID, nil)
	if err != nil {
		return model.Customer{}, fmt.Errorf("failed to create request for customers service: %w", err)
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return model.Customer{}, fmt.Errorf("failed to call customers service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResponse map[string]string
		if decodeErr := json.NewDecoder(resp.Body).Decode(&errResponse); decodeErr == nil {
			return model.Customer{}, fmt.Errorf("customers service responded with status %d: %s", resp.StatusCode, errResponse["error"])
		}
		return model.Customer{}, fmt.Errorf("customers service responded with status %d", resp.StatusCode)
	}
	var customer model.Customer
	if err := json.NewDecoder(resp.Body).Decode(&customer); err != nil {
		return model.Customer{}, fmt.Errorf("failed to decode customer response: %w", err)
	}
	return customer, nil
}

// CreateInventoryItem creates a new inventory record.
func (s *inventoryServiceImpl) CreateInventoryItem(ctx context.Context, newItem model.InventoryItem) (model.InventoryItem, error) {
	if newItem.WarehouseID == "" || newItem.CommodityID == "" || newItem.Quantity <= 0 {
		return model.InventoryItem{}, fmt.Errorf("warehouse ID, commodity ID, and positive quantity are required")
	}

	_, err := s.checkWarehouseExists(ctx, newItem.WarehouseID)
	if err != nil {
		return model.InventoryItem{}, fmt.Errorf("invalid warehouse ID: %w", err)
	}
	_, err = s.checkCommodityExists(ctx, newItem.CommodityID)
	if err != nil {
		return model.InventoryItem{}, fmt.Errorf("invalid commodity ID: %w", err)
	}

	if newItem.CustomerID != "" {
		_, err := s.checkCustomerExists(ctx, newItem.CustomerID)
		if err != nil {
			return model.InventoryItem{}, fmt.Errorf("invalid customer ID: %w", err)
		}
	}

	newItem.ID = uuid.New().String()
	newItem.LastUpdated = time.Now()

	if err := s.repo.Create(ctx, newItem); err != nil {
		return model.InventoryItem{}, fmt.Errorf("failed to create inventory item: %w", err)
	}
	return newItem, nil
}

// GetInventoryItemByID retrieves a single inventory item.
func (s *inventoryServiceImpl) GetInventoryItemByID(ctx context.Context, id string) (model.InventoryItem, error) {
	item, err := s.repo.FindByID(ctx, id)
	if err == mongo.ErrNoDocuments {
		return model.InventoryItem{}, fmt.Errorf("inventory item not found")
	}
	return item, err
}

// GetAllInventoryItems retrieves all inventory items, with optional filters.
func (s *inventoryServiceImpl) GetAllInventoryItems(ctx context.Context, warehouseID, commodityID, customerID string) ([]model.InventoryItem, error) {
	filter := bson.M{}
	if warehouseID != "" {
		filter["warehouse_id"] = warehouseID
	}
	if commodityID != "" {
		filter["commodity_id"] = commodityID
	}
	if customerID != "" {
		filter["customer_id"] = customerID
	}

	items, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch inventory items: %w", err)
	}
	return items, nil
}

// UpdateInventoryItem updates an existing inventory item.
func (s *inventoryServiceImpl) UpdateInventoryItem(ctx context.Context, id string, updatedData map[string]interface{}) (model.InventoryItem, error) {
	if len(updatedData) == 0 {
		return model.InventoryItem{}, fmt.Errorf("no fields provided for update")
	}

	if wid, ok := updatedData["warehouse_id"].(string); ok && wid != "" {
		_, err := s.checkWarehouseExists(ctx, wid)
		if err != nil {
			return model.InventoryItem{}, fmt.Errorf("invalid warehouse ID for update: %w", err)
		}
	}
	if cid, ok := updatedData["commodity_id"].(string); ok && cid != "" {
		_, err := s.checkCommodityExists(ctx, cid)
		if err != nil {
			return model.InventoryItem{}, fmt.Errorf("invalid commodity ID for update: %w", err)
		}
	}
	if custID, ok := updatedData["customer_id"].(string); ok && custID != "" {
		_, err := s.checkCustomerExists(ctx, custID)
		if err != nil {
			return model.InventoryItem{}, fmt.Errorf("invalid customer ID for update: %w", err)
		}
	}

	updatedData["last_updated"] = time.Now()

	if err := s.repo.Update(ctx, id, updatedData); err != nil {
		if err == mongo.ErrNoDocuments {
			return model.InventoryItem{}, fmt.Errorf("inventory item not found")
		}
		return model.InventoryItem{}, fmt.Errorf("failed to update inventory item: %w", err)
	}

	updatedItem, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return model.InventoryItem{}, fmt.Errorf("failed to fetch updated inventory item: %w", err)
	}
	return updatedItem, nil
}

// AdjustInventoryQuantity adjusts the quantity of a specific inventory item.
func (s *inventoryServiceImpl) AdjustInventoryQuantity(ctx context.Context, id string, quantityChange int) (model.InventoryItem, error) {
	currentItem, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.InventoryItem{}, fmt.Errorf("inventory item not found")
		}
		return model.InventoryItem{}, fmt.Errorf("failed to retrieve current quantity: %w", err)
	}

	newQuantity := currentItem.Quantity + quantityChange
	if newQuantity < 0 {
		return model.InventoryItem{}, fmt.Errorf("insufficient stock: cannot reduce quantity below zero")
	}

	updateMap := map[string]interface{}{
		"quantity":     newQuantity,
		"last_updated": time.Now(),
	}
	err = s.repo.Update(ctx, id, updateMap)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.InventoryItem{}, fmt.Errorf("inventory item not found during quantity update")
		}
		return model.InventoryItem{}, fmt.Errorf("failed to adjust quantity: %w", err)
	}

	return s.repo.FindByID(ctx, id)
}
