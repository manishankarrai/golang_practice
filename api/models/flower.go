package models

import "go.mongodb.org/mongo-driver/v2/bson"

// Flower represents a flower product.
type Flower struct {
	FlowerID  bson.ObjectID `json:"flowerid,omitempty" bson:"flowerid,omitempty"`
	Name      string        `json:"name" bson:"name"`
	Type      string        `json:"type" bson:"type"`
	Color     string        `json:"color" bson:"color"`
	Location  string        `json:"location" bson:"location"`
	Pricing   float64       `json:"pricing" bson:"pricing"`
	Details   string        `json:"details" bson:"details"`
}

// FlowerRequest is the payload accepted when creating or updating a flower.
type FlowerRequest struct {
	Name     string  `json:"name" binding:"required"`
	Type     string  `json:"type" binding:"required"`
	Color    string  `json:"color" binding:"required"`
	Location string  `json:"location" binding:"required"`
	Pricing  float64 `json:"pricing" binding:"required,gt=0"`
	Details  string  `json:"details"`
}
