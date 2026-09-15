package models

import "go.mongodb.org/mongo-driver/v2/bson"

// Fruit represents a fruit product.
type Fruit struct {
	FruitID  bson.ObjectID `json:"fruitid,omitempty" bson:"fruitid,omitempty"`
	Name     string        `json:"name" bson:"name"`
	Type     string        `json:"type" bson:"type"`
	Location string        `json:"location" bson:"location"`
	Pricing  float64       `json:"pricing" bson:"pricing"`
	Details  string        `json:"details" bson:"details"`
}

// FruitRequest is the payload accepted when creating or updating a fruit.
type FruitRequest struct {
	Name     string  `json:"name" binding:"required"`
	Type     string  `json:"type" binding:"required"`
	Location string  `json:"location" binding:"required"`
	Pricing  float64 `json:"pricing" binding:"required,gt=0"`
	Details  string  `json:"details"`
}
