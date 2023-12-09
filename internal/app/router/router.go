package router

import (
	"main/internal/app/controllers"

	"github.com/gin-gonic/gin"
)

func Init(app *gin.Engine) {
	rootRoute := app.Group("/")

	rootRoute.GET("/", controllers.RootController)
}
