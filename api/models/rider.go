package models

import "go.mongodb.org/mongo-driver/v2/bson"

// Rider represents a person who can currently ride a vehicle.
type Rider struct {
	ID    bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name  string        `json:"name" bson:"name"`
	Email string        `json:"email" bson:"email"`
	Phone string        `json:"phone" bson:"phone"`
}

// RiderRequest is the payload accepted when creating or updating a rider.
type RiderRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// AttachRiderRequest is used to attach a rider to a vehicle.
type AttachRiderRequest struct {
	RiderID string `json:"rider_id" binding:"required"`
}
