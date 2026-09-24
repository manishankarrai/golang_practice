package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// ShopOrderItem identifies a flower and quantity purchased from a flower shop.
type ShopOrderItem struct {
	FlowerID bson.ObjectID `json:"flowerid" bson:"flowerid"`
	Name     string        `json:"name" bson:"name"`
	Quantity int           `json:"quantity" bson:"quantity"`
	Price    float64       `json:"price" bson:"price"`
	Subtotal float64       `json:"subtotal" bson:"subtotal"`
}

// ShopOrder represents an order placed by a flower buyer for a flower shop.
type ShopOrder struct {
	OrderID   bson.ObjectID  `json:"orderid,omitempty" bson:"orderid,omitempty"`
	BuyerID   bson.ObjectID  `json:"buyerid" bson:"buyerid"`
	ShopID    bson.ObjectID  `json:"shopid" bson:"shopid"`
	Items     []ShopOrderItem `json:"items" bson:"items"`
	Total     float64        `json:"total" bson:"total"`
	Status    string         `json:"status" bson:"status"`
	CreatedAt time.Time      `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" bson:"updated_at"`
}

// ShopOrderItemRequest is one requested flower and quantity.
type ShopOrderItemRequest struct {
	FlowerID string `json:"flowerid" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,gt=0"`
}

// ShopOrderRequest is the payload used to create or update an order.
type ShopOrderRequest struct {
	BuyerID string                 `json:"buyerid" binding:"required"`
	ShopID  string                 `json:"shopid" binding:"required"`
	Items   []ShopOrderItemRequest `json:"items" binding:"required,min=1"`
	Status  string                 `json:"status"`
}
