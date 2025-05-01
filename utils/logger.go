package utils

import (
	"log"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var ZapLogger *zap.Logger
var SugaredLogger *zap.SugaredLogger

func LoadLogger() {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Development = true
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	var err error
	ZapLogger, err = loggerConfig.Build()
	if err != nil {
		log.Fatal(err)
	}
	SugaredLogger = ZapLogger.Sugar()
}
