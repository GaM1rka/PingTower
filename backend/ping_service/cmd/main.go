package cmd

import (
	"net/http"
	"ping_service/configs"
	"ping_service/internal"
)

func main() {
	configs.Configure()

	http.HandleFunc("/ping", internal.PingHandler)

	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		return
	}
}
