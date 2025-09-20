package internal

import (
	"context"
	"fmt"
	"log"
	"notification_service/configs"
	"notification_service/models"
	"time"

	"github.com/mailersend/mailersend-go"
)

type EmailService struct {
	client *mailersend.Mailersend
}

func NewEmailService() *EmailService {
	ms := mailersend.NewMailersend(configs.MailerSendAPIKey)
	return &EmailService{
		client: ms,
	}
}

func (es *EmailService) SendNotification(ctx context.Context, req models.NotificationRequest) (*models.NotificationResponse, error) {
	content := GenerateEmailContent(req)
	
	message := es.client.Email.NewMessage()
	
	from := mailersend.From{
		Name:  configs.FromName,
		Email: configs.FromEmail,
	}
	
	recipients := []mailersend.Recipient{
		{
			Name:  "",
			Email: req.Email,
		},
	}
	
	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject(content.Subject)
	message.SetHTML(content.HTML)
	message.SetText(content.Text)
	
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	
	res, err := es.client.Email.Send(ctxWithTimeout, message)
	if err != nil {
		log.Printf("Failed to send email to %s: %v", req.Email, err)
		return &models.NotificationResponse{
			Status: "failed",
			Error:  err.Error(),
		}, err
	}
	
	var messageID string
	if res != nil {
		messageID = res.Header.Get("X-Message-Id")
	}
	
	log.Printf("Email sent successfully to %s, Message ID: %s", req.Email, messageID)
	
	return &models.NotificationResponse{
		Status:    "sent",
		MessageID: messageID,
	}, nil
}

func (es *EmailService) SendNotificationWithRetry(ctx context.Context, req models.NotificationRequest, maxRetries int) (*models.NotificationResponse, error) {
	var lastErr error
	
	for attempt := 1; attempt <= maxRetries; attempt++ {
		response, err := es.SendNotification(ctx, req)
		if err == nil {
			return response, nil
		}
		
		lastErr = err
		log.Printf("Attempt %d/%d failed for email %s: %v", attempt, maxRetries, req.Email, err)
		
		if attempt < maxRetries {
			backoffDuration := time.Duration(attempt) * 2 * time.Second
			log.Printf("Retrying in %v...", backoffDuration)
			
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoffDuration):
			}
		}
	}
	
	return &models.NotificationResponse{
		Status: "failed",
		Error:  fmt.Sprintf("Failed after %d attempts: %v", maxRetries, lastErr),
	}, lastErr
}