package util

import (
	"log"
	"os"
	"sync"
)

// Logger struct.
type Logger struct {
	filename string
	*log.Logger
}

var logger *Logger
var once sync.Once

// GetLogger creates singleton logger for envssh.
func GetLogger() *Logger {
	once.Do(func() {
		logger = createLogger("envssh.log")
	})
	return logger
}

func createLogger(fname string) *Logger {
	file, _ := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0700)

	return &Logger{
		filename: fname,
		Logger:   log.New(file, "envssh ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
