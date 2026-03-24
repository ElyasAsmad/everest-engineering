package logger

import (
	"os"
	"strings"

	"github.com/charmbracelet/log"
)

type AppLogger struct {
	logger *log.Logger
}

// check log level from env, if not set, default to info level
func parseLogLevel() log.Level {
	logLevel := os.Getenv("EE_LOG_LEVEL")
	logLevel = strings.ToLower(logLevel)

	switch logLevel {
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	case "warn":
		return log.WarnLevel
	case "error":
		return log.ErrorLevel
	case "fatal":
		return log.FatalLevel
	default:
		return log.InfoLevel
	}
}

func NewLogger() *AppLogger {
	logger := log.New(os.Stderr)
	logger.SetLevel(parseLogLevel())
	return &AppLogger{logger: logger}
}

func (l *AppLogger) Debug(msg any, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *AppLogger) Debugf(format string, args ...any) {
	l.logger.Debugf(format, args...)
}

func (l *AppLogger) Info(msg any, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *AppLogger) Warn(msg any, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *AppLogger) Error(msg any, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *AppLogger) Print(msg any, args ...any) {
	l.logger.Print(msg, args...)
}

func (l *AppLogger) Infof(format string, args ...any) {
	l.logger.Infof(format, args...)
}

func (l *AppLogger) Warnf(format string, args ...any) {
	l.logger.Warnf(format, args...)
}

func (l *AppLogger) Errorf(format string, args ...any) {
	l.logger.Errorf(format, args...)
}

func (l *AppLogger) Printf(format string, args ...any) {
	l.logger.Printf(format, args...)
}

func (l *AppLogger) Fatal(msg any, args ...any) {
	l.logger.Fatal(msg, args...)
}

func (l *AppLogger) Fatalf(format string, args ...any) {
	l.logger.Fatalf(format, args...)
}
