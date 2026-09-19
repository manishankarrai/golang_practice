package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Ride represents a customer's vehicle booking.
type Ride struct {
	ID        bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CustomerID bson.ObjectID `json:"customer_id" bson:"customer_id"`
	VehicleID  bson.ObjectID `json:"vehicle_id" bson:"vehicle_id"`
	Latitude   float64      `json:"latitude" bson:"latitude"`
	Longitude  float64      `json:"longitude" bson:"longitude"`
	Status     string       `json:"status" bson:"status"`
	CreatedAt  time.Time    `json:"created_at" bson:"created_at"`
}

// BookRideRequest is accepted when a customer books a vehicle.
type BookRideRequest struct {
	VehicleID string  `json:"vehicle_id" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required,gte=-90,lte=90"`
	Longitude float64 `json:"longitude" binding:"required,gte=-180,lte=180"`
}

// AvailableVehicle contains a vehicle and its current location.
type AvailableVehicle struct {
	Vehicle  MotorVehicle  `json:"vehicle"`
	Location VehicleLocation `json:"location"`
	Distance float64        `json:"distance_km"`
}
