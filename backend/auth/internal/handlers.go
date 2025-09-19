package internal

import (
	"auth/configs"
	"auth/models"
	"encoding/json"
	"net/http"
	"strings"
)

func GenerateTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req models.GenerateTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		configs.Logger.Println("decode generate token request failed:", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password required", http.StatusBadRequest)
		return
	}

	accessToken, err := GenerateToken(req.Email)
	if err != nil {
		configs.Logger.Println("generate access token failed:", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	refreshToken, err := GenerateRefreshToken(req.Email)
	if err != nil {
		configs.Logger.Println("generate refresh token failed:", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	resp := models.GenerateTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func ValidateHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		var req models.ValidateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			configs.Logger.Println("no auth header and decode failed:", err)
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		authHeader = "Bearer " + req.Token
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		resp := models.ValidateResponse{Valid: false}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := ValidateToken(token)
	if err != nil {
		configs.Logger.Println("validate token failed:", err)
		resp := models.ValidateResponse{Valid: false}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := models.ValidateResponse{
		Valid: true,
		Email: claims.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		configs.Logger.Println("decode refresh request failed:", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.RefreshToken == "" {
		http.Error(w, "refresh_token required", http.StatusBadRequest)
		return
	}

	email, err := ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		configs.Logger.Println("validate refresh token failed:", err)
		http.Error(w, "invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	if err := RevokeRefreshToken(req.RefreshToken); err != nil {
		configs.Logger.Println("revoke refresh token failed:", err)
	}

	newAccessToken, err := GenerateToken(email)
	if err != nil {
		configs.Logger.Println("generate new access token failed:", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	newRefreshToken, err := GenerateRefreshToken(email)
	if err != nil {
		configs.Logger.Println("generate new refresh token failed:", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	resp := models.RefreshResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}