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

	http.HandleFunc("/authorize", wrappedHandlerAuth)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
