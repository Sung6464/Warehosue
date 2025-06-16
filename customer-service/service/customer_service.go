package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"Customer-Services/config"     // Fixed import path
	"Customer-Services/model"      // Fixed import path
	"Customer-Services/repository" // Fixed import path

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

// CustomerService defines the interface for customer business logic.
type CustomerService interface {
	CreateCustomer(ctx context.Context, newCustomer model.Customer) (model.Customer, error)
	GetCustomerByID(ctx context.Context, id string) (model.Customer, error)
	GetAllCustomers(ctx context.Context) ([]model.Customer, error)
	UpdateCustomer(ctx context.Context, id string, updatedData map[string]interface{}) (model.Customer, error)
	AddWarehouseToCustomer(ctx context.Context, customerID, warehouseID string) error
	RemoveWarehouseFromCustomer(ctx context.Context, customerID, warehouseID string) error
	GetCustomersByWarehouseID(ctx context.Context, warehouseID string) ([]model.Customer, error)
}

// customerServiceImpl implements CustomerService.
type customerServiceImpl struct {
	repo       repository.CustomerRepository
	httpClient *http.Client
}

// NewCustomerService creates a new instance of CustomerService.
func NewCustomerService(repo repository.CustomerRepository) CustomerService {
	return &customerServiceImpl{
		repo: repo,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// checkWarehouseExists calls the Warehouse Service to verify a warehouse by its ID.
func (s *customerServiceImpl) checkWarehouseExists(ctx context.Context, warehouseID string) (model.Warehouse, error) {
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

// CreateCustomer implements CustomerService.
func (s *customerServiceImpl) CreateCustomer(ctx context.Context, newCustomer model.Customer) (model.Customer, error) {
	if newCustomer.Name == "" {
		return model.Customer{}, fmt.Errorf("customer name is required")
	}

	for _, wid := range newCustomer.WarehouseIDs {
		_, err := s.checkWarehouseExists(ctx, wid)
		if err != nil {
			return model.Customer{}, fmt.Errorf("invalid warehouse ID '%s': %w", wid, err)
		}
	}

	newCustomer.ID = uuid.New().String()
	if err := s.repo.Create(ctx, newCustomer); err != nil {
		return model.Customer{}, fmt.Errorf("failed to create customer: %w", err)
	}
	return newCustomer, nil
}

// GetCustomerByID implements CustomerService.
func (s *customerServiceImpl) GetCustomerByID(ctx context.Context, id string) (model.Customer, error) {
	customer, err := s.repo.FindByID(ctx, id)
	if err == mongo.ErrNoDocuments {
		return model.Customer{}, fmt.Errorf("customer not found")
	}
	return customer, err
}

// GetAllCustomers implements CustomerService.
func (s *customerServiceImpl) GetAllCustomers(ctx context.Context) ([]model.Customer, error) {
	customers, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all customers: %w", err)
	}
	return customers, nil
}

// UpdateCustomer implements CustomerService.
func (s *customerServiceImpl) UpdateCustomer(ctx context.Context, id string, updatedData map[string]interface{}) (model.Customer, error) {
	if len(updatedData) == 0 {
		return model.Customer{}, fmt.Errorf("no fields provided for update")
	}

	if warehouseIDs, ok := updatedData["warehouse_ids"].([]interface{}); ok {
		for _, widIfc := range warehouseIDs {
			if wid, isString := widIfc.(string); isString {
				_, err := s.checkWarehouseExists(ctx, wid)
				if err != nil {
					return model.Customer{}, fmt.Errorf("invalid warehouse ID '%s' in update: %w", wid, err)
				}
			} else {
				return model.Customer{}, fmt.Errorf("warehouse_ids must be an array of strings")
			}
		}
	}

	if err := s.repo.Update(ctx, id, updatedData); err != nil {
		if err == mongo.ErrNoDocuments {
			return model.Customer{}, fmt.Errorf("customer not found")
		}
		return model.Customer{}, fmt.Errorf("failed to update customer: %w", err)
	}

	updatedCustomer, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return model.Customer{}, fmt.Errorf("failed to fetch updated customer: %w", err)
	}
	return updatedCustomer, nil
}

// AddWarehouseToCustomer adds a warehouse ID to a customer's list of associated warehouses.
func (s *customerServiceImpl) AddWarehouseToCustomer(ctx context.Context, customerID, warehouseID string) error {
	_, err := s.checkWarehouseExists(ctx, warehouseID)
	if err != nil {
		return fmt.Errorf("invalid warehouse ID to add: %w", err)
	}

	err = s.repo.AddWarehouseToCustomer(ctx, customerID, warehouseID)
	if err == mongo.ErrNoDocuments {
		return fmt.Errorf("customer not found")
	}
	return err
}

// RemoveWarehouseFromCustomer removes a warehouse ID from a customer's list.
func (s *customerServiceImpl) RemoveWarehouseFromCustomer(ctx context.Context, customerID, warehouseID string) error {
	err := s.repo.RemoveWarehouseFromCustomer(ctx, customerID, warehouseID)
	if err == mongo.ErrNoDocuments {
		return fmt.Errorf("customer not found")
	}
	return err
}

// GetCustomersByWarehouseID implements CustomerService to find customers associated with a specific warehouse.
func (s *customerServiceImpl) GetCustomersByWarehouseID(ctx context.Context, warehouseID string) ([]model.Customer, error) {
	customers, err := s.repo.FindCustomersByWarehouseID(ctx, warehouseID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch customers by warehouse ID: %w", err)
	}
	return customers, nil
}
