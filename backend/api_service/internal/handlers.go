package internal

import (
	"api_service/configs"
	"api_service/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

	// Валидация email
	if !isValidEmail(authReq.Email) {
		http.Error(resp, "invalid email format", http.StatusBadRequest)
		return
	}

	// Валидация пароля
	if len(authReq.Password) < 6 {
		http.Error(resp, "password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	// Проверка существования пользователя
	userExists, err := h.checkUserExists(authReq.Email)
	if err != nil {
		configs.APILogger.Println("check user exists failed:", err)
		http.Error(resp, "internal error", http.StatusInternalServerError)
		return
	}

	if req.URL.Path == "/register" {
		if userExists {
			http.Error(resp, "user already exists", http.StatusConflict)
			return
		}

		// Регистрация пользователя
		userID, err := h.registerUser(authReq.Email, authReq.Password)
		if err != nil {
			configs.APILogger.Println("register user failed:", err)
			http.Error(resp, "registration failed", http.StatusInternalServerError)
			return
		}

		// ✅ Генерация JWT (передаём пароль третьим аргументом)
		jwtBody, err := h.getJWTToken(authReq.Email, userID, authReq.Password)
		if err != nil {
			configs.APILogger.Println("get jwt token failed:", err)
			http.Error(resp, "jwt generation failed", http.StatusInternalServerError)
			return
		}

		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusCreated)
		json.NewEncoder(resp).Encode(jwtBody)

	} else if req.URL.Path == "/login" {
		if !userExists {
			http.Error(resp, "user not found", http.StatusNotFound)
			return
		}

		// Проверка логина (пока заглушка)
		userID, err := h.verifyLogin(authReq.Email, authReq.Password)
		if err != nil {
			configs.APILogger.Println("login verification failed:", err)
			http.Error(resp, "login failed", http.StatusUnauthorized)
			return
		}

		// ✅ Генерация JWT (передаём пароль третьим аргументом)
		jwtBody, err := h.getJWTToken(authReq.Email, userID, authReq.Password)
		if err != nil {
			configs.APILogger.Println("get jwt token failed:", err)
			http.Error(resp, "jwt generation failed", http.StatusInternalServerError)
			return
		}

		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		json.NewEncoder(resp).Encode(jwtBody)
	}
}

func (h *Handler) CheckerHandler(resp http.ResponseWriter, req *http.Request) {
	// Проверка JWT
	userID, err := h.verifyJWT(req)
	if err != nil {
		http.Error(resp, "unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Извлекаем ID сайта из URL
	pathParts := strings.Split(req.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(resp, "invalid URL", http.StatusBadRequest)
		return
	}

	siteID, err := strconv.Atoi(pathParts[len(pathParts)-1])
	if err != nil {
		http.Error(resp, "invalid site ID", http.StatusBadRequest)
		return
	}

	// Получаем логи сайта
	logs, err := h.getSiteLogs(userID, siteID)
	if err != nil {
		configs.APILogger.Println("get site logs failed:", err)
		http.Error(resp, "internal error", http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resp).Encode(logs)
}

func (h *Handler) CheckersHandler(resp http.ResponseWriter, req *http.Request) {
	// Проверка JWT
	userID, err := h.verifyJWT(req)
	if err != nil {
		http.Error(resp, "unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	switch req.Method {
	case http.MethodGet:
		h.getUserSites(resp, userID)
	case http.MethodPost:
		h.addUserSite(resp, req, userID)
	default:
		http.Error(resp, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) PingAllHandler(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(resp, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем всех пользователей и их сайты
	usersSites, err := h.getAllUsersSites()
	if err != nil {
		configs.APILogger.Println("get all users sites failed:", err)
		http.Error(resp, "internal error", http.StatusInternalServerError)
		return
	}

	successCount := 0
	failCount := 0

	// Пингуем все сайты
	for _, userSite := range usersSites {
		for _, site := range userSite.Sites {
			pingResult, err := h.pingSite(site.URL)
			if err != nil {
				configs.APILogger.Printf("ping site %s failed: %v", site.URL, err)
				failCount++
				continue
			}

			// Сохраняем лог
			err = h.savePingLog(userSite.UserID, site.URL, pingResult.ResponseTime, pingResult.Status)
			if err != nil {
				configs.APILogger.Printf("save ping log failed: %v", err)
				failCount++
				continue
			}

			// Если статус нерабочий, отправляем уведомление
			if pingResult.Status != "success" && pingResult.Status != "ok" {
				userEmail, err := h.getUserEmail(userSite.UserID)
				if err != nil {
					configs.APILogger.Printf("get user email failed: %v", err)
					continue
				}

				err = configs.SendKafkaNotification(userEmail, site.URL, pingResult.ResponseTime)
				if err != nil {
					configs.APILogger.Printf("send kafka notification failed: %v", err)
				}
			}

			successCount++
		}
	}

	response := map[string]interface{}{
		"message":    "Ping completed",
		"successful": successCount,
		"failed":     failCount,
		"total":      successCount + failCount,
	}

	resp.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resp).Encode(response)
}

// Вспомогательные методы
func (h *Handler) checkUserExists(email string) (bool, error) {
	req, err := http.NewRequest(http.MethodGet, configs.DBURL+"/user?email="+email, nil)
	if err != nil {
		return false, err
	}

	resp, err := configs.Client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

func (h *Handler) registerUser(email, password string) (int, error) {
	userData := map[string]string{
		"email":    email,
		"password": password,
	}

	jsonData, _ := json.Marshal(userData)
	req, err := http.NewRequest(http.MethodPost, configs.DBURL+"/user", bytes.NewReader(jsonData))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := configs.Client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return 0, fmt.Errorf("registration failed")
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return int(result["id"].(float64)), nil
}

func (h *Handler) verifyLogin(email, password string) (int, error) {
	// Реализация проверки логина через auth service
	return 1, nil // Заглушка
}

func (h *Handler) getJWTToken(email string, userID int, password string) (*models.AuthResp, error) {
	authData := models.AuthReq{
		Email:    email,
		Password: password, // ВАЖНО: не пустой
	}
	jsonData, _ := json.Marshal(authData)

	req, err := http.NewRequest(http.MethodPost, configs.JWTURL, bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := configs.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("auth service error: status=%d body=%s", resp.StatusCode, string(b))
	}

	var jwtBody models.AuthResp // структура должна совпадать с JSON от auth_service (/generate)
	if err := json.NewDecoder(resp.Body).Decode(&jwtBody); err != nil {
		return nil, err
	}

	// Допишем служебные поля
	jwtBody.ID = userID
	jwtBody.Email = email
	return &jwtBody, nil
}

func (h *Handler) verifyJWT(req *http.Request) (int, error) {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return 0, fmt.Errorf("missing authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil // Заменить на реальный секрет (или читать из ENV)
	})

	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	claims := token.Claims.(jwt.MapClaims)
	return int(claims["user_id"].(float64)), nil
}

func (h *Handler) getSiteLogs(userID, siteID int) ([]models.PingLog, error) {
	req, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf("%s/checker/%d?user_id=%d", configs.DBURL, siteID, userID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := configs.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var logs []models.PingLog
	err = json.NewDecoder(resp.Body).Decode(&logs)
	return logs, err
}

func (h *Handler) getUserSites(resp http.ResponseWriter, userID int) {
	req, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf("%s/checkers?user_id=%d", configs.DBURL, userID), nil)
	if err != nil {
		http.Error(resp, "internal error", http.StatusInternalServerError)
		return
	}

	sitesResp, err := configs.Client.Do(req)
	if err != nil {
		http.Error(resp, "internal error", http.StatusInternalServerError)
		return
	}
	defer sitesResp.Body.Close()

	if sitesResp.StatusCode != http.StatusOK {
		http.Error(resp, "failed to get sites", http.StatusInternalServerError)
		return
	}

	io.Copy(resp, sitesResp.Body)
}

func (h *Handler) addUserSite(resp http.ResponseWriter, req *http.Request, userID int) {
	var siteReq models.AddSiteRequest
	if err := json.NewDecoder(req.Body).Decode(&siteReq); err != nil {
		http.Error(resp, "invalid JSON", http.StatusBadRequest)
		return
	}

	siteReq.UserID = userID
	jsonData, _ := json.Marshal(siteReq)

	dbReq, err := http.NewRequest(http.MethodPost, configs.DBURL+"/checkers", bytes.NewReader(jsonData))
	if err != nil {
		http.Error(resp, "internal error", http.StatusInternalServerError)
		return
	}
	dbReq.Header.Set("Content-Type", "application/json")

	dbResp, err := configs.Client.Do(dbReq)
	if err != nil {
		http.Error(resp, "internal error", http.StatusInternalServerError)
		return
	}
	defer dbResp.Body.Close()

	resp.WriteHeader(dbResp.StatusCode)
	io.Copy(resp, dbResp.Body)
}

func (h *Handler) getAllUsersSites() ([]models.UserSites, error) {
	// Запрос к DB service для получения всех пользователей и их сайтов
	req, err := http.NewRequest(http.MethodGet, configs.DBURL+"/all-users-sites", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := configs.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get users sites: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("db service returned status: %d", resp.StatusCode)
	}

	var usersSites []models.UserSites
	if err := json.NewDecoder(resp.Body).Decode(&usersSites); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return usersSites, nil
}

func (h *Handler) getUserEmail(userID int) (string, error) {
	// Запрос к DB service для получения email пользователя
	req, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf("%s/user/%d/email", configs.DBURL, userID), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := configs.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get user email: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("db service returned status: %d", resp.StatusCode)
	}

	var result struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	return result.Email, nil
}

// internal/handlers.go
func (h *Handler) pingSite(url string) (*PingResult, error) {
	// логируем
	configs.APILogger.Printf("ping site: %s", url)

	pingRequest := models.PingRequest{Site: url}
	jsonData, err := json.Marshal(pingRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ping request: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, configs.PingURL+"/ping", bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create ping request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := configs.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ping request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ping service returned status: %d (%s)", resp.StatusCode, strings.TrimSpace(string(b)))
	}

	var pingResponse models.PingResponse
	if err := json.NewDecoder(resp.Body).Decode(&pingResponse); err != nil {
		return nil, fmt.Errorf("failed to parse ping response: %v", err)
	}

	if pingResponse.ResponseTime < 0 {
		return nil, fmt.Errorf("ping service signaled failure (response_time=%d)", pingResponse.ResponseTime)
	}

	// ping_service не шлёт status — подставим "ok" по умолчанию
	status := pingResponse.Status
	if status == "" {
		status = "ok"
	}

	return &PingResult{
		ResponseTime: pingResponse.ResponseTime,
		Status:       status,
	}, nil
}

func (h *Handler) savePingLog(userID int, site string, responseTime int64, status string) error {
	logData := map[string]interface{}{
		"user_id":   userID,
		"site":      site,
		"resp_time": responseTime,
		"status":    status,
	}

	jsonData, _ := json.Marshal(logData)
	req, err := http.NewRequest(http.MethodPost, configs.DBURL+"/ping", bytes.NewReader(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := configs.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to save ping log")
	}

	return nil
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	// Дополнительная проверка regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

type PingResult struct {
	ResponseTime int64
	Status       string
}
