package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var whiteListedOrigin map[string]bool = map[string]bool{
	"http://localhost:5173": true,
}

func Cors() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		if whiteListedOrigin[ginContext.Request.Header.Get("Origin")] {
			ginContext.Writer.Header().Set("Access-Control-Allow-Origin", ginContext.Request.Header.Get("Origin"))
			ginContext.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			ginContext.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
			ginContext.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if ginContext.Request.Method == "OPTIONS" {
			ginContext.AbortWithStatus(http.StatusOK)
			return
		}
		ginContext.Next()
	}
}
