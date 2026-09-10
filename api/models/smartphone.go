package models

import "go.mongodb.org/mongo-driver/v2/bson"

// Smartphone represents a single mobile phone product.
type Smartphone struct {
	ID       bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string        `json:"name" bson:"name"`
	Brand    string        `json:"brand" bson:"brand"`
	Category string        `json:"category" bson:"category"`
	Price    float64       `json:"price" bson:"price"`
	Stock    int           `json:"stock" bson:"stock"`
}

// SmartphoneRequest is the payload accepted when creating or updating a smartphone.
type SmartphoneRequest struct {
	Name     string  `json:"name" binding:"required"`
	Brand    string  `json:"brand" binding:"required"`
	Category string  `json:"category"`
	Price    float64 `json:"price" binding:"required,gt=0"`
	Stock    int     `json:"stock" binding:"gte=0"`
}
