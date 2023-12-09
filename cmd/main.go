package main

import (
	"main/configs"
	"main/internal/app"
	"main/internal/app/router"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	configs.InitEnvConfig()

	appInstance := app.Init()

	router.Init(appInstance)
	appInstance.Run("0.0.0.0:9000")
}
