package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLogger : type alias for zap.Logger
type ZapLogger = zap.SugaredLogger

// Init : Initializes the logger
func Init(isDevelopment bool, serviceName string) *ZapLogger {
	var loggerConfig zap.Config

	if isDevelopment {
		loggerConfig = zap.NewDevelopmentConfig()
	} else {
		loggerConfig = zap.NewProductionConfig()
		loggerConfig.OutputPaths = []string{"stdout"}
		loggerConfig.ErrorOutputPaths = []string{"stderr"}
		loggerConfig.InitialFields = map[string]interface{}{"name": serviceName}
		loggerConfig.EncoderConfig.TimeKey = "time"
		loggerConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	log, err := loggerConfig.Build()

	if err != nil {
		panic(err)
	}

	return log.Sugar()
}
