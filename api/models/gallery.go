package models

import "go.mongodb.org/mongo-driver/v2/bson"

// Gallery represents an uploaded gallery photo and its metadata.
type Gallery struct {
	ID          bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string        `json:"title" bson:"title"`
	Description string        `json:"description" bson:"description"`
	FileName    string        `json:"file_name" bson:"file_name"`
	URL         string        `json:"url" bson:"url"`
	ContentType string        `json:"content_type" bson:"content_type"`
	Size        int64         `json:"size" bson:"size"`
}
