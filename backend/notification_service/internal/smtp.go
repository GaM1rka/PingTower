package internal

import (
	"context"
	"fmt"
	"log"
	"notification_service/configs"
	"notification_service/models"
	"time"

	"gopkg.in/gomail.v2"
)

type SMTPService struct {
	dialer *gomail.Dialer
}

func NewSMTPService() *SMTPService {
	d := gomail.NewDialer(
		configs.SMTPHost,
		configs.SMTPPort,
		configs.SMTPUsername,
		configs.SMTPPassword,
	)
	
	return &SMTPService{
		dialer: d,
	}
}

func (s *SMTPService) SendNotification(ctx context.Context, req models.NotificationRequest) (*models.NotificationResponse, error) {
	content := GenerateEmailContent(req)
	
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", configs.FromName, configs.SMTPUsername))
	m.SetHeader("To", req.Email)
	m.SetHeader("Subject", content.Subject)
	m.SetBody("text/plain", content.Text)
	m.AddAlternative("text/html", content.HTML)
	
	if err := s.dialer.DialAndSend(m); err != nil {
		log.Printf("Failed to send email to %s: %v", req.Email, err)
		return &models.NotificationResponse{
			Status: "failed",
			Error:  err.Error(),
		}, err
	}
	
	log.Printf("Email sent successfully to %s via SMTP", req.Email)
	
	return &models.NotificationResponse{
		Status:    "sent",
		MessageID: fmt.Sprintf("smtp-%d", req.GetHashCode()),
	}, nil
}

func (s *SMTPService) SendNotificationWithRetry(ctx context.Context, req models.NotificationRequest, maxRetries int) (*models.NotificationResponse, error) {
	var lastErr error
	
	for attempt := 1; attempt <= maxRetries; attempt++ {
		response, err := s.SendNotification(ctx, req)
		if err == nil {
			return response, nil
		}
		
		lastErr = err
		log.Printf("SMTP attempt %d/%d failed for email %s: %v", attempt, maxRetries, req.Email, err)
		
		if attempt < maxRetries {
			backoffDuration := time.Duration(attempt) * 2 * time.Second
			log.Printf("Retrying SMTP in %v...", backoffDuration)
			
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoffDuration):
			}
		}
	}
	
	return &models.NotificationResponse{
		Status: "failed",
		Error:  fmt.Sprintf("SMTP failed after %d attempts: %v", maxRetries, lastErr),
	}, lastErr
}