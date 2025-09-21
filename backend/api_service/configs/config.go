package configs

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/segmentio/kafka-go"
)

const (
	JWTURL      = "http://auth_service:8081/generate"
	JWTValidate = "http://auth_service:8081/validate"
	DBURL       = "http://db_service:8083"
	PingURL     = "http://ping_service:8082"
	KafkaAddr   = "kafka1:29092"
)

var APILogger *log.Logger
var Client *http.Client
var KafkaWriter *kafka.Writer

func Configure() {
	APILogger = log.New(os.Stdout, "API_SERVICE: ", log.LstdFlags)
	Client = &http.Client{
		Timeout: 30 * time.Second,
	}

	// Инициализация Kafka writer
	KafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP(KafkaAddr),
		Topic:    "notification-alerts",
		Balancer: &kafka.LeastBytes{},
	}
}

func StartCronScheduler(apiBaseURL string) {
	s := gocron.NewScheduler(time.UTC)

	s.Every(1).Minute().Do(func() {
		APILogger.Println("Starting cron job: pingAll")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiBaseURL+"/pingAll", nil)
		if err != nil {
			APILogger.Println("Error creating pingAll request:", err)
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			APILogger.Println("Error calling /pingAll:", err)
			return
		}
		defer resp.Body.Close()

		APILogger.Println("/pingAll called, response code:", resp.StatusCode)
	})

	s.StartAsync()
}

func SendKafkaNotification(email, site string, responseTime int64) error {
	message := map[string]interface{}{
		"email":  email,
		"site":   site,
		"time":   responseTime,
		"status": "down",
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return KafkaWriter.WriteMessages(context.Background(),
		kafka.Message{
			Value: jsonData,
		},
	)
}
