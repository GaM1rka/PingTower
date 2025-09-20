package internal

import (
	"encoding/json"
	"log"
	"net/http"
	"notification_service/models"
	"strings"
)

type NotificationHandler struct {
	emailService *EmailService
}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{
		emailService: NewEmailService(),
	}
}

func (nh *NotificationHandler) writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}

func (nh *NotificationHandler) SendNotificationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		nh.writeJSON(w, http.StatusMethodNotAllowed, models.NotificationResponse{
			Status: "error",
			Error:  "method not allowed",
		})
		return
	}

	var req models.NotificationRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		nh.writeJSON(w, http.StatusBadRequest, models.NotificationResponse{
			Status: "error",
			Error:  "invalid JSON: " + err.Error(),
		})
		return
	}

	if err := nh.validateRequest(req); err != nil {
		nh.writeJSON(w, http.StatusBadRequest, models.NotificationResponse{
			Status: "error",
			Error:  err.Error(),
		})
		return
	}

	log.Printf("Processing notification request for email: %s, site: %s", req.Email, req.Site)

	response, err := nh.emailService.SendNotificationWithRetry(r.Context(), req, 3)
	if err != nil {
		log.Printf("Failed to send notification: %v", err)
		nh.writeJSON(w, http.StatusInternalServerError, response)
		return
	}

	log.Printf("Notification sent successfully to %s", req.Email)
	nh.writeJSON(w, http.StatusOK, response)
}

func (nh *NotificationHandler) validateRequest(req models.NotificationRequest) error {
	if strings.TrimSpace(req.Email) == "" {
		return &ValidationError{Field: "email", Message: "email is required"}
	}

	if !isValidEmail(req.Email) {
		return &ValidationError{Field: "email", Message: "invalid email format"}
	}

	if strings.TrimSpace(req.Site) == "" {
		return &ValidationError{Field: "site", Message: "site is required"}
	}

	if strings.TrimSpace(req.Time) == "" {
		return &ValidationError{Field: "time", Message: "time is required"}
	}

	return nil
}

func isValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	if len(email) == 0 {
		return false
	}

	atIndex := strings.LastIndex(email, "@")
	if atIndex <= 0 || atIndex == len(email)-1 {
		return false
	}

	localPart := email[:atIndex]
	domain := email[atIndex+1:]

	if len(localPart) == 0 || len(domain) == 0 {
		return false
	}

	if !strings.Contains(domain, ".") {
		return false
	}

	return true
}

func (nh *NotificationHandler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		nh.writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		})
		return
	}

	nh.writeJSON(w, http.StatusOK, map[string]string{
		"status":  "healthy",
		"service": "notification_service",
	})
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}