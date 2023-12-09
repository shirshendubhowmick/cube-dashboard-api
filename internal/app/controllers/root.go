package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RootController(ginContext *gin.Context) {

	ginContext.JSON(http.StatusOK, gin.H{
		"message": "Welcome to cube dashboard API",
	})
}
