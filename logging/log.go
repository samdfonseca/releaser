package logging

import (
	"fmt"
	"log"
	"os"
)

var (
	isDebug     = (os.Getenv("DEBUG") != "0" && os.Getenv("DEBUG") != "")
	loggerFlags = log.Ldate | log.Ltime | log.Lshortfile
)

type Logger struct {
	*log.Logger
}

func New(command string) *Logger {
	return &Logger{
		log.New(os.Stdout, fmt.Sprintf("releaser - %s - ", command), loggerFlags),
	}
}

func (l *Logger) Debug(v ...interface{}) {
	if isDebug {
		l.Print(v...)
	}
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if isDebug {
		l.Printf(format, v...)
	}
}

func (l *Logger) Debugln(msg string) {
	if isDebug {
		l.Println(msg)
	}
}
