package middlewares

import (
	"main/configs"
	apperrors "main/internal/app/appErrors"
	"main/internal/app/constants"
	"main/internal/db"
	"main/internal/services/jwt"
	"main/internal/services/logger"

	"github.com/gin-gonic/gin"
)

func Authorization() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		token, err := ginContext.Cookie(constants.AuthCookieName)
		if err != nil {
			logger.AppMiddlewareLog.Infow("Error getting auth cookie", "error", err)
			errorResponse := apperrors.New("E005")
			ginContext.AbortWithStatusJSON(errorResponse.HttpStatusCode, errorResponse)
			return
		}
		payload, err := jwt.Verify(token, configs.APIJWTSecret)
		if err != nil {
			logger.AppMiddlewareLog.Infow("Error verifying jwt token", "error", err)
			errorResponse := apperrors.New("E005")
			ginContext.AbortWithStatusJSON(errorResponse.HttpStatusCode, errorResponse)
			return
		}

		username, ok := payload["username"]

		if !ok {
			logger.AppMiddlewareLog.Errorw("Error getting username from jwt payload", "error", err)
			errorResponse := apperrors.New("E001")
			ginContext.AbortWithStatusJSON(errorResponse.HttpStatusCode, errorResponse)
			return
		}

		csrfToken, ok := payload["csrfToken"]

		if !ok {
			logger.AppMiddlewareLog.Errorw("Error getting csrfToken from jwt payload", "error", err)
			errorResponse := apperrors.New("E001")
			ginContext.AbortWithStatusJSON(errorResponse.HttpStatusCode, errorResponse)
			return
		}

		csrfTokenFromRequest := ginContext.GetHeader(constants.CSRFTokenHeaderName)

		if csrfTokenFromRequest == "" || csrfTokenFromRequest != csrfToken {
			logger.AppMiddlewareLog.Infow("Error verifying csrfToken", "jwtCsrfToken", csrfToken, "csrfTokenFromRequest", csrfTokenFromRequest)
			errorResponse := apperrors.New("E005")
			ginContext.AbortWithStatusJSON(errorResponse.HttpStatusCode, errorResponse)
			return
		}

		user := db.Users{
			UserName: username.(string),
			Active:   true,
		}

		tx := db.DB.Find(&user)

		if tx.Error != nil {
			logger.AppMiddlewareLog.Errorw("Error while fetching user", "error", tx.Error)
			errorResponse := apperrors.New("E001")
			ginContext.AbortWithStatusJSON(errorResponse.HttpStatusCode, errorResponse)
			return
		}

		if tx.RowsAffected == 0 {
			logger.AppMiddlewareLog.Infow("No active user found", "username", username)
			errorResponse := apperrors.New("E004")
			ginContext.AbortWithStatusJSON(errorResponse.HttpStatusCode, errorResponse)
			return
		}

		ginContext.Set(constants.UserContextKey, user)
		ginContext.Next()
	}
}
