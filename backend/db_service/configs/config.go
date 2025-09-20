package configs

import (
	"log"
	"os"
)

// не nil по умолчанию
var DBLogger = log.New(os.Stdout, "[db_service] ", log.LstdFlags|log.Lshortfile)

// Configure можно оставить, если хочешь менять формат/вывод в рантайме:
func Configure() {
	DBLogger.SetPrefix("[db_service] ")
	DBLogger.SetFlags(log.LstdFlags | log.Lshortfile)
	DBLogger.SetOutput(os.Stdout)
}
