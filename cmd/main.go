package main

import (
	"main/configs"
	"main/internal/app"
	"main/internal/app/router"
	"main/internal/db"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	configs.InitEnvConfig()

	_, err = db.Init()

	if err != nil {
		panic(err)
	}

	appInstance := app.Init()

	router.Init(appInstance)
	appInstance.Run("0.0.0.0:9000")
}
