package initializers

import (
	"log"
	"os"
)

func InitLogger() {
	logFile := "D:/Logs/student-api.log"
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	log.SetOutput(file)
}
