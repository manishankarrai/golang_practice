package handlers

import (
	"errors"
	"net/http"

	"practice/api/models"
	"practice/api/services"
	"github.com/gin-gonic/gin"
)

func CreateFlowerBuyer(c *gin.Context) { var req models.FlowerBuyerRequest; if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()}); return }; buyer, err := services.CreateFlowerBuyer(c, models.FlowerBuyer{Name: req.Name, Email: req.Email, Phone: req.Phone, Address: req.Address}); if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()}); return }; c.JSON(http.StatusCreated, gin.H{"success": true, "data": buyer}) }
func GetFlowerBuyers(c *gin.Context) { buyers, err := services.GetFlowerBuyers(c); if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()}); return }; c.JSON(http.StatusOK, gin.H{"success": true, "data": buyers}) }
func GetFlowerBuyer(c *gin.Context) { buyer, err := services.GetFlowerBuyerByID(c, c.Param("id")); if errors.Is(err, services.ErrFlowerBuyerNotFound) { c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()}); return }; if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()}); return }; c.JSON(http.StatusOK, gin.H{"success": true, "data": buyer}) }
func UpdateFlowerBuyer(c *gin.Context) { var req models.FlowerBuyerRequest; if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()}); return }; buyer, err := services.UpdateFlowerBuyer(c, c.Param("id"), models.FlowerBuyer{Name: req.Name, Email: req.Email, Phone: req.Phone, Address: req.Address}); if errors.Is(err, services.ErrFlowerBuyerNotFound) { c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()}); return }; if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()}); return }; c.JSON(http.StatusOK, gin.H{"success": true, "data": buyer}) }
func DeleteFlowerBuyer(c *gin.Context) { err := services.DeleteFlowerBuyer(c, c.Param("id")); if errors.Is(err, services.ErrFlowerBuyerNotFound) { c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()}); return }; if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()}); return }; c.JSON(http.StatusOK, gin.H{"success": true, "message": "flower buyer deleted successfully"}) }
