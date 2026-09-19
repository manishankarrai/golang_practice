package models

import "go.mongodb.org/mongo-driver/v2/bson"

// OwnerAddress stores the address supplied when a vehicle is registered.
type OwnerAddress struct {
	Address  string `json:"address" bson:"address"`
	City     string `json:"city" bson:"city"`
	State    string `json:"state" bson:"state"`
	Country  string `json:"country" bson:"country"`
	ZipCode  string `json:"zip_code" bson:"zip_code"`
}

// Registration links a vehicle to its owner and registration address.
type Registration struct {
	ID                 bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	VehicleID          bson.ObjectID `json:"vehicle_id" bson:"vehicle_id"`
	OwnerName          string        `json:"owner_name" bson:"owner_name"`
	OwnerEmail         string        `json:"owner_email" bson:"owner_email"`
	OwnerPhone         string        `json:"owner_phone" bson:"owner_phone"`
	RegistrationNumber string        `json:"registration_number" bson:"registration_number"`
	RegistrationDate   string        `json:"registration_date" bson:"registration_date"`
	ExpiryDate         string        `json:"expiry_date" bson:"expiry_date"`
	OwnerAddress       OwnerAddress  `json:"owner_address" bson:"owner_address"`
}

// RegistrationRequest is the payload accepted when creating or updating a registration.
type RegistrationRequest struct {
	VehicleID          string       `json:"vehicle_id" binding:"required"`
	OwnerName          string       `json:"owner_name" binding:"required"`
	OwnerEmail         string       `json:"owner_email"`
	OwnerPhone         string       `json:"owner_phone"`
	RegistrationNumber string       `json:"registration_number" binding:"required"`
	RegistrationDate   string       `json:"registration_date" binding:"required"`
	ExpiryDate         string       `json:"expiry_date"`
	OwnerAddress       OwnerAddress `json:"owner_address" binding:"required"`
}
