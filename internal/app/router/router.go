package router

import (
	"main/internal/app/controllers"
	"main/internal/app/middlewares"

	"github.com/gin-gonic/gin"
)

func Init(app *gin.Engine) {
	app.Use(middlewares.Cors())
	rootRoute := app.Group("/")
	rootRoute.GET("/", controllers.RootController)

	dataRoute := app.Group("/data")
	dataRoute.PUT("/ingest/meteriote", controllers.HandleMeteoriteData)

	userRoute := app.Group("/user")
	userRoute.POST("/session", controllers.CreateUserSession)
	userRoute.GET("/session", middlewares.Authorization(), controllers.GetUserSession)
	userRoute.DELETE("/session", controllers.DeleteUserSession)
}
