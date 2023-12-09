package app

import "github.com/gin-gonic/gin"

var (
	singeltonApp *gin.Engine
)

func Init() *gin.Engine {
	if singeltonApp == nil {
		singeltonApp = gin.Default()
	}

	return singeltonApp
}
