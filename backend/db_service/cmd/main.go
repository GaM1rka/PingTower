package main

import (
	"db_service/configs"
	"db_service/internal"
	"log"
	"net/http"
)

func main() {
	configs.Configure()
	// Инициализация баз данных
	psqlDB, err := internal.InitPostgreSQL()
	if err != nil {
		log.Fatal("PostgreSQL init error:", err)
	}
	defer psqlDB.Close()

	chDB, err := internal.InitClickHouse()
	if err != nil {
		log.Fatal("ClickHouse init error:", err)
	}
	defer chDB.Close()

	store := internal.NewStorage(psqlDB, chDB)

	handler := internal.NewHandler(store)

	http.HandleFunc("/user", handler.UserHandler)
	http.HandleFunc("/user/sites/", handler.UserSitesHandler)
	http.HandleFunc("/checker/", handler.CheckerHandler)
	http.HandleFunc("/checkers", handler.CheckersHandler)
	http.HandleFunc("/all-users-sites", handler.AllUsersSitesHandler) // GET
	http.HandleFunc("/ping", handler.PingHandler)                     // POST
	http.HandleFunc("/user/", handler.UserEmailHandler)

	configs.DBLogger.Println("Server starting on :8083")
	err = http.ListenAndServe(":8083", nil)
	if err != nil {
		configs.DBLogger.Println("Error while listening DB endpoints: ", err.Error())
	}
}
