package models

import "go.mongodb.org/mongo-driver/v2/bson"

// State represents a state/region within a country (e.g. Bihar, Goa in India).
type State struct {
	ID      bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name    string        `json:"name" bson:"name"`
	Code    string        `json:"code" bson:"code"`
	Country string        `json:"country" bson:"country"`
}

// StateRequest is the payload accepted when creating or updating a state.
type StateRequest struct {
	Name    string `json:"name" binding:"required"`
	Code    string `json:"code"`
	Country string `json:"country" binding:"required"`
}
