package logger

import (
	"log"
)

func Init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func Info(msg ...interface{}) {
	log.Println("[INFO]", msg)
}

func Error(msg ...interface{}) {
	log.Println("[ERROR]", msg)
}

func Fatal(msg ...interface{}) {
	log.Fatal("[FATAL]", msg)
}
