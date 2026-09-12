package models

import "go.mongodb.org/mongo-driver/v2/bson"

// Stage represents a single stage.
type Stage struct {
	ID          bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string        `json:"name" bson:"name"`
	Description string        `json:"description" bson:"description"`
	Sequence    int           `json:"sequence" bson:"sequence"`
	Status      string        `json:"status" bson:"status"`
}

// StageRequest is the payload accepted when creating or updating a stage.
type StageRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Sequence    int    `json:"sequence"`
	Status      string `json:"status"`
}
