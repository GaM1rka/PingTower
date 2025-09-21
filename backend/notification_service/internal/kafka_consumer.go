package internal

import (
	"context"
	"encoding/json"
	"log"
	"notification_service/configs"
	"notification_service/models"
	"sync"

	"github.com/IBM/sarama"
)

type KafkaConsumer struct {
	consumerGroup sarama.ConsumerGroup
	smtpService   *SMTPService
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
}

type ConsumerGroupHandler struct {
	smtpService *SMTPService
}

func NewKafkaConsumer() (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Group.Session.Timeout = configs.KafkaSessionTimeout
	config.Consumer.Group.Heartbeat.Interval = configs.KafkaHeartbeatInterval

	consumerGroup, err := sarama.NewConsumerGroup(configs.KafkaBrokers, configs.KafkaConsumerGroup, config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &KafkaConsumer{
		consumerGroup: consumerGroup,
		smtpService:   NewSMTPService(),
		ctx:           ctx,
		cancel:        cancel,
	}, nil
}

func (kc *KafkaConsumer) Start() error {
	handler := &ConsumerGroupHandler{
		smtpService: kc.smtpService,
	}

	kc.wg.Add(1)
	go func() {
		defer kc.wg.Done()
		for {
			select {
			case <-kc.ctx.Done():
				log.Println("Kafka consumer context cancelled")
				return
			default:
				if err := kc.consumerGroup.Consume(kc.ctx, []string{configs.KafkaTopic}, handler); err != nil {
					log.Printf("Error consuming from Kafka: %v", err)
					return
				}
			}
		}
	}()

	log.Printf("Kafka consumer started, listening to topic: %s", configs.KafkaTopic)
	return nil
}

func (kc *KafkaConsumer) Stop() {
	log.Println("Stopping Kafka consumer...")
	kc.cancel()
	kc.wg.Wait()
	
	if err := kc.consumerGroup.Close(); err != nil {
		log.Printf("Error closing consumer group: %v", err)
	}
	log.Println("Kafka consumer stopped")
}

func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Println("Message channel closed")
				return nil
			}

			if err := h.processMessage(message); err != nil {
				log.Printf("Error processing message: %v", err)
			} else {
				session.MarkMessage(message, "")
			}

		case <-session.Context().Done():
			return nil
		}
	}
}

func (h *ConsumerGroupHandler) processMessage(message *sarama.ConsumerMessage) error {
	log.Printf("Received message from topic %s, partition %d, offset %d", 
		message.Topic, message.Partition, message.Offset)

	var notificationReq models.NotificationRequest
	if err := json.Unmarshal(message.Value, &notificationReq); err != nil {
		log.Printf("Failed to unmarshal message: %v, raw message: %s", err, string(message.Value))
		return err
	}

	if err := validateKafkaMessage(notificationReq); err != nil {
		log.Printf("Invalid message format: %v", err)
		return err
	}

	log.Printf("Processing notification for email: %s, site: %s, time: %s", 
		notificationReq.Email, notificationReq.Site, notificationReq.Time)

	ctx := context.Background()
	response, err := h.smtpService.SendNotificationWithRetry(ctx, notificationReq, configs.MaxRetries)
	if err != nil {
		log.Printf("Failed to send notification: %v", err)
		return err
	}

	log.Printf("Notification sent successfully: status=%s, message_id=%s", 
		response.Status, response.MessageID)
	return nil
}

func validateKafkaMessage(req models.NotificationRequest) error {
	if req.Email == "" {
		return &ValidationError{Field: "email", Message: "email is required"}
	}
	if req.Site == "" {
		return &ValidationError{Field: "site", Message: "site is required"}
	}
	if req.Time == "" {
		return &ValidationError{Field: "time", Message: "time is required"}
	}
	return nil
}