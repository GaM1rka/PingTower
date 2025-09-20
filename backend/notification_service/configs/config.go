package configs

import (
	"log"
	"os"
)

var (
	MailerSendAPIKey string
	FromEmail        string
	FromName         string
	ServerPort       string
	MaxRetries       int = 3
)

func Configure() {
	MailerSendAPIKey = os.Getenv("MAILERSEND_API_KEY")
	if MailerSendAPIKey == "" {
		log.Fatal("MAILERSEND_API_KEY environment variable is required")
	}

	FromEmail = os.Getenv("FROM_EMAIL")
	if FromEmail == "" {
		FromEmail = "noreply@yourdomain.com"
		log.Printf("FROM_EMAIL not set, using default: %s", FromEmail)
	}

	FromName = os.Getenv("FROM_NAME")
	if FromName == "" {
		FromName = "PingTower Alert System"
	}

	ServerPort = os.Getenv("SERVER_PORT")
	if ServerPort == "" {
		ServerPort = "8084"
	}

	log.Printf("Notification Service configured:")
	log.Printf("- Server Port: %s", ServerPort)
	log.Printf("- From Email: %s", FromEmail)
	log.Printf("- From Name: %s", FromName)
	log.Printf("- Max Retries: %d", MaxRetries)
}