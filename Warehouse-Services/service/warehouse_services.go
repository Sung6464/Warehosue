package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"Warehouse-Services/config"
	"Warehouse-Services/model"
	"Warehouse-Services/repository"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson" // Needed for specific MongoDB error (ErrNoDocuments)
	"go.mongodb.org/mongo-driver/mongo"
)

// WarehouseService defines the interface for warehouse business logic.
type WarehouseService interface {
	CreateWarehouse(ctx context.Context, newWarehouse model.Warehouse) (model.Warehouse, error)
	GetWarehouseByID(ctx context.Context, id string) (model.Warehouse, error)
	GetAllWarehouses(ctx context.Context) ([]model.Warehouse, error)
	UpdateWarehouse(ctx context.Context, id string, updatedData map[string]interface{}) (model.Warehouse, error)
	GetInventoryInWarehouse(ctx context.Context, warehouseID string) ([]model.InventoryItem, error)
	// New methods for managing customer booking
	BookWarehouseForCustomer(ctx context.Context, warehouseID, customerID string) (model.Warehouse, error)
	UnbookWarehouseFromCustomer(ctx context.Context, warehouseID, customerID string) (model.Warehouse, error)
}

// warehouseServiceImpl implements WarehouseService.
type warehouseServiceImpl struct {
	repo       repository.WarehouseRepository
	httpClient *http.Client
}

// NewWarehouseService creates a new instance of WarehouseService.
func NewWarehouseService(repo repository.WarehouseRepository) WarehouseService {
	return &warehouseServiceImpl{
		repo: repo,
		httpClient: &http.Client{
			Timeout: 5 * time.Second, // Timeout for external HTTP calls
		},
	}
}

// checkCommodityExists calls the Commodities Service to verify a commodity by its ID.
func (s *warehouseServiceImpl) checkCommodityExists(ctx context.Context, commodityID string) (model.Commodity, error) {
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

// checkCustomerExists calls the Customers Service to verify a customer by its ID.
func (s *warehouseServiceImpl) checkCustomerExists(ctx context.Context, customerID string) (model.Customer, error) {
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

// CreateWarehouse implements WarehouseService.
func (s *warehouseServiceImpl) CreateWarehouse(ctx context.Context, newWarehouse model.Warehouse) (model.Warehouse, error) {
	if newWarehouse.Name == "" {
		return model.Warehouse{}, fmt.Errorf("warehouse name is required")
	}
	if newWarehouse.Address == "" {
		return model.Warehouse{}, fmt.Errorf("warehouse location (address) is required")
	}

	if newWarehouse.CommodityID != "" {
		_, err := s.checkCommodityExists(ctx, newWarehouse.CommodityID)
		if err != nil {
			return model.Warehouse{}, fmt.Errorf("invalid commodity ID: %w", err)
		}
	}
	// If a customer_id is provided at creation, validate and ensure it's not already booked
	if newWarehouse.CustomerID != "" {
		_, err := s.checkCustomerExists(ctx, newWarehouse.CustomerID)
		if err != nil {
			return model.Warehouse{}, fmt.Errorf("invalid customer ID: %w", err)
		}
		// When creating, the ID is new, so we don't need to check existing booking yet.
		// The booking logic (BookWarehouseForCustomer) handles this correctly for existing warehouses.
	}

	newWarehouse.ID = uuid.New().String()
	if err := s.repo.Create(ctx, newWarehouse); err != nil {
		return model.Warehouse{}, fmt.Errorf("failed to create warehouse: %w", err)
	}
	return newWarehouse, nil
}

// GetWarehouseByID implements WarehouseService.
func (s *warehouseServiceImpl) GetWarehouseByID(ctx context.Context, id string) (model.Warehouse, error) {
	warehouse, err := s.repo.FindByID(ctx, id)
	if err == mongo.ErrNoDocuments {
		return model.Warehouse{}, fmt.Errorf("warehouse not found")
	}
	return warehouse, err
}

// GetAllWarehouses implements WarehouseService.
func (s *warehouseServiceImpl) GetAllWarehouses(ctx context.Context) ([]model.Warehouse, error) {
	warehouses, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all warehouses: %w", err)
	}
	return warehouses, nil
}

// UpdateWarehouse implements WarehouseService.
func (s *warehouseServiceImpl) UpdateWarehouse(ctx context.Context, id string, updatedData map[string]interface{}) (model.Warehouse, error) {
	if len(updatedData) == 0 {
		return model.Warehouse{}, fmt.Errorf("no fields provided for update")
	}

	// Fetch existing warehouse to check booking status if customer_id is being updated
	existingWarehouse, err := s.repo.FindByID(ctx, id)
	if err == mongo.ErrNoDocuments {
		return model.Warehouse{}, fmt.Errorf("warehouse not found")
	}
	if err != nil {
		return model.Warehouse{}, fmt.Errorf("failed to fetch existing warehouse for update: %w", err)
	}

	if customerIDUpdate, ok := updatedData["customer_id"].(string); ok {
		// If attempting to set a customer_id (not empty string)
		if customerIDUpdate != "" {
			_, err := s.checkCustomerExists(ctx, customerIDUpdate)
			if err != nil {
				return model.Warehouse{}, fmt.Errorf("invalid customer ID for update: %w", err)
			}
			// Check if already booked by another customer
			if existingWarehouse.CustomerID != "" && existingWarehouse.CustomerID != customerIDUpdate {
				return model.Warehouse{}, fmt.Errorf("warehouse is already booked by another customer")
			}
		} else { // If customer_id is explicitly being set to empty string (unbooking)
			// No specific check needed, allows clearing customer_id
		}
	}

	if commodityID, ok := updatedData["commodity_id"].(string); ok && commodityID != "" {
		_, err := s.checkCommodityExists(ctx, commodityID)
		if err != nil {
			return model.Warehouse{}, fmt.Errorf("invalid commodity ID for update: %w", err)
		}
	}

	if err := s.repo.Update(ctx, id, updatedData); err != nil {
		if err == mongo.ErrNoDocuments {
			return model.Warehouse{}, fmt.Errorf("warehouse not found")
		}
		return model.Warehouse{}, fmt.Errorf("failed to update warehouse: %w", err)
	}

	updatedWarehouse, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return model.Warehouse{}, fmt.Errorf("failed to fetch updated warehouse: %w", err)
	}
	return updatedWarehouse, nil
}

// GetInventoryInWarehouse calls the Inventory Service to get all commodities in a specific warehouse.
func (s *warehouseServiceImpl) GetInventoryInWarehouse(ctx context.Context, warehouseID string) ([]model.InventoryItem, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/inventory?warehouse_id=%s", config.Cfg.InventoryServiceURL, warehouseID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for inventory service: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call inventory service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResponse map[string]string
		if decodeErr := json.NewDecoder(resp.Body).Decode(&errResponse); decodeErr == nil {
			return nil, fmt.Errorf("inventory service responded with status %d: %s", resp.StatusCode, errResponse["error"])
		}
		return nil, fmt.Errorf("inventory service responded with status %d", resp.StatusCode)
	}

	var inventoryItems []model.InventoryItem
	if err := json.NewDecoder(resp.Body).Decode(&inventoryItems); err != nil {
		return nil, fmt.Errorf("failed to decode inventory response: %w", err)
	}

	return inventoryItems, nil
}

// BookWarehouseForCustomer sets the CustomerID field on a warehouse.
// This is the core logic for the 1-to-many relationship (from customer's perspective).
func (s *warehouseServiceImpl) BookWarehouseForCustomer(ctx context.Context, warehouseID, customerID string) (model.Warehouse, error) {
	// 1. Validate customer ID exists (call Customer Service)
	_, err := s.checkCustomerExists(ctx, customerID)
	if err != nil {
		return model.Warehouse{}, fmt.Errorf("invalid customer ID '%s': %w", customerID, err)
	}

	// 2. Fetch the warehouse
	warehouse, err := s.repo.FindByID(ctx, warehouseID)
	if err == mongo.ErrNoDocuments {
		return model.Warehouse{}, fmt.Errorf("warehouse not found")
	}
	if err != nil {
		return model.Warehouse{}, fmt.Errorf("failed to retrieve warehouse for booking: %w", err)
	}

	// 3. Check if it's already booked by someone else
	if warehouse.CustomerID != "" && warehouse.CustomerID != customerID {
		return model.Warehouse{}, fmt.Errorf("warehouse '%s' is already booked by customer '%s'", warehouseID, warehouse.CustomerID)
	}
	// If it's already booked by the SAME customer, it's idempotent, just return success
	if warehouse.CustomerID == customerID {
		return warehouse, nil
	}

	// 4. Book the warehouse
	update := bson.M{"customer_id": customerID}
	if err := s.repo.Update(ctx, warehouseID, update); err != nil {
		if err == mongo.ErrNoDocuments { // Redundant but safe check
			return model.Warehouse{}, fmt.Errorf("warehouse not found during booking update")
		}
		return model.Warehouse{}, fmt.Errorf("failed to book warehouse: %w", err)
	}

	// 5. Fetch and return the updated warehouse
	return s.repo.FindByID(ctx, warehouseID)
}

// UnbookWarehouseFromCustomer removes the CustomerID from a warehouse.
func (s *warehouseServiceImpl) UnbookWarehouseFromCustomer(ctx context.Context, warehouseID, customerID string) (model.Warehouse, error) {
	// 1. Fetch the warehouse
	warehouse, err := s.repo.FindByID(ctx, warehouseID)
	if err == mongo.ErrNoDocuments {
		return model.Warehouse{}, fmt.Errorf("warehouse not found")
	}
	if err != nil {
		return model.Warehouse{}, fmt.Errorf("failed to retrieve warehouse for unbooking: %w", err)
	}

	// 2. Check if it's booked by the specified customer
	if warehouse.CustomerID != customerID {
		return model.Warehouse{}, fmt.Errorf("warehouse '%s' is not booked by customer '%s'", warehouseID, customerID)
	}

	// 3. Unbook the warehouse (set CustomerID to empty string)
	update := bson.M{"customer_id": ""}
	if err := s.repo.Update(ctx, warehouseID, update); err != nil {
		if err == mongo.ErrNoDocuments { // Redundant but safe check
			return model.Warehouse{}, fmt.Errorf("warehouse not found during unbooking update")
		}
		return model.Warehouse{}, fmt.Errorf("failed to unbook warehouse: %w", err)
	}

	// 4. Fetch and return the updated warehouse (with empty CustomerID)
	return s.repo.FindByID(ctx, warehouseID)
}
