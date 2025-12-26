package commonfun

import (
	"encoding/json"
	"log"
	"practice/api/models"
)

func SaveActivity(activityLog models.ActivityLog) {
	// convert struct → JSON
	data, _ := json.Marshal(activityLog)
	log.Println("data", string(data))
	// example SQL
	// INSERT INTO activity_logs (payload) VALUES (?)
	//_ = data
	// put db logic here
}
