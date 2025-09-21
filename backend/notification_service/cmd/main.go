package main

import (
	"context"
	"log"
	"net/http"
	"notification_service/configs"
	"notification_service/internal"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	configs.Configure()

	// Initialize Kafka consumer
	kafkaConsumer, err := internal.NewKafkaConsumer()
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}

	// Start Kafka consumer
	if err := kafkaConsumer.Start(); err != nil {
		log.Fatalf("Failed to start Kafka consumer: %v", err)
	}

	// Start minimal health check HTTP server
	healthHandler := internal.NewHealthHandler()
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler.HealthHandler)

	middlewareStack := internal.LoggingMiddleware(mux)

	server := &http.Server{
		Addr:         ":" + configs.HealthPort,
		Handler:      middlewareStack,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("Health check server starting on port %s", configs.HealthPort)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Printf("Health server error: %v", err)
		kafkaConsumer.Stop()

	case sig := <-shutdown:
		log.Printf("Received signal %v, starting graceful shutdown", sig)

		// Stop Kafka consumer first
		kafkaConsumer.Stop()

		// Then stop health server
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Health server shutdown failed: %v", err)
			if err := server.Close(); err != nil {
				log.Printf("Force shutdown failed: %v", err)
			}
		}
		log.Println("Notification Service stopped")
	}
}