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

type Builder struct {
	minLevel zapcore.Level
	isDev    bool
}

func NewLoggerBuilder(level string, isDev bool) *Builder {
	var zapLevel zapcore.Level
	switch strings.ToLower(level) {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "fatal":
		zapLevel = zapcore.FatalLevel
	default:
		zapLevel = zapcore.InfoLevel
	}
	return &Builder{
		minLevel: zapLevel,
		isDev:    isDev,
	}
}

func (l *Builder) BuildAdvancedLogger(output io.Writer) *zap.Logger {
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

func (l *Builder) getEncoder() zapcore.Encoder {
	encoderCfg := zap.NewDevelopmentEncoderConfig()
	if !l.isDev {
		encoderCfg = zap.NewProductionEncoderConfig()
	}

	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder

	return zapcore.NewConsoleEncoder(encoderCfg)
}

func (l *Builder) getConsoleWriter() zapcore.WriteSyncer {
	return zapcore.AddSync(os.Stdout)
}

func (l *Builder) getWriter(writer io.Writer) zapcore.WriteSyncer {
	return zapcore.AddSync(writer)
}

func (l *Builder) getFileWriter(fileName string) zapcore.WriteSyncer {
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
