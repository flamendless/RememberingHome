package logger

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var loggerSync sync.Once

func Log() *zap.Logger {
	loggerSync.Do(func() {
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
	})
	return logger
}
