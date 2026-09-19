package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// VehicleLocation is the current GPS position of a vehicle and its rider.
type VehicleLocation struct {
	VehicleID string         `json:"vehicle_id" bson:"vehicle_id"`
	RiderID   *bson.ObjectID `json:"rider_id,omitempty" bson:"rider_id,omitempty"`
	Latitude  float64        `json:"latitude" bson:"latitude"`
	Longitude float64        `json:"longitude" bson:"longitude"`
	UpdatedAt time.Time      `json:"updated_at" bson:"updated_at"`
}

// VehicleLocationRequest is accepted by the vehicle location endpoint.
type VehicleLocationRequest struct {
	RiderID   string  `json:"rider_id" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required,gte=-90,lte=90"`
	Longitude float64 `json:"longitude" binding:"required,gte=-180,lte=180"`
}
