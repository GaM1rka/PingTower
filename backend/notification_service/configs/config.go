package configs

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	// MailerSend Configuration
	MailerSendAPIKey string
	FromEmail        string
	FromName         string
	MaxRetries       int = 3

	// SMTP Configuration
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string

	// Kafka Configuration
	KafkaBrokers           []string
	KafkaTopic             string
	KafkaConsumerGroup     string
	KafkaSessionTimeout    time.Duration
	KafkaHeartbeatInterval time.Duration

	// Health Check Server
	HealthPort string
)

func Configure() {
	// MailerSend Configuration
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

	// SMTP Configuration
	SMTPHost = os.Getenv("SMTP_HOST")
	if SMTPHost == "" {
		SMTPHost = "smtp.mailersend.net"
	}

	smtpPortStr := os.Getenv("SMTP_PORT")
	if smtpPortStr == "" {
		SMTPPort = 587
	} else {
		if port, err := strconv.Atoi(smtpPortStr); err == nil {
			SMTPPort = port
		} else {
			SMTPPort = 587
		}
	}

	SMTPUsername = os.Getenv("SMTP_USERNAME")
	if SMTPUsername == "" {
		log.Fatal("SMTP_USERNAME environment variable is required")
	}

	SMTPPassword = os.Getenv("SMTP_PASSWORD")
	if SMTPPassword == "" {
		log.Fatal("SMTP_PASSWORD environment variable is required")
	}

	// Kafka Configuration
	kafkaBrokersStr := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokersStr == "" {
		kafkaBrokersStr = "localhost:19092"
		log.Printf("KAFKA_BROKERS not set, using default: %s", kafkaBrokersStr)
	}
	KafkaBrokers = strings.Split(kafkaBrokersStr, ",")

	KafkaTopic = os.Getenv("KAFKA_TOPIC")
	if KafkaTopic == "" {
		KafkaTopic = "notification-alerts"
		log.Printf("KAFKA_TOPIC not set, using default: %s", KafkaTopic)
	}

	KafkaConsumerGroup = os.Getenv("KAFKA_CONSUMER_GROUP")
	if KafkaConsumerGroup == "" {
		KafkaConsumerGroup = "notification-service"
		log.Printf("KAFKA_CONSUMER_GROUP not set, using default: %s", KafkaConsumerGroup)
	}

	KafkaSessionTimeout = 10 * time.Second
	KafkaHeartbeatInterval = 3 * time.Second

	// Health Check Server
	HealthPort = os.Getenv("HEALTH_PORT")
	if HealthPort == "" {
		HealthPort = "8084"
	}

	log.Printf("Notification Service configured:")
	log.Printf("- Health Port: %s", HealthPort)
	log.Printf("- From Email: %s", FromEmail)
	log.Printf("- From Name: %s", FromName)
	log.Printf("- Max Retries: %d", MaxRetries)
	log.Printf("- SMTP Host: %s:%d", SMTPHost, SMTPPort)
	log.Printf("- SMTP Username: %s", SMTPUsername)
	log.Printf("- Kafka Brokers: %v", KafkaBrokers)
	log.Printf("- Kafka Topic: %s", KafkaTopic)
	log.Printf("- Kafka Consumer Group: %s", KafkaConsumerGroup)
}
