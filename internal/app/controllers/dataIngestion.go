package controllers

import (
	"main/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleMeteoriteData(ginContext *gin.Context) {
	go services.ProcessMeteoriteData()
	ginContext.Status(http.StatusAccepted)
}
