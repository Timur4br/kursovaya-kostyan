package configs

import (
	"log"
	"os"
)

func InitLogging() {
	logsOutputFile, err := os.OpenFile(os.Getenv("LOGS_FILE_PATH"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logsOutputFile)
}
