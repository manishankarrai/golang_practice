package models

import "go.mongodb.org/mongo-driver/v2/bson"

// MotorPart represents a replacement or accessory part with its current price.
type MotorPart struct {
	ID          bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string        `json:"name" bson:"name"`
	Description string        `json:"description" bson:"description"`
	Price       float64       `json:"price" bson:"price"`
}

// MotorPartRequest is the payload accepted when creating or updating a motor part.
type MotorPartRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"gte=0"`
}
