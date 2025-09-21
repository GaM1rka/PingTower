package main

import (
	"api_service/configs"
	"api_service/internal"
	"encoding/json"
	"net/http"
)

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		configs.APILogger.Printf("Received %s request for %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		configs.APILogger.Printf("Headers: %v", r.Header)

		next(w, r)
	}
}

func main() {
	configs.Configure()

	handler := internal.NewHandler()
	swaggerHandler := internal.NewSwaggerHandler()

	// Обертки для middleware
	wrappedRegister := enableCORS(loggingMiddleware(handler.AuthHandler))
	wrappedLogin := enableCORS(loggingMiddleware(handler.AuthHandler))
	wrappedChecker := enableCORS(loggingMiddleware(handler.CheckerHandler))
	wrappedCheckers := enableCORS(loggingMiddleware(handler.CheckersHandler))
	wrappedPingAll := enableCORS(loggingMiddleware(handler.PingAllHandler))
	
	// Обертки для Swagger
	wrappedSwaggerSpec := enableCORS(loggingMiddleware(swaggerHandler.ServeSwaggerSpec))
	wrappedSwaggerUI := enableCORS(loggingMiddleware(swaggerHandler.ServeSwaggerUI))

	// Регистрируем все обработчики
	http.HandleFunc("/register", wrappedRegister)
	http.HandleFunc("/login", wrappedLogin)
	http.HandleFunc("/checker/", wrappedChecker)
	http.HandleFunc("/checkers", wrappedCheckers)
	http.HandleFunc("/pingAll", wrappedPingAll)

	// Добавляем health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Serve Swagger endpoints
	http.HandleFunc("/swagger/spec", wrappedSwaggerSpec)
	http.HandleFunc("/swagger/", wrappedSwaggerUI)
	http.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
		// Redirect /swagger to /swagger/
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})

	configs.APILogger.Println("API Service starting on :8080")

	go configs.StartCronScheduler("http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		configs.APILogger.Panic("Error starting server:", err)
	}
}