package models

import "hash/fnv"

type NotificationRequest struct {
	Email string `json:"email"`
	Site  string `json:"site"`
	Time  string `json:"time"`
}

func (r NotificationRequest) GetHashCode() uint32 {
	h := fnv.New32a()
	h.Write([]byte(r.Email + r.Site + r.Time))
	return h.Sum32()
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