package configs

import (
	"log"
	"net/http"
	"os"
)

var PingLogger *log.Logger
var Client *http.Client

func Configure() {
	PingLogger = log.New(os.Stdout, "LOGGER: ", log.LstdFlags)
	Client = &http.Client{}
}
