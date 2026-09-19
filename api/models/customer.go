package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Customer represents a person who can book rides.
type Customer struct {
	ID        bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string        `json:"name" bson:"name"`
	Email     string        `json:"email" bson:"email"`
	Phone     string        `json:"phone" bson:"phone"`
	Latitude  *float64      `json:"latitude,omitempty" bson:"latitude,omitempty"`
	Longitude *float64      `json:"longitude,omitempty" bson:"longitude,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

// CustomerRequest is accepted when creating or updating a customer.
type CustomerRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// CustomerLocationRequest updates the current customer's GPS location.
type CustomerLocationRequest struct {
	Latitude  float64 `json:"latitude" binding:"required,gte=-90,lte=90"`
	Longitude float64 `json:"longitude" binding:"required,gte=-180,lte=180"`
}
