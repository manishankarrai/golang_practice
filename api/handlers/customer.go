package handlers

import (
	"errors"
	"net/http"

	"practice/api/models"
	"practice/api/services"

	"github.com/gin-gonic/gin"
)

func CreateCustomer(c *gin.Context) {
	var req models.CustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()}); return }
	customer, err := services.CreateCustomer(c.Request.Context(), models.Customer{Name: req.Name, Email: req.Email, Phone: req.Phone})
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()}); return }
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "customer created successfully", "data": customer})
}

func GetCustomers(c *gin.Context) {
	customers, err := services.GetCustomers(c.Request.Context())
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"success": true, "data": customers})
}

func GetCustomer(c *gin.Context) {
	customer, err := services.GetCustomerByID(c.Request.Context(), c.Param("id"))
	if errors.Is(err, services.ErrCustomerNotFound) { c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "customer not found"}); return }
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"success": true, "data": customer})
}

func UpdateCustomer(c *gin.Context) {
	var req models.CustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()}); return }
	customer, err := services.UpdateCustomer(c.Request.Context(), c.Param("id"), models.Customer{Name: req.Name, Email: req.Email, Phone: req.Phone})
	if errors.Is(err, services.ErrCustomerNotFound) { c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "customer not found"}); return }
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"success": true, "data": customer})
}

func UpdateCustomerLocation(c *gin.Context) {
	var req models.CustomerLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()}); return }
	customer, err := services.UpdateCustomerLocation(c.Request.Context(), c.Param("id"), req.Latitude, req.Longitude)
	if errors.Is(err, services.ErrCustomerNotFound) { c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "customer not found"}); return }
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"success": true, "data": customer})
}

func DeleteCustomer(c *gin.Context) {
	if err := services.DeleteCustomer(c.Request.Context(), c.Param("id")); errors.Is(err, services.ErrCustomerNotFound) { c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "customer not found"}); return } else if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "customer deleted successfully"})
}
