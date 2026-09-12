package models

import "go.mongodb.org/mongo-driver/v2/bson"

// Country represents a single country.
type Country struct {
	ID   bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name string        `json:"name" bson:"name"`
	Code string        `json:"code" bson:"code"`
}

// CountryRequest is the payload accepted when creating or updating a country.
type CountryRequest struct {
	Name string `json:"name" binding:"required"`
	Code string `json:"code"`
}
