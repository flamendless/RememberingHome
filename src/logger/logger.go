package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func InitLog() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	newLogger, err := config.Build()
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		err := newLogger.Sync()
		if err != nil {
			fmt.Println(err)
		}
	}()
	logger = newLogger
}

func Log() *zap.Logger {
	return logger
}
