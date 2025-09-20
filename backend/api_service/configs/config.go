package configs

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-co-op/gocron"
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

func StartCronScheduler(apiBaseURL string) {
	s := gocron.NewScheduler(time.UTC)

	s.Cron("*/5 * * * *").Do(func() {
		APILogger.Println("New cron request")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiBaseURL+"/pingAll", nil)
		if err != nil {
			APILogger.Println("Error creating request:", err)
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			APILogger.Println("Error calling /pingAll:", err)
			return
		}
		defer resp.Body.Close()

		APILogger.Println("/pingAll called, response code:", resp.StatusCode)
	})

	s.StartAsync()
}
