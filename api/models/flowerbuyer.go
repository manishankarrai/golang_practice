package models

import "go.mongodb.org/mongo-driver/v2/bson"

// FlowerBuyer represents a customer who purchases flowers from a flower shop.
type FlowerBuyer struct {
	BuyerID  bson.ObjectID `json:"buyerid,omitempty" bson:"buyerid,omitempty"`
	Name     string        `json:"name" bson:"name"`
	Email    string        `json:"email" bson:"email"`
	Phone    string        `json:"phone" bson:"phone"`
	Address  string        `json:"address" bson:"address"`
}

// FlowerBuyerRequest is the payload accepted when creating or updating a buyer.
type FlowerBuyerRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Phone   string `json:"phone" binding:"required"`
	Address string `json:"address" binding:"required"`
}
