package controllers

import (
	"main/configs"
	apperrors "main/internal/app/appErrors"
	"main/internal/app/constants"
	"main/internal/app/utils"
	"main/internal/db"
	"main/internal/services"
	"main/internal/services/jwt"
	"main/internal/services/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginRequestBody struct {
	User     string `json:"user" validate:"required,min:8,max:16"`
	Password string `json:"password" validate:"required,min:8,max:16"`
}

func CreateUserSession(ginContext *gin.Context) {
	var body LoginRequestBody
	if err := ginContext.ShouldBindJSON(&body); err != nil {
		errResponse := apperrors.New("E003")
		ginContext.AbortWithStatusJSON(errResponse.HttpStatusCode, errResponse)
		return
	}
	user, userJwt, csrfToken, maxAge, cubeJwt, errorResponse := services.VerifyUserCredentials(body.User, body.Password)

	if errorResponse != nil {
		ginContext.AbortWithStatusJSON(errorResponse.HttpStatusCode, errorResponse)
		return
	}

	ginContext.SetSameSite(http.SameSiteStrictMode)
	ginContext.SetCookie(constants.AuthCookieName, userJwt, int(maxAge), "/", "", false, true)

	ginContext.JSON(http.StatusCreated, utils.GenerateSuccessResponse(map[string]interface{}{
		"user":         user,
		"cubeApiToken": cubeJwt,
		"csrfToken":    csrfToken,
	}))

}

func GetUserSession(ginContext *gin.Context) {
	value, ok := ginContext.Get("user")
	if !ok {
		logger.AppControllerLog.Errorw("Error getting user from gin context")
		errResponse := apperrors.New("E001")
		ginContext.AbortWithStatusJSON(errResponse.HttpStatusCode, errResponse)
		return
	}

	user, ok := value.(db.Users)

	if !ok {
		logger.AppControllerLog.Errorw("Error asserting user from gin context")
		errResponse := apperrors.New("E001")
		ginContext.AbortWithStatusJSON(errResponse.HttpStatusCode, errResponse)
		return
	}

	cubeApiToken, _, err := jwt.Sign(jwt.GenerationData{
		Key: configs.CubeAPISecret,
	})

	if err != nil {
		logger.AppControllerLog.Errorw("Error generating cube api token", "error", err)
		errResponse := apperrors.New("E001")
		ginContext.AbortWithStatusJSON(errResponse.HttpStatusCode, errResponse)
		return
	}

	ginContext.JSON(http.StatusOK, utils.GenerateSuccessResponse(map[string]interface{}{
		"user":         user,
		"cubeApiToken": cubeApiToken,
	}))
}

func DeleteUserSession(ginContext *gin.Context) {
	ginContext.SetCookie(constants.AuthCookieName, "", -1, "/", "", false, true)
	ginContext.JSON(http.StatusNoContent, utils.GenerateSuccessResponse(nil))
}
