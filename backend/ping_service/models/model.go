package models

type PingRequest struct {
	Site string `json:"site"`
}

type PingResponse struct {
	PingTime     string `json:"ping_time"`
	ResponseTime int64  `json:"response_time"`
	Error        string `json:"error,omitempty"`
}
