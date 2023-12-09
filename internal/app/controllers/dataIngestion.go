package controllers

import (
	"main/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleMeteoriteData(ginContext *gin.Context) {
	services.DownloadMeteoriteData()
	ginContext.Status(http.StatusAccepted)
}
