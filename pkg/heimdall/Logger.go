package heimdall

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type Logger struct {
	prefix string
	flags  int
}

var (
	loggerInstance *Logger
	loggerOnce     sync.Once
)

func init() {
	loggerOnce.Do(func() {
		loggerInstance = NewLogger("MyLogger")
	})
}

func Debug(msg string) {
	if !(DebugMode) {
		return
	}
	loggerInstance.log(msg, "DEBUG")
}

func Info(msg string) {
	loggerInstance.log(msg, "INFO")
}

func Warn(msg string) {
	loggerInstance.log(msg, "WARN")
}

func Error(msg string) {
	loggerInstance.log(msg, "ERROR")
}

func Fatal(msg string) {
	loggerInstance.log(msg, "FATAL")
	os.Exit(1)
}

func NewLogger(prefix string) *Logger {
	return &Logger{
		prefix: prefix,
		flags:  0,
	}
}

func (l *Logger) log(msg, level string) {
	formattedMessage := fmt.Sprintf("[%s] [%s] %s", time.Now().Format("2006-01-02T15:04:05:123"), level, msg)
	log.SetFlags(l.flags)
	log.Printf(formattedMessage)
}
