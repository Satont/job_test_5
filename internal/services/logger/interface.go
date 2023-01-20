package logger

type Logger interface {
	Info(msg string, args ...string)
	Error(args ...any)
}
