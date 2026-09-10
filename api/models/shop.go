package models

import "go.mongodb.org/mongo-driver/v2/bson"

// Shop represents a single shop.
type Shop struct {
	ID      bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name    string        `json:"name" bson:"name"`
	Address string        `json:"address" bson:"address"`
	City    string        `json:"city" bson:"city"`
	Phone   string        `json:"phone" bson:"phone"`
}

// ShopRequest is the payload accepted when creating or updating a shop.
type ShopRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
	City    string `json:"city"`
	Phone   string `json:"phone"`
}
