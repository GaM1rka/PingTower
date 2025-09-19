package main

import (
	"auth/configs"
	"auth/internal"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	configs.InitLogger()
	configs.InitConfig()
	configs.Logger.Println("Starting JWT auth service...")

	if err := configs.InitRedis(); err != nil {
		configs.Logger.Fatal("Failed to connect to Redis:", err)
	}

	http.HandleFunc("/generate", internal.CORSMiddleware(internal.LoggingMiddleware(internal.RateLimitMiddleware(internal.GenerateTokenHandler))))
	http.HandleFunc("/validate", internal.CORSMiddleware(internal.LoggingMiddleware(internal.ValidateHandler)))
	http.HandleFunc("/refresh", internal.CORSMiddleware(internal.LoggingMiddleware(internal.RateLimitMiddleware(internal.RefreshHandler))))
	http.HandleFunc("/health", internal.HealthHandler)

	server := &http.Server{
		Addr: ":" + configs.ServerPort,
	}

	go func() {
		configs.Logger.Println("Auth service starting on :" + configs.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			configs.Logger.Fatal("Server failed to start:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	configs.Logger.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		configs.Logger.Fatal("Server forced to shutdown:", err)
	}

	configs.Logger.Println("Server exited")
}