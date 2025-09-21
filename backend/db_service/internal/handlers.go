package internal

import (
	"db_service/configs"
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
	configs.DBLogger.Printf("➡️ UserHandler %s %s", r.Method, r.URL.String())

	switch r.Method {
	case http.MethodPost:
		var user struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			configs.DBLogger.Println("❌ UserHandler POST: decode error:", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		configs.DBLogger.Printf("📥 CreateUser email=%s", user.Email)

		userID, err := h.store.CreateUser(user.Email, user.Password)
		if err != nil {
			configs.DBLogger.Println("❌ UserHandler POST: CreateUser error:", err)
			http.Error(w, fmt.Sprintf("Error creating user: %v", err), http.StatusInternalServerError)
			return
		}
		configs.DBLogger.Printf("✅ User created id=%d email=%s", userID, user.Email)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{"id": userID})

	case http.MethodGet:
		email := r.URL.Query().Get("email")
		if email == "" {
			configs.DBLogger.Println("❌ UserHandler GET: missing email query param")
			http.Error(w, "email query param required", http.StatusBadRequest)
			return
		}
		configs.DBLogger.Printf("📥 GetUserIDByEmail email=%s", email)

		id, err := h.store.GetUserIDByEmail(email)
		if err != nil {
			configs.DBLogger.Printf("❌ UserHandler GET: user not found email=%s", email)
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		configs.DBLogger.Printf("✅ Found user id=%d email=%s", id, email)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"id": id, "email": email})

	default:
		configs.DBLogger.Printf("❌ UserHandler: method not allowed %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GET /all-users-sites
func (h *Handler) AllUsersSitesHandler(w http.ResponseWriter, r *http.Request) {
	configs.DBLogger.Printf("➡️ AllUsersSitesHandler %s %s", r.Method, r.URL.String())

	if r.Method != http.MethodGet {
		configs.DBLogger.Printf("❌ AllUsersSitesHandler: method not allowed %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	items, err := h.store.GetAllUsersSites()
	if err != nil {
		configs.DBLogger.Println("❌ AllUsersSitesHandler: GetAllUsersSites error:", err)
		http.Error(w, fmt.Sprintf("Error getting users sites: %v", err), http.StatusInternalServerError)
		return
	}
	configs.DBLogger.Printf("✅ AllUsersSitesHandler: users=%d", len(items))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// POST /ping  {user_id, site, resp_time, status}
func (h *Handler) PingHandler(w http.ResponseWriter, r *http.Request) {
	configs.DBLogger.Printf("➡️ PingHandler %s %s", r.Method, r.URL.String())

	if r.Method != http.MethodPost {
		configs.DBLogger.Printf("❌ PingHandler: method not allowed %s", r.Method)
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
		configs.DBLogger.Println("❌ PingHandler: decode error:", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	configs.DBLogger.Printf("📥 PingHandler payload user_id=%d site=%s resp_time=%d status=%s",
		body.UserID, body.Site, body.RespTime, body.Status)

	if body.UserID == 0 || body.Site == "" {
		configs.DBLogger.Println("❌ PingHandler: missing user_id or site")
		http.Error(w, "user_id and site are required", http.StatusBadRequest)
		return
	}
	if err := h.store.AddPingLog(body.UserID, body.Site, body.RespTime, body.Status); err != nil {
		configs.DBLogger.Println("❌ PingHandler: AddPingLog error:", err)
		http.Error(w, fmt.Sprintf("Error saving ping log: %v", err), http.StatusInternalServerError)
		return
	}
	configs.DBLogger.Println("✅ PingHandler: log saved")

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"message": "ping saved"})
}

// GET /user/{id}/email
func (h *Handler) UserEmailHandler(w http.ResponseWriter, r *http.Request) {
	configs.DBLogger.Printf("➡️ UserEmailHandler %s %s", r.Method, r.URL.String())

	if r.Method != http.MethodGet {
		configs.DBLogger.Printf("❌ UserEmailHandler: method not allowed %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 3 || parts[0] != "user" || parts[2] != "email" {
		configs.DBLogger.Printf("❌ UserEmailHandler: invalid URL %s", r.URL.Path)
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(parts[1])
	if err != nil {
		configs.DBLogger.Printf("❌ UserEmailHandler: invalid user id %q", parts[1])
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}
	configs.DBLogger.Printf("📥 UserEmailHandler: get email for user_id=%d", userID)

	email, err := h.store.GetUserEmail(userID)
	if err != nil {
		configs.DBLogger.Println("❌ UserEmailHandler: user not found:", err)
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	configs.DBLogger.Printf("✅ UserEmailHandler: email=%s", email)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"email": email})
}

// UserSitesHandler добавляет сайт пользователю
func (h *Handler) UserSitesHandler(w http.ResponseWriter, r *http.Request) {
	configs.DBLogger.Printf("➡️ UserSitesHandler %s %s", r.Method, r.URL.String())

	if r.Method != http.MethodPost {
		configs.DBLogger.Printf("❌ UserSitesHandler: method not allowed %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем ID из URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		configs.DBLogger.Printf("❌ UserSitesHandler: invalid URL %s", r.URL.Path)
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(pathParts[len(pathParts)-1])
	if err != nil {
		configs.DBLogger.Printf("❌ UserSitesHandler: invalid user ID %q", pathParts[len(pathParts)-1])
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var siteData struct {
		Site string `json:"site"`
	}

	if err := json.NewDecoder(r.Body).Decode(&siteData); err != nil {
		configs.DBLogger.Println("❌ UserSitesHandler: decode error:", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	configs.DBLogger.Printf("📥 UserSitesHandler: add site user_id=%d site=%s", userID, siteData.Site)

	err = h.store.AddUserSite(userID, siteData.Site)
	if err != nil {
		configs.DBLogger.Println("❌ UserSitesHandler: AddUserSite error:", err)
		http.Error(w, fmt.Sprintf("Error adding site: %v", err), http.StatusInternalServerError)
		return
	}
	configs.DBLogger.Println("✅ UserSitesHandler: site added")

	response := map[string]string{"message": "Site added successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CheckerHandler получает логи конкретного сайта
func (h *Handler) CheckerHandler(w http.ResponseWriter, r *http.Request) {
	configs.DBLogger.Printf("➡️ CheckerHandler %s %s", r.Method, r.URL.String())

	if r.Method != http.MethodGet {
		configs.DBLogger.Printf("❌ CheckerHandler: method not allowed %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		configs.DBLogger.Printf("❌ CheckerHandler: invalid URL %s", r.URL.Path)
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	siteID, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		configs.DBLogger.Printf("❌ CheckerHandler: invalid site ID %q", parts[len(parts)-1])
		http.Error(w, "Invalid site ID", http.StatusBadRequest)
		return
	}

	// сначала попробуем query ?user_id=
	var userID int
	if uid := r.URL.Query().Get("user_id"); uid != "" {
		userID, err = strconv.Atoi(uid)
		if err != nil {
			configs.DBLogger.Printf("❌ CheckerHandler: invalid user_id query %q", uid)
			http.Error(w, "Invalid user_id", http.StatusBadRequest)
			return
		}
	} else {
		// fallback: тело JSON (на случай, если кто-то всё же шлёт body)
		var userData struct {
			UserID int `json:"user_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
			configs.DBLogger.Println("❌ CheckerHandler: decode body user_id error:", err)
			http.Error(w, "Invalid JSON (need user_id)", http.StatusBadRequest)
			return
		}
		userID = userData.UserID
	}
	configs.DBLogger.Printf("📥 CheckerHandler: user_id=%d site_id=%d", userID, siteID)

	logs, err := h.store.GetSiteLogs(userID, siteID)
	if err != nil {
		configs.DBLogger.Println("❌ CheckerHandler: GetSiteLogs error:", err)
		http.Error(w, fmt.Sprintf("Error getting logs: %v", err), http.StatusInternalServerError)
		return
	}
	configs.DBLogger.Printf("✅ CheckerHandler: logs=%d", len(logs))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

// CheckersHandler обрабатывает GET и POST для /checkers
func (h *Handler) CheckersHandler(w http.ResponseWriter, r *http.Request) {
	configs.DBLogger.Printf("➡️ CheckersHandler %s %s", r.Method, r.URL.String())

	switch r.Method {
	case http.MethodGet:
		h.getUserLogs(w, r)
	case http.MethodPost:
		h.addSiteWithCheck(w, r)
	default:
		configs.DBLogger.Printf("❌ CheckersHandler: method not allowed %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) getUserLogs(w http.ResponseWriter, r *http.Request) {
	// Логироваться будет из вызывающего CheckersHandler
	var userID int
	if uid := r.URL.Query().Get("user_id"); uid != "" {
		id, err := strconv.Atoi(uid)
		if err != nil {
			configs.DBLogger.Printf("❌ getUserLogs: invalid user_id query %q", uid)
			http.Error(w, "Invalid user_id", http.StatusBadRequest)
			return
		}
		userID = id
	} else {
		var userData struct {
			UserID int `json:"user_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
			configs.DBLogger.Println("❌ getUserLogs: decode user_id error:", err)
			http.Error(w, "Invalid JSON (need user_id)", http.StatusBadRequest)
			return
		}
		userID = userData.UserID
	}
	configs.DBLogger.Printf("📥 getUserLogs: user_id=%d", userID)

	logs, err := h.store.GetAllUserLogs(userID)
	if err != nil {
		configs.DBLogger.Println("❌ getUserLogs: GetAllUserLogs error:", err)
		http.Error(w, fmt.Sprintf("Error getting logs: %v", err), http.StatusInternalServerError)
		return
	}
	configs.DBLogger.Printf("✅ getUserLogs: logs=%d", len(logs))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func (h *Handler) addSiteWithCheck(w http.ResponseWriter, r *http.Request) {
	// Логироваться будет из вызывающего CheckersHandler
	var siteData struct {
		Site   string `json:"site"`
		Time   int    `json:"time"`
		UserID int    `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&siteData); err != nil {
		configs.DBLogger.Println("❌ addSiteWithCheck: decode error:", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	configs.DBLogger.Printf("📥 addSiteWithCheck: user_id=%d site=%s time=%d", siteData.UserID, siteData.Site, siteData.Time)

	// Добавляем сайт и сразу создаем запись для проверки
	err := h.store.AddSiteWithCheck(siteData.UserID, siteData.Site, siteData.Time)
	if err != nil {
		configs.DBLogger.Println("❌ addSiteWithCheck: AddSiteWithCheck error:", err)
		http.Error(w, fmt.Sprintf("Error adding site: %v", err), http.StatusInternalServerError)
		return
	}
	configs.DBLogger.Println("✅ addSiteWithCheck: done")

	response := map[string]string{"message": "Site and check added successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
