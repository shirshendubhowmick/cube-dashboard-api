package main

import (
	"main/internal/app"
	"main/internal/app/router"
)

func main() {
	appInstance := app.Init()

	router.Init(appInstance)
	appInstance.Run("0.0.0.0:9000")
}
