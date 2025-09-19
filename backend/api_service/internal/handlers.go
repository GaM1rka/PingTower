package internal

import (
	"api_service/configs"
	"api_service/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) AuthHandler(resp http.ResponseWriter, req *http.Request) {
	configs.APILogger.Println("Request for authentication.")

	body, err := io.ReadAll(req.Body)
	if err != nil {
		configs.APILogger.Println("read body failed:", err)
		http.Error(resp, "invalid request body", http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	var authReq models.AuthReq
	if err := json.Unmarshal(body, &authReq); err != nil {
		configs.APILogger.Println("json unmarshal failed:", err)
		http.Error(resp, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Сделать провреку на существование юзера, если такой есть, то банить запрос.

	// Добавить пользователя в БД.

	jwtReq, err := http.NewRequest(http.MethodPost, configs.JWTURL, req.Body)
	if err != nil {
		configs.APILogger.Println(err.Error())
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	jwtReq.Header.Set("Content-Type", "application/json")

	jwtResp, err := configs.Client.Do(jwtReq)
	if err != nil {
		configs.APILogger.Println(err.Error())
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	defer jwtResp.Body.Close()

	if jwtResp.StatusCode < 200 || jwtResp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(jwtResp.Body)
		configs.APILogger.Println("jwt service error:", jwtResp.Status, string(respBody))
		http.Error(resp, fmt.Sprintf("jwt error: %s", http.StatusText(jwtResp.StatusCode)), http.StatusBadGateway)
		return
	}

	var jwtBody models.AuthResp
	dec := json.NewDecoder(jwtResp.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&jwtBody); err != nil {
		configs.APILogger.Println("decode jwt json failed:", err)
		http.Error(resp, "invalid jwt response", http.StatusBadGateway)
		return
	}
	if jwtBody.Token == "" {
		configs.APILogger.Println("jwt token is empty")
		http.Error(resp, "invalid jwt response", http.StatusBadGateway)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(resp).Encode(jwtBody); err != nil {
		configs.APILogger.Println("encode response failed:", err)
	}
}
