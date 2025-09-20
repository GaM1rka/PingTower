package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	store *Storage
}

func NewHandler(store *Storage) *Handler {
	return &Handler{store: store}
}

// UserHandler добавляет пользователя
func (h *Handler) UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var user struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		userID, err := h.store.CreateUser(user.Email, user.Password)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating user: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{"id": userID})

	case http.MethodGet:
		email := r.URL.Query().Get("email")
		if email == "" {
			http.Error(w, "email query param required", http.StatusBadRequest)
			return
		}
		id, err := h.store.GetUserIDByEmail(email)
		if err != nil {
			// не нашли — 404
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"id": id, "email": email})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GET /all-users-sites
func (h *Handler) AllUsersSitesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	items, err := h.store.GetAllUsersSites()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting users sites: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// POST /ping  {user_id, site, resp_time, status}
func (h *Handler) PingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		UserID   int    `json:"user_id"`
		Site     string `json:"site"`
		RespTime int64  `json:"resp_time"`
		Status   string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if body.UserID == 0 || body.Site == "" {
		http.Error(w, "user_id and site are required", http.StatusBadRequest)
		return
	}
	if err := h.store.AddPingLog(body.UserID, body.Site, body.RespTime, body.Status); err != nil {
		http.Error(w, fmt.Sprintf("Error saving ping log: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"message": "ping saved"})
}

// GET /user/{id}/email
func (h *Handler) UserEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	// ожидаем: ["user", "{id}", "email"]
	if len(parts) != 3 || parts[0] != "user" || parts[2] != "email" {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(parts[1])
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}
	email, err := h.store.GetUserEmail(userID)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"email": email})
}

// UserSitesHandler добавляет сайт пользователю
func (h *Handler) UserSitesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем ID из URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(pathParts[len(pathParts)-1])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var siteData struct {
		Site string `json:"site"`
	}

	if err := json.NewDecoder(r.Body).Decode(&siteData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = h.store.AddUserSite(userID, siteData.Site)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error adding site: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Site added successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CheckerHandler получает логи конкретного сайта
func (h *Handler) CheckerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	siteID, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		http.Error(w, "Invalid site ID", http.StatusBadRequest)
		return
	}

	// сначала попробуем query ?user_id=
	var userID int
	if uid := r.URL.Query().Get("user_id"); uid != "" {
		userID, err = strconv.Atoi(uid)
		if err != nil {
			http.Error(w, "Invalid user_id", http.StatusBadRequest)
			return
		}
	} else {
		// fallback: тело JSON (на случай, если кто-то всё же шлёт body)
		var userData struct {
			UserID int `json:"user_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
			http.Error(w, "Invalid JSON (need user_id)", http.StatusBadRequest)
			return
		}
		userID = userData.UserID
	}

	logs, err := h.store.GetSiteLogs(userID, siteID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting logs: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

// CheckersHandler обрабатывает GET и POST для /checkers
func (h *Handler) CheckersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getUserLogs(w, r)
	case http.MethodPost:
		h.addSiteWithCheck(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) getUserLogs(w http.ResponseWriter, r *http.Request) {
	var userID int
	if uid := r.URL.Query().Get("user_id"); uid != "" {
		id, err := strconv.Atoi(uid)
		if err != nil {
			http.Error(w, "Invalid user_id", http.StatusBadRequest)
			return
		}
		userID = id
	} else {
		var userData struct {
			UserID int `json:"user_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
			http.Error(w, "Invalid JSON (need user_id)", http.StatusBadRequest)
			return
		}
		userID = userData.UserID
	}

	logs, err := h.store.GetAllUserLogs(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting logs: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func (h *Handler) addSiteWithCheck(w http.ResponseWriter, r *http.Request) {
	var siteData struct {
		Site   string `json:"site"`
		Time   int    `json:"time"`
		UserID int    `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&siteData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Добавляем сайт и сразу создаем запись для проверки
	err := h.store.AddSiteWithCheck(siteData.UserID, siteData.Site, siteData.Time)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error adding site: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Site and check added successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
