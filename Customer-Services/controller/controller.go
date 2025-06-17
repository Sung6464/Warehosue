package controller

import (
	"Customer-Services/model"   // Corrected import path
	"Customer-Services/service" // Corrected import path
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CustomerController handles HTTP requests related to customers.
type CustomerController struct {
	customerService service.CustomerService
}

// NewCustomerController creates a new instance of CustomerController.
func NewCustomerController(s service.CustomerService) *CustomerController {
	return &CustomerController{customerService: s}
}

// CreateCustomer handles POST /customers requests.
func (c *CustomerController) CreateCustomer(ctx *gin.Context) {
	var customer model.Customer
	if err := ctx.ShouldBindJSON(&customer); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if customer.FirstName == "" || customer.LastName == "" || customer.Email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "First name, last name, and email are required"})
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	createdCustomer, err := c.customerService.CreateCustomer(timeoutCtx, &customer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, createdCustomer)
}

// GetAllCustomers handles GET /customers requests.
func (c *CustomerController) GetAllCustomers(ctx *gin.Context) {
	timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	customers, err := c.customerService.GetAllCustomers(timeoutCtx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, customers)
}

// GetCustomerByID handles GET /customers/:id requests.
func (c *CustomerController) GetCustomerByID(ctx *gin.Context) {
	id := ctx.Param("id")

	timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	customer, err := c.customerService.GetCustomerByID(timeoutCtx, id)
	if err != nil {
		if err.Error() == "customer not found" || err.Error() == "invalid customer ID format" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, customer)
}

// UpdateCustomer handles PUT /customers/:id requests.
func (c *CustomerController) UpdateCustomer(ctx *gin.Context) {
	id := ctx.Param("id")
	var customer model.Customer
	if err := ctx.ShouldBindJSON(&customer); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	updatedCustomer, err := c.customerService.UpdateCustomer(timeoutCtx, id, &customer)
	if err != nil {
		if err.Error() == "customer not found" || err.Error() == "invalid customer ID format" || err.Error() == "customer not found or no changes made" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, updatedCustomer)
}

// DeleteCustomer handles DELETE /customers/:id requests.
func (c *CustomerController) DeleteCustomer(ctx *gin.Context) {
	id := ctx.Param("id")

	timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	err := c.customerService.DeleteCustomer(timeoutCtx, id)
	if err != nil {
		if err.Error() == "customer not found" || err.Error() == "invalid customer ID format" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusNoContent, nil) // 204 No Content for successful deletion
}
