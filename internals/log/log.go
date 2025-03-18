package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func NewLogger(level string) *zap.SugaredLogger {
	config := zap.NewProductionConfig()
	config.Level.SetLevel(getLevel(level))
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000000")
	config.Encoding = "console"
	config.DisableStacktrace = true
	logger, _ := config.Build()
	return logger.Sugar()
}

func getLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func Fatalf(template string, args ...interface{}) {
	fmt.Printf(template, args)
	fmt.Println()
	os.Exit(1)
}
