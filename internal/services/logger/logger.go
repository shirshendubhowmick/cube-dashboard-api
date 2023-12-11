package logger

import (
	"main/internal/services/logger/constants"
	logger "main/pkg"
)

const (
	contextKey = "type"
)

// Log : The main logger
var Log *logger.ZapLogger

var AppRequestResponseLog *logger.ZapLogger

var AppControllerLog *logger.ZapLogger

var AppMiddlewareLog *logger.ZapLogger

var AppServiceLog *logger.ZapLogger

var AppMiscLog *logger.ZapLogger

// Init : Initialized the logger
func Init(isDevelopment bool) *logger.ZapLogger {
	if Log == nil {
		Log = logger.Init(isDevelopment, "cube-dashboard-api")
		AppRequestResponseLog = Log.With(contextKey, constants.AppRequestResponseLoggerType)
		AppControllerLog = Log.With(contextKey, constants.AppControllerLoggerType)
		AppMiddlewareLog = Log.With(contextKey, constants.AppMiddlewareLoggerType)
		AppServiceLog = Log.With(contextKey, constants.AppServiceLoggerType)
		AppMiscLog = Log.With(contextKey, constants.AppMiscLoggerType)
	}

	return Log
}
