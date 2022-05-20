package config

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

var SugarLogger *zap.SugaredLogger

func InitializeLogger() {
	config := zap.NewDevelopmentEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	writer := zapcore.AddSync(os.Stdout)
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, getLogLevel()),
	)
	Logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	defer Logger.Sync()
	SugarLogger = Logger.Sugar()
}

func getLogLevel() zapcore.Level {
	configuration := GetConfig()
	fmt.Println("Level", configuration.Log.Level)
	switch strings.ToLower(configuration.Log.Level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	default:
		return zapcore.InfoLevel
	}
}
