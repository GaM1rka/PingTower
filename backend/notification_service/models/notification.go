package models

type NotificationRequest struct {
	Email string `json:"email"`
	Site  string `json:"site"`
	Time  string `json:"time"`
}

type NotificationResponse struct {
	Status    string `json:"status"`
	MessageID string `json:"message_id,omitempty"`
	Error     string `json:"error,omitempty"`
}

type EmailContent struct {
	Subject string
	HTML    string
	Text    string
}