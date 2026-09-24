package handlers

import (
	"errors"
	"net/http"

	"practice/api/models"
	"practice/api/services"
	"github.com/gin-gonic/gin"
)

func CreateShopOrder(c *gin.Context) { var req models.ShopOrderRequest; if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()}); return }; order, err := services.CreateShopOrder(c, req); if errors.Is(err, services.ErrInvalidShopOrder) { c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()}); return }; if err != nil { c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()}); return }; c.JSON(http.StatusCreated, gin.H{"success": true, "data": order}) }
func GetShopOrders(c *gin.Context) { orders, err := services.GetShopOrders(c); if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()}); return }; c.JSON(http.StatusOK, gin.H{"success": true, "data": orders}) }
func GetShopOrder(c *gin.Context) { order, err := services.GetShopOrderByID(c, c.Param("id")); if errors.Is(err, services.ErrShopOrderNotFound) { c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()}); return }; if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()}); return }; c.JSON(http.StatusOK, gin.H{"success": true, "data": order}) }
func UpdateShopOrder(c *gin.Context) { var req models.ShopOrderRequest; if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()}); return }; order, err := services.UpdateShopOrder(c, c.Param("id"), req); if errors.Is(err, services.ErrShopOrderNotFound) { c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()}); return }; if errors.Is(err, services.ErrInvalidShopOrder) { c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()}); return }; if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()}); return }; c.JSON(http.StatusOK, gin.H{"success": true, "data": order}) }
func DeleteShopOrder(c *gin.Context) { err := services.DeleteShopOrder(c, c.Param("id")); if errors.Is(err, services.ErrShopOrderNotFound) { c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()}); return }; if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()}); return }; c.JSON(http.StatusOK, gin.H{"success": true, "message": "shop order deleted successfully"}) }
