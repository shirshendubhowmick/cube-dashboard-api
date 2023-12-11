package main

import (
	"main/configs"
	"main/internal/app"
	"main/internal/app/router"
	"main/internal/db"
	"main/internal/services/logger"

	"github.com/joho/godotenv"
)

func main() {
	log := logger.Init(true)
	logger.Log.Info("Logger initialized")
	defer log.Sync()

	err := godotenv.Load()
	if err != nil {
		logger.Log.Panicw("Error loading .env file", "error", err)
	}
	configs.InitEnvConfig()

	_, err = db.Init()

	if err != nil {
		logger.Log.Panicw("Error initializing database", "error", err)
	}

	appInstance := app.Init()

	router.Init(appInstance)
	appInstance.Run("0.0.0.0:9000")
}
