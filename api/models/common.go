package models

import "time"

type ActivityLog struct {
	Method              string              `json:"method"`
	Path                string              `json:"path"`
	Query               string              `json:"query"`
	IP                  string              `json:"ip"`
	UserAgent           string              `json:"user_agent"`
	Authorization       string              `json:"authorization" bson:"authorization"`
	Origin              string              `json:"origin" bson:"origin"`
	Referer             string              `json:"referer" bson:"referer"`
	Headers             map[string][]string `json:"headers"`
	RequestBody         string              `json:"request_body"`
	StatusCode          int                 `json:"status_code"`
	ResponseBody        string              `json:"response_body"`
	LatencySeconds      float64             `json:"latency_seconds" bson:"latency_seconds"`
	IncomingRequestTime time.Time           `json:"incoming_request_time" bson:"incoming_request_time"`
	OutgoingRequestTime time.Time           `json:"outgoing_request_time" bson:"outgoing_request_time"`
}
