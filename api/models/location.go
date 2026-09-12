package models

import "go.mongodb.org/mongo-driver/v2/bson"

// Location represents a single location.
type Location struct {
	ID        bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string        `json:"name" bson:"name"`
	Address   string        `json:"address" bson:"address"`
	City      string        `json:"city" bson:"city"`
	StateID   bson.ObjectID `json:"state_id,omitempty" bson:"state_id,omitempty"`
	State     string        `json:"state" bson:"state"`
	Country   string        `json:"country" bson:"country"`
	ZipCode   string        `json:"zip_code" bson:"zip_code"`
	Latitude  float64       `json:"latitude" bson:"latitude"`
	Longitude float64       `json:"longitude" bson:"longitude"`
}

// LocationRequest is the payload accepted when creating or updating a location.
type LocationRequest struct {
	Name      string  `json:"name" binding:"required"`
	Address   string  `json:"address" binding:"required"`
	City      string  `json:"city" binding:"required"`
	StateID   string  `json:"state_id"`
	State     string  `json:"state"`
	Country   string  `json:"country" binding:"required"`
	ZipCode   string  `json:"zip_code"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
