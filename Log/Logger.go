package Log

import (
	"log"
	"os"
)

var Log *log.Logger

type Logger struct {
}

func NewLogger() *Logger {
	return &Logger{}
}

func (i *Logger) Init() {
	Log = log.New(os.Stdout, "", log.Ldate|log.Ltime)
}

func Sfatal(err error, location_error string) {
	if err != nil {
		log.Println("location_error : ", location_error)
		log.Fatal(err)
	}
}
