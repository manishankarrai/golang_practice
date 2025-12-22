package models

type ActivityLog struct {
	Method       string              `json:"method"`
	Path         string              `json:"path"`
	Query        string              `json:"query"`
	IP           string              `json:"ip"`
	UserAgent    string              `json:"user_agent"`
	Headers      map[string][]string `json:"headers"`
	RequestBody  string              `json:"request_body"`
	StatusCode   int                 `json:"status_code"`
	ResponseBody string              `json:"response_body"`
	LatencyMs    int64               `json:"latency_ms"`
}
