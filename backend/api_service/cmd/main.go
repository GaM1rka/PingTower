package main

import (
	"api_service/configs"
	"api_service/internal"
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
	wrappedHandlerAuth := enableCORS(loggingMiddleware(handler.AuthHandler))

	http.HandleFunc("/register", wrappedHandlerAuth)
	http.HandleFunc("/login")
	http.HandleFunc("/checker/") // проверять id пользователя по jwt. Возвращает инфу по id пользователя и id сайта инфу по нему
	http.HandleFunc("/checkers") // GET Возвращает по id пользователя все url и их последний статус(работает или нет).
	http.HandleFunc("/checkers") // POST добавляет для пользователя сайт для мониторинга.

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
	go configs.StartCronScheduler("http://localhost:8080") // Потом поменять для виртуальной машины.
}
