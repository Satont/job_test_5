package logger_impl

import (
	"fmt"
	"github.com/satont/test/internal/services/logger"
	"log"
	"strings"
)

type consoleLogger struct{}

func NewConsoleLogger() logger.Logger {
	return &consoleLogger{}
}

func (l *consoleLogger) Info(msg string, args ...string) {
	fmt.Printf("%s %s\n", msg, strings.Join(args, " "))
}

func (l *consoleLogger) Error(args ...any) {
	log.Println(args)
}
