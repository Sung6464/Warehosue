package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"Customer-Services/model"
	"Customer-Services/service"

	"github.com/gin-gonic/gin"
)

// CustomerController handles HTTP requests for customers.
type CustomerController struct {
	service service.CustomerService
}

// NewCustomerController creates a new instance of CustomerController.
func NewCustomerController(s service.CustomerService) *CustomerController {
	return &CustomerController{
		service: s,
	}
}

// CreateCustomer handles POST requests to create a new customer.
func (ctrl *CustomerController) CreateCustomer(c *gin.Context) {
	var newCustomer model.Customer
	if err := c.ShouldBindJSON(&newCustomer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %v", err.Error())})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	createdCustomer, err := ctrl.service.CreateCustomer(ctx, newCustomer)
	if err != nil {
		if err.Error() == "customer name is required" || (len(err.Error()) >= 20 && err.Error()[0:20] == "invalid warehouse ID") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create customer: %v", err.Error())})
		}
		return
	}

	c.JSON(http.StatusCreated, createdCustomer)
}

// GetCustomer handles GET requests to retrieve a single customer by ID.
func (ctrl *CustomerController) GetCustomer(c *gin.Context) { // FIXED: Changed * 분위기 to *gin.Context
	customerID := c.Param("id")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	customer, err := ctrl.service.GetCustomerByID(ctx, customerID)
	if err != nil {
		if err.Error() == "customer not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch customer: %v", err.Error())})
		}
		return
	}

	c.JSON(http.StatusOK, customer)
}

// GetAllCustomers handles GET requests to retrieve all customers, with optional filtering by warehouse_id.
// This endpoint supports a query parameter `?warehouse_id=XYZ` to find customers associated with a specific warehouse.
func (ctrl *CustomerController) GetAllCustomers(c *gin.Context) {
	warehouseIDFilter := c.Query("warehouse_id") // Get optional query parameter from URL

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var customers []model.Customer
	var err error

	if warehouseIDFilter != "" {
		// If warehouse_id filter is present, call the service method that queries by it
		customers, err = ctrl.service.GetCustomersByWarehouseID(ctx, warehouseIDFilter)
	} else {
		// Otherwise, get all customers
		customers, err = ctrl.service.GetAllCustomers(ctx)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch customers: %v", err.Error())})
		return
	}

	c.JSON(http.StatusOK, customers)
}

// UpdateCustomer handles PUT requests to update an existing customer by ID.
// This allows general updates to customer fields, including replacing the entire warehouse_ids array.
func (ctrl *CustomerController) UpdateCustomer(c *gin.Context) {
	customerID := c.Param("id")

	var updatedData map[string]interface{}
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %v", err.Error())})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	updatedCustomer, err := ctrl.service.UpdateCustomer(ctx, customerID, updatedData)
	if err != nil {
		if err.Error() == "no fields provided for update" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if err.Error() == "customer not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if len(err.Error()) >= 20 && err.Error()[0:20] == "invalid warehouse ID" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if err.Error() == "warehouse_ids must be an array of strings" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update customer: %v", err.Error())})
		}
		return
	}

	c.JSON(http.StatusOK, updatedCustomer)
}

// AddWarehouseToCustomer handles POST requests to add a warehouse to a customer's list.
// The customer ID and warehouse ID are obtained from URL parameters.
func (ctrl *CustomerController) AddWarehouseToCustomer(c *gin.Context) {
	customerID := c.Param("id")
	warehouseID := c.Param("warehouse_id")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err := ctrl.service.AddWarehouseToCustomer(ctx, customerID, warehouseID)
	if err != nil {
		if err.Error() == "customer not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if len(err.Error()) >= 20 && err.Error()[0:20] == "invalid warehouse ID" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to add warehouse to customer: %v", err.Error())})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// RemoveWarehouseFromCustomer handles DELETE requests to remove a warehouse from a customer's list.
// The customer ID and warehouse ID are obtained from URL parameters.
func (ctrl *CustomerController) RemoveWarehouseFromCustomer(c *gin.Context) {
	customerID := c.Param("id")
	warehouseID := c.Param("warehouse_id")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err := ctrl.service.RemoveWarehouseFromCustomer(ctx, customerID, warehouseID)
	if err != nil {
		if err.Error() == "customer not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to remove warehouse from customer: %v", err.Error())})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
