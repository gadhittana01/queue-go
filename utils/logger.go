package utils

import (
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	var err error
	config := zap.NewProductionConfig()
	logger, err = config.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
}

func LogInfo(message string, fields ...zap.Field) {
	logger.Info(message, fields...)
}

func LogFatal(message string, fields ...zap.Field) {
	logger.Fatal(message, fields...)
}

func LogPanic(message string, fields ...zap.Field) {
	logger.Panic(message, fields...)
}

func LogDebug(message string, fields ...zap.Field) {
	logger.Debug(message, fields...)
}

func LogError(message string, fields ...zap.Field) {
	logger.Error(message, fields...)
}

func LogWarning(message string, fields ...zap.Field) {
	logger.Warn(message, fields...)
}
