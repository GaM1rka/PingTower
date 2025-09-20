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

	handler := internal.NewNotificationHandler()

	mux := http.NewServeMux()
	mux.HandleFunc("/send-notification", handler.SendNotificationHandler)
	mux.HandleFunc("/health", handler.HealthHandler)

	middlewareStack := internal.LoggingMiddleware(internal.CORSMiddleware(mux))

	server := &http.Server{
		Addr:         ":" + configs.ServerPort,
		Handler:      middlewareStack,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("Notification Service starting on port %s", configs.ServerPort)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("Server failed to start: %v", err)

	case sig := <-shutdown:
		log.Printf("Received signal %v, starting graceful shutdown", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)
			if err := server.Close(); err != nil {
				log.Printf("Force shutdown failed: %v", err)
			}
		}
		log.Println("Notification Service stopped")
	}
}