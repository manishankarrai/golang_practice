package commonfun

import (
	"context"
	"log"
	"time"

	"practice/api/db"
	"practice/api/models"
)

const (
	// activityDB and activityCollection identify where request activity is stored.
	activityDB         = "golang_practice"
	activityCollection = "activities"
)

// SaveActivity persists an ActivityLog document, including the browser source
// captured from the Origin/Referer headers, into the activities collection.
func SaveActivity(activityLog models.ActivityLog) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := db.Database(activityDB).Collection(activityCollection).InsertOne(ctx, activityLog); err != nil {
		log.Println("failed to save activity:", err)
	}
}
