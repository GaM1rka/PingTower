package models

type AuthReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResp struct {
	Token string `json:"token"`
}
