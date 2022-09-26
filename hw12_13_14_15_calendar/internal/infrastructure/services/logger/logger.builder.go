package logger

import (
	"fmt"
	"io"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

type loggerBuilder struct {
	minLevel zapcore.Level
	isDev    bool
}

func NewLoggerBuilder(level string, isDev bool) *loggerBuilder {
	var zapLevel zapcore.Level
	switch strings.ToLower(level) {
	case "debug":
		zapLevel = zapcore.DebugLevel
		break
	case "error":
		zapLevel = zapcore.ErrorLevel
		break
	case "warn":
		zapLevel = zapcore.WarnLevel
		break
	case "fatal":
		zapLevel = zapcore.FatalLevel
		break
	default:
		zapLevel = zapcore.InfoLevel
		break
	}
	return &loggerBuilder{
		minLevel: zapLevel,
		isDev:    isDev,
	}
}

func (l *loggerBuilder) BuildAdvancedLogger(output io.Writer) *zap.Logger {
	encoder := l.getEncoder()
	writeSyncer := l.getConsoleWriter()
	if output != nil {
		writeSyncer = l.getWriter(output)
	}

	var core zapcore.Core
	consoleCore := zapcore.NewCore(encoder, writeSyncer, l.minLevel)

	if !l.isDev {
		writeSyncerFile := l.getFileWriter("logs")
		fileCore := zapcore.NewCore(encoder, writeSyncerFile, l.minLevel)
		core = zapcore.NewTee(consoleCore, fileCore)
	} else {
		core = consoleCore
	}

	return zap.New(core)
}

func (l *loggerBuilder) getEncoder() zapcore.Encoder {
	encoderCfg := zap.NewDevelopmentEncoderConfig()
	if !l.isDev {
		encoderCfg = zap.NewProductionEncoderConfig()
	}

	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder

	return zapcore.NewConsoleEncoder(encoderCfg)
}

func (l *loggerBuilder) getConsoleWriter() zapcore.WriteSyncer {
	return zapcore.AddSync(os.Stdout)
}

func (l *loggerBuilder) getWriter(writer io.Writer) zapcore.WriteSyncer {
	return zapcore.AddSync(writer)
}

func (l *loggerBuilder) getFileWriter(fileName string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		// TODO: вынести в конфиг?
		Filename:   fmt.Sprintf("./logs/%s.log", fileName),
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}
