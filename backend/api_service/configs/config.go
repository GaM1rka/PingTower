package configs

import (
	"log"
	"net/http"
	"os"
)

const (
	JWTURL  = "https://600a28bf-3e6d-4221-a9bc-f769f0de3536.mock.pstmn.io/authorize"
	DBURL   = ""
	PingURL = ""
)

var APILogger *log.Logger
var Client *http.Client

func Configure() {
	APILogger = log.New(os.Stdout, "LOGGER: ", log.LstdFlags)
	Client = &http.Client{}
}
