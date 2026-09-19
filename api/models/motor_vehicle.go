package models

import "go.mongodb.org/mongo-driver/v2/bson"

// MotorVehicle represents a vehicle that can be registered to an owner.
type MotorVehicle struct {
	ID          bson.ObjectID  `json:"id,omitempty" bson:"_id,omitempty"`
	Make        string         `json:"make" bson:"make"`
	Model       string         `json:"model" bson:"model"`
	Year        int            `json:"year" bson:"year"`
	Color       string         `json:"color" bson:"color"`
	VehicleType string         `json:"vehicle_type" bson:"vehicle_type"`
	VIN         string         `json:"vin" bson:"vin"`
	EngineNo    string         `json:"engine_no" bson:"engine_no"`
	RiderID     *bson.ObjectID `json:"rider_id,omitempty" bson:"rider_id,omitempty"`
}

// MotorVehicleRequest is the payload accepted when creating or updating a vehicle.
type MotorVehicleRequest struct {
	Make        string `json:"make" binding:"required"`
	Model       string `json:"model" binding:"required"`
	Year        int    `json:"year" binding:"required,gt=0"`
	Color       string `json:"color"`
	VehicleType string `json:"vehicle_type" binding:"required"`
	VIN         string `json:"vin" binding:"required"`
	EngineNo    string `json:"engine_no" binding:"required"`
}
