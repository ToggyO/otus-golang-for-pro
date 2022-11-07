package logger

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.Logger
}

func NewLogger(level string, isDev bool, output io.Writer) *Logger {
	loggerBuilder := NewLoggerBuilder(level, isDev)
	return &Logger{
		logger: loggerBuilder.BuildAdvancedLogger(output),
	}
}

func (l *Logger) Debug(message string) {
	if l.checkLevel(zapcore.DebugLevel) {
		l.logger.Debug(message)
	}
}

func (l *Logger) Info(message string) {
	if l.checkLevel(zapcore.InfoLevel) {
		l.logger.Info(message)
	}
}

func (l *Logger) Warn(message string) {
	if l.checkLevel(zapcore.WarnLevel) {
		l.logger.Warn(message)
	}
}

func (l *Logger) Error(message string) {
	if l.checkLevel(zapcore.ErrorLevel) {
		l.logger.Error(message)
	}
}

func (l *Logger) Fatal(message string) {
	if l.checkLevel(zapcore.FatalLevel) {
		l.logger.Fatal(message)
	}
}

func (l *Logger) Flush() error {
	return l.logger.Sync()
}

func (l *Logger) checkLevel(level zapcore.Level) bool {
	return l.logger.Core().Enabled(level)
}
