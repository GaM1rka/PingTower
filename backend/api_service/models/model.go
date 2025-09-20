package models

type AuthReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ID           int    `json:"id,omitempty"`
	Email        string `json:"email,omitempty"`
}

type Site struct {
	ID            int    `json:"id"`
	URL           string `json:"url"` // ВАЖНО: именно "url"
	CheckInterval int    `json:"check_interval"`
}

type UserSites struct {
	UserID int    `json:"user_id"`
	Sites  []Site `json:"sites"`
}

type PingLog struct {
	ReqTime  string `json:"req_time"`
	RespTime int64  `json:"resp_time"`
	Status   string `json:"status"`
	Site     string `json:"site"`
}

type CheckerRequest struct {
	UserID int `json:"user_id"`
}

type AddSiteRequest struct {
	UserID int    `json:"user_id"`
	Site   string `json:"site"`
	Time   int    `json:"time,omitempty"`
}

type Notification struct {
	Email string `json:"email"`
	Site  string `json:"site"`
	Time  int64  `json:"time"`
}

type PingRequest struct {
	Site string `json:"site"`
}

type PingResponse struct {
	ResponseTime int64  `json:"response_time"`
	Status       string `json:"status"`
}
